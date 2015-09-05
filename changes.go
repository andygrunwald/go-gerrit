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

// QueryOptions specifies global parameters to query changes / reviewers.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type QueryOptions struct {
	// Query parameter
	// Clients are allowed to specify more than one query by setting the q parameter multiple times.
	// In this case the result is an array of arrays, one per query in the same order the queries were given in.
	//
	// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/user-search.html#_search_operators
	Query []string `url:"q,omitempty"`

	// The n parameter can be used to limit the returned results.
	// If the n query parameter is supplied and additional changes exist that match the query beyond the end, the last change object has a _more_changes: true JSON field set.
	Limit int `url:"n,omitempty"`
}

// QueryChangeOptions specifies the parameters to the ChangesService.QueryChanges.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-changes
type QueryChangeOptions struct {
	QueryOptions

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

// ChangeEditDetailOptions specifies the parameters to the ChangesService.GetChangeEditDetails.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-detail
type ChangeEditDetailOptions struct {
	// When request parameter list is provided the response also includes the file list.
	List bool `url:"list,omitempty"`
	// When base request parameter is provided the file list is computed against this base revision.
	Base bool `url:"base,omitempty"`
	// When request parameter download-commands is provided fetch info map is also included.
	DownloadCommands bool `url:"download-commands,omitempty"`
}

// DiffOptions specifies the parameters for GetDiff call.
type DiffOptions struct {
	// If the intraline parameter is specified, intraline differences are included in the diff.
	Intraline bool `url:"intraline,omitempty"`

	//The base parameter can be specified to control the base patch set from which the diff should be generated.
	Base bool `url:"base,omitempty"`

	// If the weblinks-only parameter is specified, only the diff web links are returned.
	WeblinksOnly bool `url:"weblinks-only,omitempty"`

	// The ignore-whitespace parameter can be specified to control how whitespace differences are reported in the result. Valid values are NONE, TRAILING, CHANGED or ALL.
	IgnoreWhitespace string `url:"ignore-whitespace,omitempty"`

	// The context parameter can be specified to control the number of lines of surrounding context in the diff.
	// Valid values are ALL or number of lines.
	Context string `url:"context,omitempty"`
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

// ChangesSubmittedTogether returns a list of all changes which are submitted when {submit} is called for this change, including the current change itself.
// An empty list is returned if this change will be submitted by itself (no other changes).
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#submitted_together
func (s *ChangesService) ChangesSubmittedTogether(changeID string) (*[]ChangeInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/submitted_together", changeID)

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

// GetIncludedIn retrieves the branches and tags in which a change is included.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-included-in
func (s *ChangesService) GetIncludedIn(changeID string) (*IncludedInInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/in", changeID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(IncludedInInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// ListChangeComments lists the published comments of all revisions of the change.
// The entries in the map are sorted by file path, and the comments for each path are sorted by patch set number.
// Each comment has the patch_set and author fields set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-comments
func (s *ChangesService) ListChangeComments(changeID string) (*map[string]CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/in", changeID)
	return s.getCommentInfoMapResponse(u)
}

// ListChangeDrafts lLists the draft comments of all revisions of the change that belong to the calling user.
// The entries in the map are sorted by file path, and the comments for each path are sorted by patch set number.
// Each comment has the patch_set field set, and no author.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-change-drafts
func (s *ChangesService) ListChangeDrafts(changeID string) (*map[string]CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/drafts", changeID)
	return s.getCommentInfoMapResponse(u)
}

// getCommentInfoMapResponse retrieved a map of CommentInfo Response for a GET request
func (s *ChangesService) getCommentInfoMapResponse(u string) (*map[string]CommentInfo, *Response, error) {
	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(map[string]CommentInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// CheckChange performs consistency checks on the change, and returns a ChangeInfo entity with the problems field set to a list of ProblemInfo entities.
// Depending on the type of problem, some fields not marked optional may be missing from the result.
// At least id, project, branch, and _number will be present.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#check-change
func (s *ChangesService) CheckChange(changeID string) (*ChangeInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/check", changeID)
	return s.getChangeInfoResponse(u, nil)
}

// GetChangeEditDetails retrieves a change edit details.
// As response an EditInfo entity is returned that describes the change edit, or “204 No Content” when change edit doesn’t exist for this change.
// Change edits are stored on special branches and there can be max one edit per user per change.
// Edits aren’t tracked in the database.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-detail
func (s *ChangesService) GetChangeEditDetails(changeID string, opt *ChangeEditDetailOptions) (*EditInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/edit", changeID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(EditInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// RetrieveMetaDataOfAFileFromChangeEdit retrieves meta data of a file from a change edit.
// Currently only web links are returned.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-meta-data
func (s *ChangesService) RetrieveMetaDataOfAFileFromChangeEdit(changeID, filePath string) (*EditFileInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/edit/%s/meta", changeID, filePath)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(EditFileInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// RetrieveCommitMessageFromChangeEdit retrieves commit message from change edit.
// The commit message is returned as base64 encoded string.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-edit-message
func (s *ChangesService) RetrieveCommitMessageFromChangeEdit(changeID string) (*string, *Response, error) {
	u := fmt.Sprintf("changes/%s/edit:message", changeID)
	return getStringResponseWithoutOptions(s.client, u)
}

// ListReviewers lists the reviewers of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-reviewers
func (s *ChangesService) ListReviewers(changeID string) (*[]ReviewerInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/reviewers/", changeID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]ReviewerInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// SuggestReviewers suggest the reviewers for a given query q and result limit n.
// If result limit is not passed, then the default 10 is used.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#suggest-reviewers
func (s *ChangesService) SuggestReviewers(changeID string, opt *QueryOptions) (*[]SuggestedReviewerInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/suggest_reviewers", changeID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]SuggestedReviewerInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetReviewer retrieves a reviewer of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-reviewer
func (s *ChangesService) GetReviewer(changeID, accountID string) (*ReviewerInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/reviewers/%s", changeID, accountID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(ReviewerInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetDiff gets the diff of a file from a certain revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-diff
func (s *ChangesService) GetDiff(changeID, revisionID, fileID string, opt *DiffOptions) (*DiffInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/files/%s/diff", changeID, revisionID, fileID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(DiffInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

/*
Get calls
	Get Commit
	Get Revision Actions
	Get Review
	Get Related Changes
	Get Mergeable
	Get Submit Type
	List Revision Drafts
	Get Draft
	List Revision Comments
	Get Comment
	List Files

Missing Change Endpoints
	Create Change
	Set Topic
	Delete Topic
	Abandon Change
	Restore Change
	Rebase Change
	Revert Change
	Submit Change
	Publish Draft Change
	Delete Draft Change
	Index Change
	Fix change

Missing Change Edit Endpoints
	Change file content in Change Edit
	Restore file content or rename files in Change Edit
	Change commit message in Change Edit
	Delete file in Change Edit
	Retrieve file content from Change Edit
	Publish Change Edit
	Rebase Change Edit
	Delete Change Edit

Missing Reviewer Endpoints
	Add Reviewer
	Delete Reviewer

Missing Revision Endpoints
	Set Review
	Rebase Revision
	Submit Revision
	Publish Draft Revision
	Delete Draft Revision
	Get Patch
	Test Submit Type
	Test Submit Rule
	Create Draft
	Update Draft
	Delete Draft
	Get Content
	Set Reviewed
	Delete Reviewed
	Cherry Pick Revision
*/
