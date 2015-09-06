package diffy

import (
	"fmt"
)

// ReviewerInfo entity contains information about a reviewer and its votes on a change.
type ReviewerInfo struct {
	AccountInfo
	Approvals string `json:"approvals"`
}

// SuggestedReviewerInfo entity contains information about a reviewer that can be added to a change (an account or a group).
type SuggestedReviewerInfo struct {
	Account AccountInfo   `json:"account,omitempty"`
	Group   GroupBaseInfo `json:"group,omitempty"`
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

/*
Missing Reviewer Endpoints
	Add Reviewer
	Delete Reviewer
*/
