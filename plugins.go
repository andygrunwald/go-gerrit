package diffy

// PluginsService contains Plugin related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-plugins.html
type PluginsService struct {
	client *Client
}

// PluginInfo entity describes a plugin.
type PluginInfo struct {
	ID       string `json:"id"`
	Version  string `json:"version"`
	IndexURL string `json:"index_url,omitempty"`
	Disabled bool   `json:"disabled,omitempty"`
}

// PluginInput entity describes a plugin that should be installed.
type PluginInput struct {
	URL string `json:"url"`
}
