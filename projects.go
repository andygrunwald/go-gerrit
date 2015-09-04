package diffy

// ProjectsService contains Project related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html
type ProjectsService struct {
	client *Client
}

// ProjectInfo entity contains information about a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#project-info
type ProjectInfo struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Parent      string            `json:"parent"`
	Description string            `json:"description"`
	State       string            `json:"state"`
	Branches    map[string]string `json:"branches"`
	WebLinks    []WebLinkInfo     `json:"web_links"`
}

// ProjectInput entity contains information for the creation of a new project.
type ProjectInput struct {
	Name                             string                       `json:"name"`
	Parent                           string                       `json:"parent"`
	Description                      string                       `json:"description"`
	PermissionsOnly                  bool                         `json:"permissions_only"`
	CreateEmptyCommit                bool                         `json:"create_empty_commit"`
	SubmitType                       string                       `json:"submit_type"`
	Branches                         []string                     `json:"branches"`
	Owners                           []string                     `json:"owners"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target"`
	UseContentMerge                  string                       `json:"use_content_merge"`
	RequireChangeID                  string                       `json:"require_change_id"`
	MaxObjectSizeLimit               string                       `json:"max_object_size_limit"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values"`
}

// GCInput entity contains information to run the Git garbage collection.
type GCInput struct {
	ShowProgress bool `json:"show_progress"`
	Aggressive   bool `json:"aggressive"`
}

// HeadInput entity contains information for setting HEAD for a project.
type HeadInput struct {
	Ref string `json:"ref"`
}

// DeleteBranchesInput entity contains information about branches that should be deleted.
type DeleteBranchesInput struct {
	Branches []string `json:"DeleteBranchesInput"`
}

// DashboardSectionInfo entity contains information about a section in a dashboard.
type DashboardSectionInfo struct {
	Name  string `json:"name"`
	Query string `json:"query"`
}

// DashboardInput entity contains information to create/update a project dashboard.
type DashboardInput struct {
	ID            string `json:"id"`
	CommitMessage string `json:"commit_message"`
}

// DashboardInfo entity contains information about a project dashboard.
type DashboardInfo struct {
	ID              string                 `json:"id"`
	Project         string                 `json:"project"`
	DefiningProject string                 `json:"defining_project"`
	Ref             string                 `json:"ref"`
	Path            string                 `json:"path"`
	Description     string                 `json:"description"`
	Foreach         string                 `json:"foreach"`
	URL             string                 `json:"url"`
	Default         bool                   `json:"default"`
	Title           string                 `json:"title"`
	Sections        []DashboardSectionInfo `json:"sections"`
}

// BanInput entity contains information for banning commits in a project.
type BanInput struct {
	Commits []string `json:"commits"`
	Reason  string   `json:"reason"`
}

// BanResultInfo entity describes the result of banning commits.
type BanResultInfo struct {
	NewlyBanned   []string `json:"newly_banned"`
	AlreadyBanned []string `json:"already_banned"`
	Ignored       []string `json:"ignored"`
}

// BranchInfo entity contains information about a branch.
type BranchInfo struct {
	Ref       string        `json:"ref"`
	Revision  string        `json:"revision"`
	CanDelete bool          `json:"can_delete"`
	WebLinks  []WebLinkInfo `json:"web_links"`
}

// BranchInput entity contains information for the creation of a new branch.
type BranchInput struct {
	Ref      string `json:"ref"`
	Revision string `json:"revision"`
}

// ThemeInfo entity describes a theme.
type ThemeInfo struct {
	CSS    string `type:"css"`
	Header string `type:"header"`
	Footer string `type:"footer"`
}

// TagInfo entity contains information about a tag.
type TagInfo struct {
	Ref      string        `json:"ref"`
	Revision string        `json:"revision"`
	Object   string        `json:"object"`
	Message  string        `json:"message"`
	Tagger   GitPersonInfo `json:"tagger"`
}

