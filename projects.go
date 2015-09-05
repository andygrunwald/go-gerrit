package diffy

import (
	"fmt"
)

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
	Parent      string            `json:"parent,omitempty"`
	Description string            `json:"description,omitempty"`
	State       string            `json:"state,omitempty"`
	Branches    map[string]string `json:"branches,omitempty"`
	WebLinks    []WebLinkInfo     `json:"web_links,omitempty"`
}

// ProjectInput entity contains information for the creation of a new project.
type ProjectInput struct {
	Name                             string                       `json:"name,omitempty"`
	Parent                           string                       `json:"parent,omitempty"`
	Description                      string                       `json:"description,omitempty"`
	PermissionsOnly                  bool                         `json:"permissions_only"`
	CreateEmptyCommit                bool                         `json:"create_empty_commit"`
	SubmitType                       string                       `json:"submit_type,omitempty"`
	Branches                         []string                     `json:"branches,omitempty"`
	Owners                           []string                     `json:"owners,omitempty"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target"`
	UseContentMerge                  string                       `json:"use_content_merge"`
	RequireChangeID                  string                       `json:"require_change_id"`
	MaxObjectSizeLimit               string                       `json:"max_object_size_limit,omitempty"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values,omitempty"`
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
	ID            string `json:"id,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
}

// DashboardInfo entity contains information about a project dashboard.
type DashboardInfo struct {
	ID              string                 `json:"id"`
	Project         string                 `json:"project"`
	DefiningProject string                 `json:"defining_project"`
	Ref             string                 `json:"ref"`
	Path            string                 `json:"path"`
	Description     string                 `json:"description,omitempty"`
	Foreach         string                 `json:"foreach,omitempty"`
	URL             string                 `json:"url"`
	Default         bool                   `json:"default"`
	Title           string                 `json:"title,omitempty"`
	Sections        []DashboardSectionInfo `json:"sections"`
}

// BanInput entity contains information for banning commits in a project.
type BanInput struct {
	Commits []string `json:"commits"`
	Reason  string   `json:"reason,omitempty"`
}

// BanResultInfo entity describes the result of banning commits.
type BanResultInfo struct {
	NewlyBanned   []string `json:"newly_banned,omitempty"`
	AlreadyBanned []string `json:"already_banned,omitempty"`
	Ignored       []string `json:"ignored,omitempty"`
}

// BranchInfo entity contains information about a branch.
type BranchInfo struct {
	Ref       string        `json:"ref"`
	Revision  string        `json:"revision"`
	CanDelete bool          `json:"can_delete"`
	WebLinks  []WebLinkInfo `json:"web_links,omitempty"`
}

// BranchInput entity contains information for the creation of a new branch.
type BranchInput struct {
	Ref      string `json:"ref,omitempty"`
	Revision string `json:"revision,omitempty"`
}

// ThemeInfo entity describes a theme.
type ThemeInfo struct {
	CSS    string `type:"css,omitempty"`
	Header string `type:"header,omitempty"`
	Footer string `type:"footer,omitempty"`
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
	Parent        string `json:"parent"`
	CommitMessage string `json:"commit_message,omitempty"`
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
	InheritedValue  bool   `json:"inherited_value,omitempty"`
}

// MaxObjectSizeLimitInfo entity contains information about the max object size limit of a project.
type MaxObjectSizeLimitInfo struct {
	Value           string `json:"value,omitempty"`
	ConfiguredValue string `json:"configured_value,omitempty"`
	InheritedValue  string `json:"inherited_value,omitempty"`
}

// ConfigParameterInfo entity describes a project configuration parameter.
type ConfigParameterInfo struct {
	DisplayName string   `json:"display_name,omitempty"`
	Description string   `json:"description,omitempty"`
	Warning     string   `json:"warning,omitempty"`
	Type        string   `json:"type"`
	Value       string   `json:"value,omitempty"`
	Values      []string `json:"values,omitempty"`
	// TODO: 5 fields are missing here, because the documentation seems to be fucked up
	// See https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#config-parameter-info
}

// ProjectDescriptionInput entity contains information for setting a project description.
type ProjectDescriptionInput struct {
	Description   string `json:"description,omitempty"`
	CommitMessage string `json:"commit_message,omitempty"`
}

// ConfigInfo entity contains information about the effective project configuration.
type ConfigInfo struct {
	Description                      string                         `json:"description,omitempty"`
	UseContributorAgreements         InheritedBooleanInfo           `json:"use_contributor_agreements,omitempty"`
	UseContentMerge                  InheritedBooleanInfo           `json:"use_content_merge,omitempty"`
	UseSignedOffBy                   InheritedBooleanInfo           `json:"use_signed_off_by,omitempty"`
	CreateNewChangeForAllNotInTarget InheritedBooleanInfo           `json:"create_new_change_for_all_not_in_target,omitempty"`
	RequireChangeID                  InheritedBooleanInfo           `json:"require_change_id,omitempty"`
	EnableSignedPush                 InheritedBooleanInfo           `json:"enable_signed_push,omitempty"`
	MaxObjectSizeLimit               MaxObjectSizeLimitInfo         `json:"max_object_size_limit"`
	SubmitType                       string                         `json:"submit_type"`
	State                            string                         `json:"state,omitempty"`
	Commentlinks                     map[string]string              `json:"commentlinks"`
	Theme                            ThemeInfo                      `json:"theme,omitempty"`
	PluginConfig                     map[string]ConfigParameterInfo `json:"plugin_config,omitempty"`
	Actions                          map[string]ActionInfo          `json:"actions,omitempty"`
}

// ConfigInput entity describes a new project configuration.
type ConfigInput struct {
	Description                      string                       `json:"description,omitempty"`
	UseContributorAgreements         string                       `json:"use_contributor_agreements,omitempty"`
	UseContentMerge                  string                       `json:"use_content_merge,omitempty"`
	UseSignedOffBy                   string                       `json:"use_signed_off_by,omitempty"`
	CreateNewChangeForAllNotInTarget string                       `json:"create_new_change_for_all_not_in_target,omitempty"`
	RequireChangeID                  string                       `json:"require_change_id,omitempty"`
	MaxObjectSizeLimit               MaxObjectSizeLimitInfo       `json:"max_object_size_limit,omitempty"`
	SubmitType                       string                       `json:"submit_type,omitempty"`
	State                            string                       `json:"state,omitempty"`
	PluginConfigValues               map[string]map[string]string `json:"plugin_config_values,omitempty"`
}

// ProjectDescription entity describes a project description.
// In some API functions no JSON will be returned, only strings.
// E.g. GetProjectDescription
type ProjectDescription string

func (d *ProjectDescription) String() string {
	return string(*d)
}

// ProjectOptions specifies the parameters to the ProjectsService.ListProjects.
type ProjectOptions struct {
	// Limit the results to the projects having the specified branch and include the sha1 of the branch in the results.
	Branch string `url:"b,omitempty"`

	// Include project description in the results.
	Description bool `url:"d,omitempty"`

	// Limit the number of projects to be included in the results.
	Limit int `url:"n,omitempty"`

	// Limit the results to those projects that start with the specified prefix.
	Prefix string `url:"p,omitempty"`

	// Limit the results to those projects that match the specified regex.
	// Boundary matchers '^' and '$' are implicit. For example: the regex 'test.*' will match any projects that start with 'test' and regex '.*test' will match any project that end with 'test'.
	Regex string `url:"r,omitempty"`

	// Skip the given number of projects from the beginning of the list.
	Skip string `url:"S,omitempty"`

	// Limit the results to those projects that match the specified substring.
	Substring string `url:"m,omitempty"`

	// Get projects inheritance in a tree-like format.
	// This option does not work together with the branch option.
	Tree string `url:"t,omitempty"`

	// Get projects with specified type: ALL, CODE, PERMISSIONS.
	Type string `url:"type,omitempty"`
}

// ListProjects lists the projects accessible by the caller.
// This is the same as using the ls-projects command over SSH, and accepts the same options as query parameters.
// The entries in the map are sorted by project name.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-projects
func (s *ProjectsService) ListProjects(opt *ProjectOptions) (map[string]ProjectInfo, *Response, error) {
	u := "projects/"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectInfo := new(map[string]ProjectInfo)
	resp, err := s.client.Do(req, projectInfo)
	if err != nil {
		return nil, resp, err
	}

	return *projectInfo, resp, err
}

// GetProject retrieves a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-project
func (s *ProjectsService) GetProject(name string) (*ProjectInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/", name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	projectInfo := new(ProjectInfo)
	resp, err := s.client.Do(req, projectInfo)
	if err != nil {
		return nil, resp, err
	}

	return projectInfo, resp, err
}

// CreateProject creates a new project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-project
func (s *ProjectsService) CreateProject(name string, input *ProjectInput) (*ProjectInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/", name)

	req, err := s.client.NewRequest("PUT", u, input)
	if err != nil {
		return nil, nil, err
	}

	projectInfo := new(ProjectInfo)
	resp, err := s.client.Do(req, projectInfo)
	if err != nil {
		return nil, resp, err
	}

	return projectInfo, resp, err
}

// GetProjectDescription retrieves the description of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-project-description
func (s *ProjectsService) GetProjectDescription(name string) (*string, *Response, error) {
	u := fmt.Sprintf("/projects/%s/description/", name)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	desc := new(ProjectDescription)

	resp, err := s.client.Do(req, desc)
	if err != nil {
		return nil, resp, err
	}
	descStr := desc.String()
	return &descStr, resp, err
}

/**
Missing Project Endpoints
	Set Project Description
	Delete Project Description
	Get Project Parent
	Set Project Parent
	Get HEAD
	Set HEAD
	Get Repository Statistics
	Get Config
	Set Config
	Run GC
	Ban Commit

Missing Branch Endpoints
	List Branches
	Get Branch
	Create Branch
	Delete Branch
	Delete Branches
	Get Content
	Get Reflog

Missing Child Project Endpoints
	List Child Projects
	Get Child Project

Missing Tag Endpoints
	List Tags
	Get Tag

Missing Commit Endpoints
	Get Commit
	Get Content

Missing Dashboard Endpoints
	List Dashboards
	Get Dashboard
	Set Dashboard
	Delete Dashboard
*/
