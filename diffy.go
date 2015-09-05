package diffy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"

	"github.com/google/go-querystring/query"
)

// A Client manages communication with the Gerrit API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// Diffy service for authentication
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
func NewClient(gerritInstance string, httpClient *http.Client) (*Client, error) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	if len(gerritInstance) == 0 {
		return nil, fmt.Errorf("No Gerrit instance given.")
	}
	baseURL, err := url.Parse(gerritInstance)
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
	// If we are authenticated, lets apply the a/ prefix
	if c.Authentication.HasAuth() == true {
		urlStr = "a/" + urlStr
	}

	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
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
