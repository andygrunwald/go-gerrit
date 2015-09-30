package gerrit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/google/go-querystring/query"
)

// TODO Try to reduce the code duplications of a std API req
// Maybe with http://play.golang.org/p/j-667shCCB
// and https://groups.google.com/forum/#!topic/golang-nuts/D-gIr24k5uY

// A Client manages communication with the Gerrit API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// Gerrit service for authentication
	Authentication *AuthenticationService

	// Services used for talking to different parts of the Gerrit API.
	Access   *AccessService
	Accounts *AccountsService
	Changes  *ChangesService
	Config   *ConfigService
	Groups   *GroupsService
	Plugins  *PluginsService
	Projects *ProjectsService
}

// Response is a Gerrit API response.
// This wraps the standard http.Response returned from Gerrit.
type Response struct {
	*http.Response
}

// NewClient returns a new Gerrit API client.
// gerritInstance has to be the HTTP endpoint of the Gerrit instance.
// If a nil httpClient is provided, http.DefaultClient will be used.
func NewClient(gerritURL string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if len(gerritURL) == 0 {
		return nil, fmt.Errorf("No Gerrit instance given.")
	}
	baseURL, err := url.Parse(gerritURL)
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:  httpClient,
		baseURL: baseURL,
	}
	c.Authentication = &AuthenticationService{client: c}
	c.Access = &AccessService{client: c}
	c.Accounts = &AccountsService{client: c}
	c.Changes = &ChangesService{client: c}
	c.Config = &ConfigService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Plugins = &PluginsService{client: c}
	c.Projects = &ProjectsService{client: c}

	return c, nil
}

// NewRequest creates an API request.
// A relative URL can be provided in urlStr, in which case it is resolved relative to the baseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	// Build URL for request
	u, err := c.buildURLForRequest(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}

	// Apply Authentication
	c.addAuthentication(req)

	// Request compact JSON
	// See https://gerrit-review.googlesource.com/Documentation/rest-api.html#output
	req.Header.Add("Accept", "application/json")

	// TODO: Add gzip encoding
	// Accept-Encoding request header is set to gzip
	// See https://gerrit-review.googlesource.com/Documentation/rest-api.html#output

	return req, nil
}

// Call is a combine function for Client.NewRequest and Client.Do.
//
// Most API methods are quite the same.
// Get the URL, apply options, make a request, and get the response.
// Without adding special headers or something.
// To avoid a big amount of code duplication you can Client.Call.
//
// method is the HTTP method you want to call.
// u is the URL you want to call.
// body is the HTTP body.
// v is the HTTP response.
//
// For more information read https://github.com/google/go-github/issues/234
func (c *Client) Call(method, u string, body interface{}, v interface{}) (*Response, error) {
	req, err := c.NewRequest(method, u, body)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req, v)
	if err != nil {
		return resp, err
	}

	return resp, err
}

// buildURLForRequest will build the URL (as string) that will be called.
// We need such a utility method, because the URL.Path needs to be escaped (partly).
//
// E.g. if a project is called via "projects/%s" and the project is named "plugin/delete-project"
// there has to be "projects/plugin%25Fdelete-project" instead of "projects/plugin/delete-project".
// The second url will return nothing.
func (c *Client) buildURLForRequest(urlStr string) (string, error) {
	u := c.baseURL.String()

	// If there is no / at the end, add one
	if strings.HasSuffix(u, "/") == false {
		u += "/"
	}

	// If there is a "/" at the start, remove it
	if strings.HasPrefix(urlStr, "/") == true {
		urlStr = urlStr[1:]
	}

	// If we are authenticated, lets apply the a/ prefix
	if c.Authentication.HasAuth() == true {
		urlStr = "a/" + urlStr
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}
	u += rel.String()

	return u, nil
}

// Do sends an API request and returns the API response.
// The API response is JSON decoded and stored in the value pointed to by v,
// or returned as an error if an API error has occurred.
// If v implements the io.Writer interface, the raw response body will be written to v,
// without attempting to first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// Wrap response
	response := &Response{Response: resp}

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, resp.Body)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				// even though there was an error, we still return the response
				// in case the caller wants to inspect it further
				return response, err
			}

			body = RemoveMagicPrefixLine(body)
			err = json.Unmarshal(body, v)
		}
	}
	return response, err
}

func (c *Client) addAuthentication(req *http.Request) {
	// Apply HTTP Basic Authentication
	if c.Authentication.HasBasicAuth() == true {
		req.SetBasicAuth(c.Authentication.name, c.Authentication.secret)
	}
	// Apply HTTP Cookie
	if c.Authentication.HasCookieAuth() == true {
		req.AddCookie(&http.Cookie{
			Name:  c.Authentication.name,
			Value: c.Authentication.secret,
		})
	}
}

// DeleteRequest sends an DELETE API Request to urlStr with optional body.
// It is a shorthand combination for Client.NewRequest with Client.Do.
//
// Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the request body.
func (c *Client) DeleteRequest(urlStr string, body interface{}) (*Response, error) {
	req, err := c.NewRequest("DELETE", urlStr, body)
	if err != nil {
		return nil, err
	}

	return c.Do(req, nil)
}

// RemoveMagicPrefixLine removed the "magic prefix line" of Gerris JSON response.
// the JSON response body starts with a magic prefix line that must be stripped before feeding the rest of the response body to a JSON parser.
// The reason for this is to prevent against Cross Site Script Inclusion (XSSI) attacks.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api.html#output
func RemoveMagicPrefixLine(body []byte) []byte {
	index := bytes.IndexByte(body, '\n')
	if index > -1 {
		// +1 to catch the \n as well
		body = body[(index + 1):]
	}

	return body
}

// CheckResponse checks the API response for errors, and returns them if present.
// A response is considered an error if it has a status code outside the 200 range.
// API error responses are expected to have no response body.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api.html#response-codes
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	// Some calls require an authentification
	// In such cases errors like:
	// 		API call to https://review.typo3.org/accounts/self failed: 403 Forbidden
	// will be thrown.

	err := fmt.Errorf("API call to %s failed: %s", r.Request.URL.String(), r.Status)
	return err
}

// addOptions adds the parameters in opt as URL query parameters to s.
// opt must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

// getStringResponseWithoutOptions retrieved a single string Response for a GET request
func getStringResponseWithoutOptions(client *Client, u string) (string, *Response, error) {
	v := new(string)
	resp, err := client.Call("GET", u, nil, v)
	return *v, resp, err
}
