package diffy

import (
	"fmt"
)

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

// RelatedChangesInfo entity contains information about related changes.
type RelatedChangesInfo struct {
	Changes []RelatedChangeAndCommitInfo `json:"changes"`
}

// FileInfo entity contains information about a file in a patch set.
type FileInfo struct {
	Status        string `json:"status,omitempty"`
	Binary        bool   `json:"binary,omitempty"`
	OldPath       string `json:"old_path,omitempty"`
	LinesInserted int    `json:"lines_inserted,omitempty"`
	LinesDeleted  int    `json:"lines_deleted,omitempty"`
}

// ActionInfo entity describes a REST API call the client can make to manipulate a resource.
// These are frequently implemented by plugins and may be discovered at runtime.
type ActionInfo struct {
	Method  string `json:"method,omitempty"`
	Label   string `json:"label,omitempty"`
	Title   string `json:"title,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
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

// MergeableInfo entity contains information about the mergeability of a change.
type MergeableInfo struct {
	SubmitType    string   `json:"submit_type"`
	Mergeable     bool     `json:"mergeable"`
	MergeableInto []string `json:"mergeable_into,omitempty"`
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

// CommitOptions specifies the parameters for GetCommit call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-commit
type CommitOptions struct {
	// Adding query parameter links (for example /changes/.../commit?links) returns a CommitInfo with the additional field web_links.
	Weblinks bool `url:"links,omitempty"`
}

// MergableOptions specifies the parameters for GetMergable call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-mergeable
type MergableOptions struct {
	// If the other-branches parameter is specified, the mergeability will also be checked for all other branches.
	OtherBranches bool `url:"other-branches,omitempty"`
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

// GetRelatedChanges retrieves related changes of a revision.
// Related changes are changes that either depend on, or are dependencies of the revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-related-changes
func (s *ChangesService) GetRelatedChanges(changeID, revisionID string) (*RelatedChangesInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/related", changeID, revisionID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(RelatedChangesInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetDraft retrieves a draft comment of a revision that belongs to the calling user.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-draft
func (s *ChangesService) GetDraft(changeID, revisionID, draftID string) (*CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/%s", changeID, revisionID, draftID)
	return s.getCommentInfoResponse(u)
}

// GetComment retrieves a published comment of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-comment
func (s *ChangesService) GetComment(changeID, revisionID, commentID string) (*CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s//comments/%s", changeID, revisionID, commentID)
	return s.getCommentInfoResponse(u)
}

// GetSubmitType gets the method the server will use to submit (merge) the change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-submit-type
func (s *ChangesService) GetSubmitType(changeID, revisionID string) (*string, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/submit_type", changeID, revisionID)
	return getStringResponseWithoutOptions(s.client, u)
}

// GetRevisionActions retrieves revision actions of the revision of a change.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-revision-actions
func (s *ChangesService) GetRevisionActions(changeID, revisionID string) (*map[string]ActionInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/actions", changeID, revisionID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(map[string]ActionInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetCommit retrieves a parsed commit of a revision.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-commit
func (s *ChangesService) GetCommit(changeID, revisionID string, opt *CommitOptions) (*CommitInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/commit", changeID, revisionID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(CommitInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetReview retrieves a review of a revision.
//
// As response a ChangeInfo entity with detailed labels and detailed accounts is returned that describes the review of the revision.
// The revision for which the review is retrieved is contained in the revisions field.
// In addition the current_revision field is set if the revision for which the review is retrieved is the current revision of the change.
// Please note that the returned labels are always for the current patch set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-review
func (s *ChangesService) GetReview(changeID, revisionID string) (*ChangeInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/review", changeID, revisionID)
	return s.getChangeInfoResponse(u, nil)
}

// GetMergeable gets the method the server will use to submit (merge) the change and an indicator if the change is currently mergeable.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-mergeable
func (s *ChangesService) GetMergeable(changeID, revisionID string, opt *MergableOptions) (*MergeableInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/mergeable", changeID, revisionID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(MergeableInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// ListRevisionDrafts lists the draft comments of a revision that belong to the calling user.
// Returns a map of file paths to lists of CommentInfo entries.
// The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-drafts
func (s *ChangesService) ListRevisionDrafts(changeID, revisionID string) (*map[string][]CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/drafts/", changeID, revisionID)
	return s.getCommentInfoMapSliceResponse(u)
}

// ListRevisionComments lists the published comments of a revision.
// As result a map is returned that maps the file path to a list of CommentInfo entries.
// The entries in the map are sorted by file path and only include file (or inline) comments.
// Use the Get Change Detail endpoint to retrieve the general change message (or comment).
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-comments
func (s *ChangesService) ListRevisionComments(changeID, revisionID string) (*map[string][]CommentInfo, *Response, error) {
	u := fmt.Sprintf("changes/%s/revisions/%s/comments/", changeID, revisionID)
	return s.getCommentInfoMapSliceResponse(u)
}

// ListFiles lists the files that were modified, added or deleted in a revision.
// As result a map is returned that maps the file path to a list of FileInfo entries.
// The entries in the map are sorted by file path.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-files
func (s *ChangesService) ListFiles(changeID, revisionID string) (*map[string]FileInfo, *Response, error) {
	// TODO: Missing q parameter
	// The request parameter q changes the response to return a list of all files (modified or unmodified) that contain that substring in the path name. This is useful to implement suggestion services finding a file by partial name.
	u := fmt.Sprintf("changes/%s/revisions/%s/files/", changeID, revisionID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(map[string]FileInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// ListFilesReviewed lists the files that were modified, added or deleted in a revision.
// The difference between ListFiles and ListFilesReviewed is that the caller has marked these files as reviewed.
// Clients that also need the FileInfo should make two requests.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#list-files
func (s *ChangesService) ListFilesReviewed(changeID, revisionID string) (*[]FileInfo, *Response, error) {
	// TODO: Missing q parameter
	// The request parameter q changes the response to return a list of all files (modified or unmodified) that contain that substring in the path name. This is useful to implement suggestion services finding a file by partial name.
	u := fmt.Sprintf("changes/%s/revisions/%s/files/", changeID, revisionID)

	opt := struct {
		// The request parameter reviewed changes the response to return a list of the paths the caller has marked as reviewed.
		Reviewed bool `url:"reviewed,omitempty"`
	}{
		Reviewed: true,
	}

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]FileInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

/*
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