// ReflogEntryInfo entity describes an entry in a reflog.
type ReflogEntryInfo struct {
	OldID   string        `json:"old_id"`
	NewID   string        `json:"new_id"`
	Who     GitPersonInfo `json:"who"`
	Comment string        `json:"comment"`
}

// ProjectParentInput entity contains information for setting a project parent.
type ProjectParentInput struct {
	Parent        string `json:"ProjectParentInput"`
	CommitMessage string `json:"commit_message"`
}

// RepositoryStatisticsInfo entity contains information about statistics of a Git repository.
type RepositoryStatisticsInfo struct {
	NumberOfLooseObjects  int `json:"number_of_loose_objects"`
	NumberOfLooseRefs     int `json:"number_of_loose_refs"`
	NumberOfPackFiles     int `json:"number_of_pack_files"`
	NumberOfPackedObjects int `json:"number_of_packed_objects"`
	NumberOfPackedRefs    int `json:"number_of_packed_refs"`
	SizeOfLooseObjects    int `json:"size_of_loose_objects"`
	SizeOfPackedObjects   int `json:"size_of_packed_objects"`
}

// InheritedBooleanInfo entity represents a boolean value that can also be inherited.
type InheritedBooleanInfo struct {
	Value           bool   `json:"value"`
	ConfiguredValue string `json:"configured_value"`
	InheritedValue  bool   `json:"inherited_value"`
}

// MaxObjectSizeLimitInfo entity contains information about the max object size limit of a project.
type MaxObjectSizeLimitInfo struct {
	Value           string `json:"value"`
	ConfiguredValue string `json:"configured_value"`
	InheritedValue  string `json:"inherited_value"`
}

// ConfigParameterInfo entity describes a project configuration parameter.
type ConfigParameterInfo struct {
	DisplayName string   `json:"display_name"`
	Description string   `json:"description"`
	Warning     string   `json:"warning"`
	Type        string   `json:"type"`
	Value       string   `json:"value"`
	Values      []string `json:"values"`
	// TODO: 5 fields are missing here, because the documentation seems to be fucked up
	// See https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#config-parameter-info
}

// ProjectDescriptionInput entity contains information for setting a project description.
type ProjectDescriptionInput struct {
	Description   string `json:"description"`
	CommitMessage string `json:"commit_message"`
}

// ConfigInfo entity contains information about the effective project configuration.
type ConfigInfo struct {
	Description                      string                         `json:"description"`
	UseContributorAgreements         InheritedBooleanInfo           `json:"use_contributor_agreements"`
	UseContentMerge                  InheritedBooleanInfo           `json:"use_content_merge"`
	UseSignedOffBy                   InheritedBooleanInfo           `json:"use_signed_off_by"`
	CreateNewChangeForAllNotInTarget InheritedBooleanInfo           `json:"create_new_change_for_all_not_in_target"`
	RequireChangeID                  InheritedBooleanInfo           `json:"require_change_id"`
	EnableSignedPush                 InheritedBooleanInfo           `json:"enable_signed_push"`
	MaxObjectSizeLimit               MaxObjectSizeLimitInfo         `json:"max_object_size_limit"`
	SubmitType                       string                         `json:"submit_type"`
	State                            string                         `json:"state"`
	Commentlinks                     map[string]string              `json:"commentlinks"`
	Theme                            ThemeInfo                      `json:"theme"`
	PluginConfig                     map[string]ConfigParameterInfo `json:"plugin_config"`
	Actions                          map[string]ActionInfo          `json:"actions"`
}

// ConfigInput entity describes a new project configuration.
type ConfigInput struct {
	Description                      string                       `json:"description"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements"`
	UseContentMerge                  string                       `json:"use_content_merge"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target"`
	RequireChangeID                  string                       `json:"require_change_id"`
	MaxObjectSizeLimit               MaxObjectSizeLimitInfo       `json:"max_object_size_limit"`
	SubmitType                       string                       `json:"submit_type"`
	State                            string                       `json:"state"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values"`
}
