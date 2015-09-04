package diffy

// ChangesService contains Change related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html
type ChangesService struct {
	client *Client
}

// WebLinkInfo entity describes a link to an external site.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#web-link-info
type WebLinkInfo struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	ImageURL string `json:"image_url"`
}

// GitPersonInfo entity contains information about the author/committer of a commit.
type GitPersonInfo struct {
	Name  string `json:"name"`
	EMail string `json:"email"`
	Date  string `json:"date"`
	TZ    int    `json:"tz"`
}

// ActionInfo entity describes a REST API call the client can make to manipulate a resource.
// These are frequently implemented by plugins and may be discovered at runtime.
type ActionInfo struct {
	Method  string `json:"method"`
	Label   string `json:"label"`
	Title   string `json:"title"`
	Enabled bool   `json:"enabled"`
}

/*
AbandonInput
AddReviewerResult
ApprovalInfo
ChangeEditInput
ChangeEditMessageInput
ChangeInfo
ChangeMessageInfo
CherryPickInput
CommentInfo
CommentInput
CommentRange
CommitInfo
DiffContent
DiffFileMetaInfo
DiffInfo
DiffIntralineInfo
DiffWebLinkInfo
EditFileInfo
EditInfo
FetchInfo
FileInfo
FixInput
GroupBaseInfo
IncludedInInfo
LabelInfo
MergeableInfo
ProblemInfo
RebaseInput
RelatedChangeAndCommitInfo
RelatedChangesInfo
RestoreInput
RevertInput
ReviewInfo
ReviewInput
ReviewerInfo
ReviewerInput
RevisionInfo
RuleInput
SubmitInfo
SubmitInput
SubmitRecord
SuggestedReviewerInfo
TopicInput
*/
