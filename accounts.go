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

// SSHKeyInfo entity contains information about an SSH key of a user.
type SSHKeyInfo struct {
	Seq          int    `json:"seq"`
	SSHPublicKey string `json:"ssh_public_key"`
	EncodedKey   string `json:"encoded_key"`
	Algorithm    string `json:"algorithm"`
	Comment      string `json:"comment,omitempty"`
	Valid        bool   `json:"valid"`
}

// UsernameInput entity contains information for setting the username for an account.
type UsernameInput struct {
	Username string `json:"username"`
}

// QueryLimitInfo entity contains information about the Query Limit of a user.
type QueryLimitInfo struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

// HTTPPasswordInput entity contains information for setting/generating an HTTP password.
type HTTPPasswordInput struct {
	Generate     bool   `json:"generate,omitempty"`
	HTTPPassword string `json:"http_password,omitempty"`
}

// GpgKeysInput entity contains information for adding/deleting GPG keys.
type GpgKeysInput struct {
	Add    []string `json:"add"`
	Delete []string `json:"delete"`
}

// GpgKeyInfo entity contains information about a GPG public key.
type GpgKeyInfo struct {
	ID          string   `json:"id,omitempty"`
	Fingerprint string   `json:"fingerprint,omitempty"`
	UserIDs     []string `json:"user_ids,omitempty"`
	Key         string   `json:"key,omitempty"`
}

// EmailInput entity contains information for registering a new email address.
type EmailInput struct {
	Email          string `json:"email"`
	Preferred      bool   `json:"preferred,omitempty"`
	NoConfirmation bool   `json:"no_confirmation,omitempty"`
}

// EmailInfo entity contains information about an email address of a user.
type EmailInfo struct {
	Email               string `json:"email"`
	Preferred           bool   `json:"preferred,omitempty"`
	PendingConfirmation bool   `json:"pending_confirmation,omitempty"`
}

// AccountInput entity contains information for the creation of a new account.
type AccountInput struct {
	Username     string   `json:"username,omitempty"`
	Name         string   `json:"name,omitempty"`
	Email        string   `json:"email,omitempty"`
	SSHKey       string   `json:"ssh_key,omitempty"`
	HTTPPassword string   `json:"http_password,omitempty"`
	Groups       []string `json:"groups,omitempty"`
}

// AccountDetailInfo entity contains detailled information about an account.
type AccountDetailInfo struct {
	AccountInfo
	RegisteredOn string `json:"registered_on"`
}

// AccountNameInput entity contains information for setting a name for an account.
type AccountNameInput struct {
	Name string `json:"name,omitempty"`
}

// CapabilityInfo entity contains information about the global capabilities of a user.
type CapabilityInfo struct {
	AccessDatabase     bool           `json:"accessDatabase,omitempty"`
	AdministrateServer bool           `json:"administrateServer,omitempty"`
	CreateAccount      bool           `json:"createAccount,omitempty"`
	CreateGroup        bool           `json:"createGroup,omitempty"`
	CreateProject      bool           `json:"createProject,omitempty"`
	EmailReviewers     bool           `json:"emailReviewers,omitempty"`
	FlushCaches        bool           `json:"flushCaches,omitempty"`
	KillTask           bool           `json:"killTask,omitempty"`
	MaintainServer     bool           `json:"maintainServer,omitempty"`
	Priority           string         `json:"priority,omitempty"`
	QueryLimit         QueryLimitInfo `json:"queryLimit"`
	RunAs              bool           `json:"runAs,omitempty"`
	RunGC              bool           `json:"runGC,omitempty"`
	StreamEvents       bool           `json:"streamEvents,omitempty"`
	ViewAllAccounts    bool           `json:"viewAllAccounts,omitempty"`
	ViewCaches         bool           `json:"viewCaches,omitempty"`
	ViewConnections    bool           `json:"viewConnections,omitempty"`
	ViewPlugins        bool           `json:"viewPlugins,omitempty"`
	ViewQueue          bool           `json:"viewQueue,omitempty"`
}

// DiffPreferencesInfo entity contains information about the diff preferences of a user.
type DiffPreferencesInfo struct {
	Context                 int    `json:"context"`
	Theme                   string `json:"theme"`
	ExpandAllComments       bool   `json:"expand_all_comments,omitempty"`
	IgnoreWhitespace        string `json:"ignore_whitespace"`
	IntralineDifference     bool   `json:"intraline_difference,omitempty"`
	LineLength              int    `json:"line_length"`
	ManualReview            bool   `json:"manual_review,omitempty"`
	RetainHeader            bool   `json:"retain_header,omitempty"`
	ShowLineEndings         bool   `json:"show_line_endings,omitempty"`
	ShowTabs                bool   `json:"show_tabs,omitempty"`
	ShowWhitespaceErrors    bool   `json:"show_whitespace_errors,omitempty"`
	SkipDeleted             bool   `json:"skip_deleted,omitempty"`
	SkipUncommented         bool   `json:"skip_uncommented,omitempty"`
	SyntaxHighlighting      bool   `json:"syntax_highlighting,omitempty"`
	HideTopMenu             bool   `json:"hide_top_menu,omitempty"`
	AutoHideDiffTableHeader bool   `json:"auto_hide_diff_table_header,omitempty"`
	HideLineNumbers         bool   `json:"hide_line_numbers,omitempty"`
	TabSize                 int    `json:"tab_size"`
	HideEmptyPane           bool   `json:"hide_empty_pane,omitempty"`
}

