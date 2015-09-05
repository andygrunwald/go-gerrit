package diffy

import (
	"fmt"
)

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
	Method  string `json:"method,omitempty"`
	Label   string `json:"label,omitempty"`
	Title   string `json:"title,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// AbandonInput entity contains information for abandoning a change.
type AbandonInput struct {
	Message string `json:"message,omitempty"`
}

// AddReviewerResult entity describes the result of adding a reviewer to a change.
type AddReviewerResult struct {
	Reviewers []ReviewerInfo `json:"reviewers,omitempty"`
	Error     string         `json:"error,omitempty"`
	Confirm   bool           `json:"confirm,omitempty"`
}

// ReviewerInfo entity contains information about a reviewer and its votes on a change.
type ReviewerInfo struct {
	AccountInfo
	Approvals string `json:"approvals"`
}

// ApprovalInfo entity contains information about an approval from a user for a label on a change.
type ApprovalInfo struct {
	AccountInfo
	Value int    `json:"value,omitempty"`
	Date  string `json:"date,omitempty"`
}

// ChangeEditInput entity contains information for restoring a path within change edit.
type ChangeEditInput struct {
	RestorePath string `json:"restore_path,omitempty"`
	OldPath     string `json:"old_path,omitempty"`
	NewPath     string `json:"new_path,omitempty"`
}

// ChangeEditMessageInput entity contains information for changing the commit message within a change edit.
type ChangeEditMessageInput struct {
	Message string `json:"message"`
}

// ChangeMessageInfo entity contains information about a message attached to a change.
type ChangeMessageInfo struct {
	ID             string      `json:"id"`
	Author         AccountInfo `json:"author,omitempty"`
	Date           string      `json:"date"`
	Message        string      `json:"message"`
	RevisionNumber int         `json:"_revision_number,omitempty"`
}

// CherryPickInput entity contains information for cherry-picking a change to a new branch.
type CherryPickInput struct {
	Message     string `json:"message"`
	Destination string `json:"destination"`
}

// CommentRange entity describes the range of an inline comment.
type CommentRange struct {
	StartLine      int `json:"start_line"`
	StartCharacter int `json:"start_character"`
	EndLine        int `json:"end_line"`
	EndCharacter   int `json:"end_character"`
}

// DiffFileMetaInfo entity contains meta information about a file diff
type DiffFileMetaInfo struct {
	Name        string        `json:"name"`
	ContentType string        `json:"content_type"`
	Lines       int           `json:"lines"`
	WebLinks    []WebLinkInfo `json:"web_links,omitempty"`
}

// DiffWebLinkInfo entity describes a link on a diff screen to an external site.
type DiffWebLinkInfo struct {
	Name                     string `json:"name"`
	URL                      string `json:"url"`
	ImageURL                 string `json:"image_url"`
	ShowOnSideBySideDiffView bool   `json:"show_on_side_by_side_diff_view"`
	ShowOnUnifiedDiffView    bool   `json:"show_on_unified_diff_view"`
}

// EditFileInfo entity contains additional information of a file within a change edit.
type EditFileInfo struct {
	WebLinks []WebLinkInfo `json:"web_links,omitempty"`
}

// EditInfo entity contains information about a change edit.
type EditInfo struct {
	Commit       CommitInfo           `json:"commit"`
	BaseRevision string               `json:"baseRevision"`
	Fetch        map[string]FetchInfo `json:"fetch"`
	Files        map[string]FileInfo  `json:"files,omitempty"`
}

// FetchInfo entity contains information about how to fetch a patch set via a certain protocol.
type FetchInfo struct {
	URL      string            `json:"url"`
	Ref      string            `json:"ref"`
	Commands map[string]string `json:"commands,omitempty"`
}

// FileInfo entity contains information about a file in a patch set.
type FileInfo struct {
	Status        string `json:"status,omitempty"`
	Binary        bool   `json:"binary,omitempty"`
	OldPath       string `json:"old_path,omitempty"`
	LinesInserted int    `json:"lines_inserted,omitempty"`
	LinesDeleted  int    `json:"lines_deleted,omitempty"`
}

// FixInput entity contains options for fixing commits using the fix change endpoint.
type FixInput struct {
	DeletePatchSetIfCommitMissing bool   `json:"delete_patch_set_if_commit_missing"`
	ExpectMergedAs                string `json:"expect_merged_as"`
}

// GroupBaseInfo entity contains base information about the group.
type GroupBaseInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// IncludedInInfo entity contains information about the branches a change was merged into and tags it was tagged with.
type IncludedInInfo struct {
	Branches []string          `json:"branches"`
	Tags     []string          `json:"tags"`
	External map[string]string `json:"external,omitempty"`
}

// MergeableInfo entity contains information about the mergeability of a change.
type MergeableInfo struct {
	SubmitType    string   `json:"submit_type"`
	Mergeable     bool     `json:"mergeable"`
	MergeableInto []string `json:"mergeable_into,omitempty"`
}

// ProblemInfo entity contains a description of a potential consistency problem with a change.
// These are not related to the code review process, but rather indicate some inconsistency in Gerrit’s database or repository metadata related to the enclosing change.
type ProblemInfo struct {
	Message string `json:"message"`
	Status  string `json:"status,omitempty"`
	Outcome string `json:"outcome,omitempty"`
}

// RebaseInput entity contains information for changing parent when rebasing.
type RebaseInput struct {
	Base string `json:"base,omitempty"`
}

// RelatedChangesInfo entity contains information about related changes.
type RelatedChangesInfo struct {
	Changes []RelatedChangeAndCommitInfo `json:"changes"`
}

// RestoreInput entity contains information for restoring a change.
type RestoreInput struct {
	Message string `json:"message,omitempty"`
}

// RevertInput entity contains information for reverting a change.
type RevertInput struct {
	Message string `json:"message,omitempty"`
}

// ReviewInfo entity contains information about a review.
type ReviewInfo struct {
	Labels map[string]int `json:"labels"`
}

// TopicInput entity contains information for setting a topic.
type TopicInput struct {
	Topic string `json:"topic,omitempty"`
}

// SubmitRecord entity describes results from a submit_rule.
type SubmitRecord struct {
	Status       string                            `json:"status"`
	Ok           map[string]map[string]AccountInfo `json:"ok,omitempty"`
	Reject       map[string]map[string]AccountInfo `json:"reject,omitempty"`
	Need         map[string]interface{}            `json:"need,omitempty"`
	May          map[string]map[string]AccountInfo `json:"may,omitempty"`
	Impossible   map[string]interface{}            `json:"impossible,omitempty"`
	ErrorMessage string                            `json:"error_message,omitempty"`
}

// SubmitInput entity contains information for submitting a change.
type SubmitInput struct {
	WaitForMerge bool `json:"wait_for_merge"`
}

// SubmitInfo entity contains information about the change status after submitting.
type SubmitInfo struct {
	Status     string `json:"status"`
	OnBehalfOf string `json:"on_behalf_of,omitempty"`
}

// RuleInput entity contains information to test a Prolog rule.
type RuleInput struct {
	Rule    string `json:"rule"`
	Filters string `json:"filters,omitempty"`
}

// ReviewerInput entity contains information for adding a reviewer to a change.
type ReviewerInput struct {
	Reviewer  string `json:"reviewer"`
	Confirmed bool   `json:"confirmed,omitempty"`
}

// ReviewInput entity contains information for adding a review to a revision.
type ReviewInput struct {
	Message      string                    `json:"message,omitempty"`
	Labels       map[string]string         `json:"labels,omitempty"`
	Comments     map[string][]CommentInput `json:"comments,omitempty"`
	StrictLabels bool                      `json:"strict_labels,omitempty"`
	Drafts       string                    `json:"drafts,omitempty"`
	Notify       string                    `json:"notify,omitempty"`
	OnBehalfOf   string                    `json:"on_behalf_of,omitempty"`
}

// RelatedChangeAndCommitInfo entity contains information about a related change and commit.
type RelatedChangeAndCommitInfo struct {
	ChangeID              string     `json:"change_id,omitempty"`
	Commit                CommitInfo `json:"commit"`
	ChangeNumber          int        `json:"_change_number,omitempty"`
	RevisionNumber        int        `json:"_revision_number,omitempty"`
	CurrentRevisionNumber int        `json:"_current_revision_number,omitempty"`
	Status                string     `json:"status,omitempty"`
}

// DiffContent entity contains information about the content differences in a file.
type DiffContent struct {
	A      string            `json:"a,omitempty"`
	B      string            `json:"b,omitempty"`
	AB     string            `json:"ab,omitempty"`
	EditA  DiffIntralineInfo `json:"edit_a,omitempty"`
	EditB  DiffIntralineInfo `json:"edit_b,omitempty"`
	Skip   int               `json:"skip,omitempty"`
	Common bool              `json:"common,omitempty"`
}

// CommitInfo entity contains information about a commit.
type CommitInfo struct {
	Commit    string        `json:"commit,omitempty"`
	Parents   []CommitInfo  `json:"parents"`
	Author    GitPersonInfo `json:"author"`
	Committer GitPersonInfo `json:"committer"`
	Subject   string        `json:"subject"`
	Message   string        `json:"message"`
	WebLinks  []WebLinkInfo `json:"web_links,omitempty"`
}

// CommentInput entity contains information for creating an inline comment.
type CommentInput struct {
	ID        string       `json:"id,omitempty"`
	Path      string       `json:"path,omitempty"`
	Side      string       `json:"side,omitempty"`
	Line      int          `json:"line,omitempty"`
	Range     CommentRange `json:"range,omitempty"`
	InReplyTo string       `json:"in_reply_to,omitempty"`
	Updated   string       `json:"updated,omitempty"`
	Message   string       `json:"message,omitempty"`
}

// DiffIntralineInfo entity contains information about intraline edits in a file.
type DiffIntralineInfo []struct {
	SkipLength int
	MarkLength int
}

// ChangeInfo entity contains information about a change.
type ChangeInfo struct {
	ID                 string                  `json:"id"`
	Project            string                  `json:"project"`
	Branch             string                  `json:"branch"`
	Topic              string                  `json:"topic,omitempty"`
	ChangeID           string                  `json:"change_id"`
	Subject            string                  `json:"subject"`
	Status             string                  `json:"status"`
	Created            string                  `json:"created"`
	Updated            string                  `json:"updated"`
	Starred            bool                    `json:"starred,omitempty"`
	Reviewed           bool                    `json:"reviewed,omitempty"`
	Mergeable          bool                    `json:"mergeable,omitempty"`
	Insertions         int                     `json:"insertions"`
	Deletions          int                     `json:"deletions"`
	Number             int                     `json:"_number"`
	Owner              AccountInfo             `json:"owner"`
	Actions            map[string]ActionInfo   `json:"actions,omitempty"`
	Labels             map[string]LabelInfo    `json:"labels,omitempty"`
	PermittedLabels    map[string][]string     `json:"permitted_labels,omitempty"`
	RemovableReviewers []AccountInfo           `json:"removable_reviewers,omitempty"`
	Messages           []ChangeMessageInfo     `json:"messages,omitempty"`
	CurrentRevision    string                  `json:"current_revision,omitempty"`
	Revisions          map[string]RevisionInfo `json:"revisions,omitempty"`
	MoreChanges        bool                    `json:"_more_changes,omitempty"`
	Problems           []ProblemInfo           `json:"problems,omitempty"`
	BaseChange         string                  `json:"base_change,omitempty"`
}

// LabelInfo entity contains information about a label on a change, always corresponding to the current patch set.
type LabelInfo struct {
	Optional bool `json:"optional,omitempty"`

	// Fields set by LABELS
	Approved     AccountInfo `json:"approved,omitempty"`
	Rejected     AccountInfo `json:"rejected,omitempty"`
	Recommended  AccountInfo `json:"recommended,omitempty"`
	Disliked     AccountInfo `json:"disliked,omitempty"`
	Blocking     bool        `json:"blocking,omitempty"`
	Value        string      `json:"value,omitempty"`
	DefaultValue string      `json:"default_value,omitempty"`

	// Fields set by DETAILED_LABELS
	All    []ApprovalInfo    `json:"all,omitempty"`
	Values map[string]string `json:"values,omitempty"`
}

// RevisionInfo entity contains information about a patch set.
type RevisionInfo struct {
	Draft             bool                  `json:"draft,omitempty"`
	Number            int                   `json:"_number"`
	Created           string                `json:"created"`
	Uploader          AccountInfo           `json:"uploader"`
	Ref               string                `json:"ref"`
	Fetch             map[string]FetchInfo  `json:"fetch"`
	Commit            CommitInfo            `json:"commit,omitempty"`
	Files             map[string]FileInfo   `json:"files,omitempty"`
	Actions           map[string]ActionInfo `json:"actions,omitempty"`
	Reviewed          bool                  `json:"reviewed,omitempty"`
	MessageWithFooter string                `json:"messageWithFooter,omitempty"`
}

// CommentInfo entity contains information about an inline comment.
type CommentInfo struct {
	PatchSet  int          `json:"patch_set,omitempty"`
	ID        string       `json:"id"`
	Path      string       `json:"path,omitempty"`
	Side      string       `json:"side,omitempty"`
	Line      int          `json:"line,omitempty"`
	Range     CommentRange `json:"range,omitempty"`
	InReplyTo string       `json:"in_reply_to,omitempty"`
	Message   string       `json:"message,omitempty"`
	Updated   string       `json:"updated"`
	Author    AccountInfo  `json:"author,omitempty"`
}

// DiffInfo entity contains information about the diff of a file in a revision.
type DiffInfo struct {
	MetaA           DiffFileMetaInfo  `json:"meta_a,omitempty"`
	MetaB           DiffFileMetaInfo  `json:"meta_b,omitempty"`
	ChangeType      string            `json:"change_type"`
	IntralineStatus string            `json:"intraline_status,omitempty"`
	DiffHeader      []string          `json:"diff_header"`
	Content         []DiffContent     `json:"content"`
	WebLinks        []DiffWebLinkInfo `json:"web_links,omitempty"`
	Binary          bool              `json:"binary,omitempty"`
}

// SuggestedReviewerInfo entity contains information about a reviewer that can be added to a change (an account or a group).
type SuggestedReviewerInfo struct {
	Account AccountInfo   `json:"account,omitempty"`
	Group   GroupBaseInfo `json:"group,omitempty"`
}

// QueryChangeOptions specifies the parameters to the ChangesService.QueryChanges.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type QueryChangeOptions struct {
	// Query parameter
	// Clients are allowed to specify more than one query by setting the q parameter multiple times.
	// In this case the result is an array of arrays, one per query in the same order the queries were given in.
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/user-search.html#_search_operators
	Query []string `url:"q,omitempty"`

	// The n parameter can be used to limit the returned results.
	// If the n query parameter is supplied and additional changes exist that match the query beyond the end, the last change object has a _more_changes: true JSON field set.
	Limit int `url:"n,omitempty"`

	// The S or start query parameter can be supplied to skip a number of changes from the list.
	Skip  int `url:"S,omitempty"`
	Start int `url:"start,omitempty"`

	ChangeOptions
}

// ChangeOptions specifies the parameters for Query changes.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type ChangeOptions struct {
	// Additional fields can be obtained by adding o parameters, each option requires more database lookups and slows down the query response time to the client so they are generally disabled by default.
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
	AdditionalFields []string `url:"o,omitempty"`
}

// QueryChanges visible to the caller.
// The query string must be provided by the q parameter.
// The n parameter can be used to limit the returned results.
//
// The change output is sorted by the last update time, most recently updated to oldest updated.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
func (s *ChangesService) QueryChanges(opt *QueryChangeOptions) (*[]ChangeInfo, *Response, error) {
	u := "changes/"

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]ChangeInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetChange retrieves a change.
// Additional fields can be obtained by adding o parameters, each option requires more database lookups and slows down the query response time to the client so they are generally disabled by default.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change
func (s *ChangesService) GetChange(changeID string, opt *ChangeOptions) (*ChangeInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s", changeID)
	return s.getChangeInfoResponse(u, opt)
}

// GetChangeDetail retrieves a change with labels, detailed labels, detailed accounts, and messages.
// Additional fields can be obtained by adding o parameters, each option requires more database lookups and slows down the query response time to the client so they are generally disabled by default.
//
// This response will contain all votes for each label and include one combined vote.
// The combined label vote is calculated in the following order (from highest to lowest): REJECTED > APPROVED > DISLIKED > RECOMMENDED.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-change
func (s *ChangesService) GetChangeDetail(changeID string, opt *ChangeOptions) (*ChangeInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/detail", changeID)
	return s.getChangeInfoResponse(u, opt)
}

// getChangeInfoResponse retrieved a single ChangeInfo Response for a GET request
func (s *ChangesService) getChangeInfoResponse(u string, opt *ChangeOptions) (*ChangeInfo, *Response, error) {
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(ChangeInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetTopic retrieves the topic of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-topic
func (s *ChangesService) GetTopic(changeID string) (*string, *Response, error) {
	u := fmt.Sprintf("changes/%s/topic", changeID)
	return getStringResponseWithoutOptions(s.client, u)
}

/*
Missing Change Endpoints
	Create Change
	Set Topic
	Delete Topic
	Abandon Change
	Restore Change
	Rebase Change
	Revert Change
	Submit Change
	Changes submitted together
	Publish Draft Change
	Delete Draft Change
	Get Included In
	Index Change
	List Change Comments
	List Change Drafts
	Check change
	Fix change

Missing Change Edit Endpoints
	Get Change Edit Details
	Change file content in Change Edit
	Restore file content or rename files in Change Edit
	Change commit message in Change Edit
	Delete file in Change Edit
	Retrieve file content from Change Edit
	Retrieve meta data of a file from Change Edit
	Retrieve commit message from Change Edit or current patch set of the change
	Publish Change Edit
	Rebase Change Edit
	Delete Change Edit

Missing Reviewer Endpoints
	List Reviewers
	Suggest Reviewers
	Get Reviewer
	Add Reviewer
	Delete Reviewer

Missing Revision Endpoints
	Get Commit
	Get Revision Actions
	Get Review
	Get Related Changes
	Set Review
	Rebase Revision
	Submit Revision
	Publish Draft Revision
	Delete Draft Revision
	Get Patch
	Get Mergeable
	Get Submit Type
	Test Submit Type
	Test Submit Rule
	List Revision Drafts
	Create Draft
	Get Draft
	Update Draft
	Delete Draft
	List Revision Comments
	Get Comment
	List Files
	Get Content
	Get Diff
	Set Reviewed
	Delete Reviewed
	Cherry Pick Revision
*/
