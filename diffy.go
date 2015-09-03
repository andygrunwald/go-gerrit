package diffy

import (
	"fmt"
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