// DiffPreferencesInput entity contains information for setting the diff preferences of a user.
// Fields which are not set will not be updated.
type DiffPreferencesInput struct {
	Context                 int    `json:"context,omitempty"`
	ExpandAllComments       bool   `json:"expand_all_comments,omitempty"`
	IgnoreWhitespace        string `json:"ignore_whitespace,omitempty"`
	IntralineDifference     bool   `json:"intraline_difference,omitempty"`
	LineLength              int    `json:"line_length,omitempty"`
	ManualReview            bool   `json:"manual_review,omitempty"`
	RetainHeader            bool   `json:"retain_header,omitempty"`
	ShowLineEndings         bool   `json:"show_line_endings,omitempty"`
	ShowTabs                bool   `json:"show_tabs,omitempty"`
	ShowWhitespaceErrors    bool   `json:"show_whitespace_errors,omitempty"`
	SkipDeleted             bool   `json:"skip_deleted,omitempty"`
	SkipUncommented         bool   `json:"skip_uncommented,omitempty"`
	SyntaxHighlighting      bool   `json:"syntax_highlighting,omitempty"`
	HideTopMenu             bool   `json:"hide_top_menu,omitempty"`
	AutoHideDiffTableHeader bool   `json:"auto_hide_diff_table_header,omitempty"`
	HideLineNumbers         bool   `json:"hide_line_numbers,omitempty"`
	TabSize                 int    `json:"tab_size,omitempty"`
}

// PreferencesInfo entity contains information about a user’s preferences.
type PreferencesInfo struct {
	ChangesPerPage            int               `json:"changes_per_page"`
	ShowSiteHeader            bool              `json:"show_site_header,omitempty"`
	UseFlashClipboard         bool              `json:"use_flash_clipboard,omitempty"`
	DownloadScheme            string            `json:"download_scheme"`
	DownloadCommand           string            `json:"download_command"`
	CopySelfOnEmail           bool              `json:"copy_self_on_email,omitempty"`
	DateFormat                string            `json:"date_format"`
	TimeFormat                string            `json:"time_format"`
	RelativeDateInChangeTable bool              `json:"relative_date_in_change_table,omitempty"`
	SizeBarInChangeTable      bool              `json:"size_bar_in_change_table,omitempty"`
	LegacycidInChangeTable    bool              `json:"legacycid_in_change_table,omitempty"`
	MuteCommonPathPrefixes    bool              `json:"mute_common_path_prefixes,omitempty"`
	ReviewCategoryStrategy    string            `json:"review_category_strategy"`
	DiffView                  string            `json:"diff_view"`
	My                        []TopMenuItemInfo `json:"my"`
	URLAliases                string            `json:"url_aliases,omitempty"`
}

// PreferencesInput entity contains information for setting the user preferences.
// Fields which are not set will not be updated.
type PreferencesInput struct {
	ChangesPerPage            int               `json:"changes_per_page,omitempty"`
	ShowSiteHeader            bool              `json:"show_site_header,omitempty"`
	UseFlashClipboard         bool              `json:"use_flash_clipboard,omitempty"`
	DownloadScheme            string            `json:"download_scheme,omitempty"`
	DownloadCommand           string            `json:"download_command,omitempty"`
	CopySelfOnEmail           bool              `json:"copy_self_on_email,omitempty"`
	DateFormat                string            `json:"date_format,omitempty"`
	TimeFormat                string            `json:"time_format,omitempty"`
	RelativeDateInChangeTable bool              `json:"relative_date_in_change_table,omitempty"`
	SizeBarInChangeTable      bool              `json:"size_bar_in_change_table,omitempty"`
	LegacycidInChangeTable    bool              `json:"legacycid_in_change_table,omitempty"`
	MuteCommonPathPrefixes    bool              `json:"mute_common_path_prefixes,omitempty"`
	ReviewCategoryStrategy    string            `json:"review_category_strategy,omitempty"`
	DiffView                  string            `json:"diff_view,omitempty"`
	My                        []TopMenuItemInfo `json:"my,omitempty"`
	URLAliases                string            `json:"url_aliases,omitempty"`
}
