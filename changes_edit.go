package diffy

import (
	"fmt"
)

// EditInfo entity contains information about a change edit.
type EditInfo struct {
	Commit       CommitInfo           `json:"commit"`
	BaseRevision string               `json:"baseRevision"`
	Fetch        map[string]FetchInfo `json:"fetch"`
	Files        map[string]FileInfo  `json:"files,omitempty"`
}

// EditFileInfo entity contains additional information of a file within a change edit.
type EditFileInfo struct {
	WebLinks []WebLinkInfo `json:"web_links,omitempty"`
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

/*
Missing Change Edit Endpoints
	Change file content in Change Edit
	Restore file content or rename files in Change Edit
	Change commit message in Change Edit
	Delete file in Change Edit
	Retrieve file content from Change Edit
	Publish Change Edit
	Rebase Change Edit
	Delete Change Edit
*/
