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
}

/*
Structs to create:
BanInput
BanResultInfo
BranchInfo
BranchInput
ConfigInfo
ConfigInput
ConfigParameterInfo
DashboardInfo
DashboardInput
DashboardSectionInfo
DeleteBranchesInput
GCInput
HeadInput
InheritedBooleanInfo
MaxObjectSizeLimitInfo
ProjectDescriptionInput
ProjectInfo
ProjectInput
ProjectParentInput
ReflogEntryInfo
RepositoryStatisticsInfo
TagInfo
ThemeInfo
*/
