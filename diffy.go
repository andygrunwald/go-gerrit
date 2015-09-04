package diffy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// A Client manages communication with the Gerrit API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests.
	// BaseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// Services used for talking to different parts of the Gerrit API.
	Access   *AccessService
	Accounts *AccountsService
	Changes  *ChangesService
	Config   *ConfigService
	Groups   *GroupsService
	Plugins  *PluginsService
	Projects *ProjectsService
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
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	err = c.checkResponse(response)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			io.Copy(w, response.Body)
		} else {
			body, err := ioutil.ReadAll(response.Body)
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
func (c *Client) checkResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}

	err := fmt.Errorf("API call failed: %s", r.Status)
	return err
}
