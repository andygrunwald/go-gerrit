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

/*
AbandonInput
ActionInfo
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
GitPersonInfo
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
