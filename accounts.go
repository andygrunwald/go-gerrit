package diffy

// AccountsService contains Account related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-accounts.html
type AccountsService struct {
	client *Client
}

// AccountInfo entity contains information about an account.
type AccountInfo struct {
	AccountID int    `json:"_account_id"`
	Name      string `json:"name,omitempty"`
	Email     string `json:"email,omitempty"`
	Username  string `json:"username,omitempty"`
}
