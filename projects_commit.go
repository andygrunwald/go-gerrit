package diffy

import (
	"fmt"
)

// GetCommit retrieves a commit of a project.
// The commit must be visible to the caller.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-commit
func (s *ProjectsService) GetCommit(projectName, commitID string) (*CommitInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/commits/%s", projectName, commitID)

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

/**
Missing Commit Endpoints
	Get Content
*/
