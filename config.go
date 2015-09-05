package diffy

// ConfigService contains Config related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-config.html
type ConfigService struct {
	client *Client
}

// TopMenuItemInfo entity contains information about a menu item in a top menu entry.
type TopMenuItemInfo struct {
	URL    string `json:"url"`
	Name   string `json:"name"`
	Target string `json:"target"`
	ID     string `json:"id,omitempty"`
}
