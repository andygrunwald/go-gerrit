package gerrit

import (
	"context"
	"fmt"
	"net/url"
)

// GetCommit retrieves a commit of a project.
// The commit must be visible to the caller.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-commit
func (s *ProjectsService) GetCommit(ctx context.Context, projectName, commitID string) (*CommitInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/commits/%s", url.QueryEscape(projectName), commitID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// GetIncludeIn Retrieves the branches and tags in which a change is included.
// Branches that are not visible to the calling user according to the project’s read permissions are filtered out from the result.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-included-in
func (s *ProjectsService) GetIncludeIn(ctx context.Context, projectName, commitID string) (*IncludedInInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/commits/%s/in", url.QueryEscape((projectName)), commitID)
	req, err := s.client.NewRequest(ctx, "GET", u, nil)
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

// GetCommitContent gets the content of a file from a certain commit.
// The content is returned as base64 encoded string.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html##get-content-from-commit
func (s *ProjectsService) GetCommitContent(ctx context.Context, projectName, commitID, fileID string) (string, *Response, error) {
	u := fmt.Sprintf("projects/%s/commits/%s/files/%s/content", url.QueryEscape(projectName), commitID, fileID)
	return getStringResponseWithoutOptions(ctx, s.client, u)
}
