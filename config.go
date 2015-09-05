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

// AuthInfo entity contains information about the authentication configuration of the Gerrit server.
type AuthInfo struct {
	Type                     string   `json:"type"`
	UseContributorAgreements bool     `json:"use_contributor_agreements,omitempty"`
	EditableAccountFields    []string `json:"editable_account_fields"`
	LoginURL                 string   `json:"login_url,omitempty"`
	LoginText                string   `json:"login_text,omitempty"`
	SwitchAccountURL         string   `json:"switch_account_url,omitempty"`
	RegisterURL              string   `json:"register_url,omitempty"`
	RegisterText             string   `json:"register_text,omitempty"`
	EditFullNameURL          string   `json:"edit_full_name_url,omitempty"`
	HTTPPasswordURL          string   `json:"http_password_url,omitempty"`
	IsGitBasicAuth           bool     `json:"is_git_basic_auth,omitempty"`
}

// CacheInfo entity contains information about a cache.
type CacheInfo struct {
	Name       string       `json:"name,omitempty"`
	Type       string       `json:"type"`
	Entries    EntriesInfo  `json:"entries"`
	AverageGet string       `json:"average_get,omitempty"`
	HitRatio   HitRatioInfo `json:"hit_ratio"`
}

// CacheOperationInput entity contains information about an operation that should be executed on caches.
type CacheOperationInput struct {
	Operation string   `json:"operation"`
	Caches    []string `json:"caches,omitempty"`
}

// CapabilityInfo entity contains information about a capability.type CapabilityInfo struct {
type CapabilityInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// HitRatioInfo entity contains information about the hit ratio of a cache.
type HitRatioInfo struct {
	Mem  int `json:"mem"`
	Disk int `json:"disk,omitempty"`
}

// EntriesInfo entity contains information about the entries in a cache.
type EntriesInfo struct {
	Mem   int    `json:"mem,omitempty"`
	Disk  int    `json:"disk,omitempty"`
	Space string `json:"space,omitempty"`
}

// UserConfigInfo entity contains information about Gerrit configuration from the user section.
type UserConfigInfo struct {
	AnonymousCowardName string `json:"anonymous_coward_name"`
}

// TopMenuEntryInfo entity contains information about a top menu entry.
type TopMenuEntryInfo struct {
	Name  string            `json:"name"`
	Items []TopMenuItemInfo `json:"items"`
}

// ThreadSummaryInfo entity contains information about the current threads.
type ThreadSummaryInfo struct {
	CPUs    int                       `json:"cpus"`
	Threads int                       `json:"threads"`
	Counts  map[string]map[string]int `json:"counts"`
}

// TaskSummaryInfo entity contains information about the current tasks.
type TaskSummaryInfo struct {
	Total    int `json:"total,omitempty"`
	Running  int `json:"running,omitempty"`
	Ready    int `json:"ready,omitempty"`
	Sleeping int `json:"sleeping,omitempty"`
}

// TaskInfo entity contains information about a task in a background work queue.
type TaskInfo struct {
	ID         string `json:"id"`
	State      string `json:"state"`
	StartTime  string `json:"start_time"`
	Delay      int    `json:"delay"`
	Command    string `json:"command"`
	RemoteName string `json:"remote_name,omitempty"`
	Project    string `json:"project,omitempty"`
}

// SummaryInfo entity contains information about the current state of the server.
type SummaryInfo struct {
	TaskSummary   TaskSummaryInfo `json:"task_summary"`
	MemSummary    MemSummaryInfo  `json:"mem_summary"`
	ThreadSummary ThemeInfo       `json:"thread_summary"`
	JVMSummary    JvmSummaryInfo  `json:"jvm_summary,omitempty"`
}

// SuggestInfo entity contains information about Gerrit configuration from the suggest section.
type SuggestInfo struct {
	From int `json:"from"`
}

// SSHdInfo entity contains information about Gerrit configuration from the sshd section.
type SSHdInfo struct{}

