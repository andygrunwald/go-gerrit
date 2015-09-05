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

// AbandonInput entity contains information for abandoning a change.
type AbandonInput struct {
	Message string `json:"message"`
}

// AddReviewerResult entity describes the result of adding a reviewer to a change.
type AddReviewerResult struct {
	Reviewers []ReviewerInfo `json:"reviewers"`
	Error     string         `json:"error"`
	Confirm   bool           `json:"confirm"`
}

// ReviewerInfo entity contains information about a reviewer and its votes on a change.
type ReviewerInfo struct {
	AccountInfo
	Approvals string `json:"approvals"`
}

// ApprovalInfo entity contains information about an approval from a user for a label on a change.
type ApprovalInfo struct {
	AccountInfo
	Value int    `json:"value"`
	Date  string `json:"date"`
}

// ChangeEditInput entity contains information for restoring a path within change edit.
type ChangeEditInput struct {
	RestorePath string `json:"restore_path"`
	OldPath     string `json:"old_path"`
	NewPath     string `json:"new_path"`
}

// ChangeEditMessageInput entity contains information for changing the commit message within a change edit.
type ChangeEditMessageInput struct {
	Message string `json:"message"`
}

// ChangeMessageInfo entity contains information about a message attached to a change.
type ChangeMessageInfo struct {
	ID             string      `json:"id"`
	Author         AccountInfo `json:"author"`
	Date           string      `json:"date"`
	Message        string      `json:"message"`
	RevisionNumber int         `json:"_revision_number"`
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
	WebLinks    []WebLinkInfo `json:"web_links"`
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
	WebLinks []WebLinkInfo `json:"web_links"`
}

// EditInfo entity contains information about a change edit.
type EditInfo struct {
	Commit       CommitInfo           `json:"commit"`
	BaseRevision string               `json:"baseRevision"`
	Fetch        map[string]FetchInfo `json:"fetch"`
	Files        map[string]FileInfo  `json:"files"`
}

// FetchInfo entity contains information about how to fetch a patch set via a certain protocol.
type FetchInfo struct {
	URL      string            `json:"url"`
	Ref      string            `json:"ref"`
	Commands map[string]string `json:"commands"`
}

// FileInfo entity contains information about a file in a patch set.
type FileInfo struct {
	Status        string `json:"status"`
	Binary        bool   `json:"binary"`
	OldPath       string `json:"old_path"`
	LinesInserted int    `json:"lines_inserted"`
	LinesDeleted  int    `json:"lines_deleted"`
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
	External map[string]string `json:"external"`
}

// MergeableInfo entity contains information about the mergeability of a change.
type MergeableInfo struct {
	SubmitType    string   `json:"submit_type"`
	Mergeable     bool     `json:"mergeable"`
	MergeableInto []string `json:"mergeable_into"`
}

// ProblemInfo entity contains a description of a potential consistency problem with a change.
// These are not related to the code review process, but rather indicate some inconsistency in Gerrit’s database or repository metadata related to the enclosing change.
type ProblemInfo struct {
	Message string `json:"message"`
	Status  string `json:"status"`
	Outcome string `json:"outcome"`
}

// RebaseInput entity contains information for changing parent when rebasing.
type RebaseInput struct {
	Base string `json:"base"`
}

// RelatedChangesInfo entity contains information about related changes.
type RelatedChangesInfo struct {
	Changes []RelatedChangeAndCommitInfo `json:"changes"`
}

// RestoreInput entity contains information for restoring a change.
type RestoreInput struct {
	Message string `json:"message"`
}

// RevertInput entity contains information for reverting a change.
type RevertInput struct {
	Message string `json:"message"`
}

// ReviewInfo entity contains information about a review.
type ReviewInfo struct {
	Labels map[string]int `json:"labels"`
}

// TopicInput entity contains information for setting a topic.
type TopicInput struct {
	Topic string `json:"topic"`
}

// SubmitRecord entity describes results from a submit_rule.
type SubmitRecord struct {
	Status       string                            `json:"status"`
	Ok           map[string]map[string]AccountInfo `json:"ok"`
	Reject       map[string]map[string]AccountInfo `json:"reject"`
	Need         map[string]interface{}            `json:"need"`
	May          map[string]map[string]AccountInfo `json:"may"`
	Impossible   map[string]interface{}            `json:"impossible"`
	ErrorMessage string                            `json:"error_message"`
}

// SubmitInput entity contains information for submitting a change.
type SubmitInput struct {
	WaitForMerge bool `json:"wait_for_merge"`
}

// SubmitInfo entity contains information about the change status after submitting.
type SubmitInfo struct {
	Status     string `json:"status"`
	OnBehalfOf string `json:"on_behalf_of"`
}

// RuleInput entity contains information to test a Prolog rule.
type RuleInput struct {
	Rule    string `json:"rule"`
	Filters string `json:"filters"`
}

// ReviewerInput entity contains information for adding a reviewer to a change.
type ReviewerInput struct {
	Reviewer  string `json:"reviewer"`
	Confirmed bool   `json:"confirmed"`
}

// ReviewInput entity contains information for adding a review to a revision.
type ReviewInput struct {
	Message      string                    `json:"message"`
	Labels       map[string]string         `json:"labels"`
	Comments     map[string][]CommentInput `json:"comments"`
	StrictLabels bool                      `json:"strict_labels"`
	Drafts       string                    `json:"drafts"`
	Notify       string                    `json:"notify"`
	OnBehalfOf   string                    `json:"on_behalf_of"`
}