// ServerInfo entity contains information about the configuration of the Gerrit server.
type ServerInfo struct {
	Auth       AuthInfo          `json:"auth"`
	Change     ChangeConfigInfo  `json:"change"`
	Download   DownloadInfo      `json:"download"`
	Gerrit     GerritInfo        `json:"gerrit"`
	Gitweb     map[string]string `json:"gitweb,omitempty"`
	Plugin     PluginConfigInfo  `json:"plugin"`
	Receive    ReceiveInfo       `json:"receive,omitempty"`
	SSHd       SSHdInfo          `json:"sshd,omitempty"`
	Suggest    SuggestInfo       `json:"suggest"`
	URLAliases map[string]string `json:"url_aliases,omitempty"`
	User       UserConfigInfo    `json:"user"`
}

// ReceiveInfo entity contains information about the configuration of git-receive-pack behavior on the server.
type ReceiveInfo struct {
	EnableSignedPush bool `json:"enableSignedPush,omitempty"`
}

// PluginConfigInfo entity contains information about Gerrit extensions by plugins.
type PluginConfigInfo struct {
	HasAvatars bool `json:"has_avatars,omitempty"`
}

// MemSummaryInfo entity contains information about the current memory usage.
type MemSummaryInfo struct {
	Total     string `json:"total"`
	Used      string `json:"used"`
	Free      string `json:"free"`
	Buffers   string `json:"buffers"`
	Max       string `json:"max"`
	OpenFiles int    `json:"open_files,omitempty"`
}

// JvmSummaryInfo entity contains information about the JVM.
type JvmSummaryInfo struct {
	VMVendor                string `json:"vm_vendor"`
	VMName                  string `json:"vm_name"`
	VMVersion               string `json:"vm_version"`
	OSName                  string `json:"os_name"`
	OSVersion               string `json:"os_version"`
	OSArch                  string `json:"os_arch"`
	User                    string `json:"user"`
	Host                    string `json:"host,omitempty"`
	CurrentWorkingDirectory string `json:"current_working_directory"`
	Site                    string `json:"site"`
}

// GerritInfo entity contains information about Gerrit configuration from the gerrit section.
type GerritInfo struct {
	AllProjectsName string `json:"all_projects_name"`
	AllUsersName    string `json:"all_users_name"`
	DocURL          string `json:"doc_url,omitempty"`
	ReportBugURL    string `json:"report_bug_url,omitempty"`
	ReportBugText   string `json:"report_bug_text,omitempty"`
}

// GitwebInfo entity contains information about the gitweb configuration.
type GitwebInfo struct {
	URL  string         `json:"url"`
	Type GitwebTypeInfo `json:"type"`
}

// GitwebTypeInfo entity contains information about the gitweb configuration.
type GitwebTypeInfo struct {
	Name          string `json:"name"`
	Revision      string `json:"revision,omitempty"`
	Project       string `json:"project,omitempty"`
	Branch        string `json:"branch,omitempty"`
	RootTree      string `json:"root_tree,omitempty"`
	File          string `json:"file,omitempty"`
	FileHistory   string `json:"file_history,omitempty"`
	PathSeparator string `json:"path_separator"`
	LinkDrafts    bool   `json:"link_drafts,omitempty"`
	URLEncode     bool   `json:"url_encode,omitempty"`
}

// EmailConfirmationInput entity contains information for confirming an email address.
type EmailConfirmationInput struct {
	Token string `json:"token"`
}

// DownloadSchemeInfo entity contains information about a supported download scheme and its commands.
type DownloadSchemeInfo struct {
	URL             string            `json:"url"`
	IsAuthRequired  bool              `json:"is_auth_required,omitempty"`
	IsAuthSupported bool              `json:"is_auth_supported,omitempty"`
	Commands        map[string]string `json:"commands"`
	CloneCommands   map[string]string `json:"clone_commands"`
}

// DownloadInfo entity contains information about supported download options.
type DownloadInfo struct {
	Schemes  map[string]DownloadSchemeInfo `json:"schemes"`
	Archives []string                      `json:"archives"`
}

// ChangeConfigInfo entity contains information about Gerrit configuration from the change section.
type ChangeConfigInfo struct {
	AllowDrafts      bool   `json:"allow_drafts,omitempty"`
	LargeChange      int    `json:"large_change"`
	ReplyLabel       string `json:"reply_label"`
	ReplyTooltip     string `json:"reply_tooltip"`
	UpdateDelay      int    `json:"update_delay"`
	SubmitWholeTopic bool   `json:"submit_whole_topic"`
}