// RelatedChangeAndCommitInfo entity contains information about a related change and commit.
type RelatedChangeAndCommitInfo struct {
	ChangeID              string     `json:"change_id"`
	Commit                CommitInfo `json:"commit"`
	ChangeNumber          int        `json:"_change_number"`
	RevisionNumber        int        `json:"_revision_number"`
	CurrentRevisionNumber int        `json:"_current_revision_number"`
	Status                string     `json:"status"`
}

// DiffContent entity contains information about the content differences in a file.
type DiffContent struct {
	A      string            `json:"a"`
	B      string            `json:"b"`
	AB     string            `json:"ab"`
	EditA  DiffIntralineInfo `json:"edit_a,omitempty"`
	EditB  DiffIntralineInfo `json:"edit_b,omitempty"`
	Skip   int               `json:"skip"`
	Common bool              `json:"common"`
}

// CommitInfo entity contains information about a commit.
type CommitInfo struct {
	Commit    string        `json:"commit"`
	Parents   []CommitInfo  `json:"parents"`
	Author    GitPersonInfo `json:"author"`
	Committer GitPersonInfo `json:"committer"`
	Subject   string        `json:"subject"`
	Message   string        `json:"message"`
	WebLinks  []WebLinkInfo `json:"web_links"`
}

// CommentInput entity contains information for creating an inline comment.
type CommentInput struct {
	ID        string       `json:"id"`
	Path      string       `json:"path"`
	Side      string       `json:"side"`
	Line      int          `json:"line"`
	Range     CommentRange `json:"range"`
	InReplyTo string       `json:"in_reply_to"`
	Updated   string       `json:"updated"`
	Message   string       `json:"message"`
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
	Topic              string                  `json:"topic"`
	ChangeID           string                  `json:"change_id"`
	Subject            string                  `json:"subject"`
	Status             string                  `json:"status"`
	Created            string                  `json:"created"`
	Updated            string                  `json:"updated"`
	Starred            bool                    `json:"starred"`
	Reviewed           bool                    `json:"reviewed"`
	Mergeable          bool                    `json:"mergeable"`
	Insertions         int                     `json:"insertions"`
	Deletions          int                     `json:"deletions"`
	Number             int                     `json:"_number"`
	Owner              AccountInfo             `json:"owner"`
	Actions            map[string]ActionInfo   `json:"actions"`
	Labels             map[string]LabelInfo    `json:"labels"`
	PermittedLabels    map[string][]string     `json:"permitted_labels"`
	RemovableReviewers []AccountInfo           `json:"removable_reviewers"`
	Messages           []ChangeMessageInfo     `json:"messages"`
	CurrentRevision    string                  `json:"current_revision"`
	Revisions          map[string]RevisionInfo `json:"revisions"`
	MoreChanges        bool                    `json:"_more_changes"`
	Problems           []ProblemInfo           `json:"problems"`
	BaseChange         string                  `json:"base_change"`
}

// LabelInfo entity contains information about a label on a change, always corresponding to the current patch set.
type LabelInfo struct {
	Optional bool `json:"optional"`

	// Fields set by LABELS
	Approved     AccountInfo `json:"approved"`
	Rejected     AccountInfo `json:"rejected"`
	Recommended  AccountInfo `json:"recommended"`
	Disliked     AccountInfo `json:"disliked"`
	Blocking     bool        `json:"blocking"`
	Value        string      `json:"value"`
	DefaultValue string      `json:"default_value"`

	// Fields set by DETAILED_LABELS
	All    []ApprovalInfo    `json:"all"`
	Values map[string]string `json:"values"`
}

// RevisionInfo entity contains information about a patch set.
type RevisionInfo struct {
	Draft             bool                  `json:"draft"`
	Number            int                   `json:"_number"`
	Created           string                `json:"created"`
	Uploader          AccountInfo           `json:"uploader"`
	Ref               string                `json:"ref"`
	Fetch             map[string]FetchInfo  `json:"fetch"`
	Commit            CommitInfo            `json:"commit"`
	Files             map[string]FileInfo   `json:"files"`
	Actions           map[string]ActionInfo `json:"actions"`
	Reviewed          bool                  `json:"reviewed"`
	MessageWithFooter string                `json:"messageWithFooter"`
}

// CommentInfo entity contains information about an inline comment.
type CommentInfo struct {
	PatchSet  int          `json:"patch_set"`
	ID        string       `json:"id"`
	Path      string       `json:"path"`
	Side      string       `json:"side"`
	Line      int          `json:"line"`
	Range     CommentRange `json:"range"`
	InReplyTo string       `json:"in_reply_to"`
	Message   string       `json:"message"`
	Updated   string       `json:"updated"`
	Author    AccountInfo  `json:"author"`
}

// DiffInfo entity contains information about the diff of a file in a revision.
type DiffInfo struct {
	MetaA           DiffFileMetaInfo  `json:"meta_a"`
	MetaB           DiffFileMetaInfo  `json:"meta_b"`
	ChangeType      string            `json:"change_type"`
	IntralineStatus string            `json:"intraline_status"`
	DiffHeader      []string          `json:"diff_header"`
	Content         []DiffContent     `json:"content"`
	WebLinks        []DiffWebLinkInfo `json:"web_links"`
	Binary          bool              `json:"binary"`
}

// SuggestedReviewerInfo entity contains information about a reviewer that can be added to a change (an account or a group).
type SuggestedReviewerInfo struct {
	Account AccountInfo   `json:"account,omitempty"`
	Group   GroupBaseInfo `json:"group,omitempty"`
}
