package gerrit

import (
	"context"
	"fmt"
	"net/url"
)

// TagInfo entity contains information about a tag.
type TagInfo struct {
	Ref      string        `json:"ref"`
	Revision string        `json:"revision"`
	Object   string        `json:"object"`
	Message  string        `json:"message"`
	Tagger   GitPersonInfo `json:"tagger"`
	Created  *Timestamp    `json:"created,omitempty"`
}

// TagInput entity for create a tag.
type TagInput struct {
	Ref      string `json:"ref"`
	Revision string `json:"revision,omitempty"`
	Message  string `json:"message,omitempty"`
}

// DeleteTagsInput entity for delete tags.
type DeleteTagsInput struct {
	Tags []string `json:"tags"`
}

// ListTags list the tags of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-tags
func (s *ProjectsService) ListTags(ctx context.Context, projectName string, opt *ProjectBaseOptions) (*[]TagInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/tags/", url.QueryEscape(projectName))
	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]TagInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetTag retrieves a tag of a project.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-tag
func (s *ProjectsService) GetTag(ctx context.Context, projectName, tagName string) (*TagInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(projectName), url.QueryEscape(tagName))

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(TagInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// CreateTag create a tag of a project
//
// Gerrit API docs:https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-tag
func (s *ProjectsService) CreateTag(ctx context.Context, projectName, tagName string, input *TagInput) (*TagInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(projectName), url.QueryEscape(tagName))

	req, err := s.client.NewRequest(ctx, "PUT", u, input)
	if err != nil {
		return nil, nil, err
	}

	v := new(TagInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// DeleteTag delete a tag of a project
//
// Gerrit API docs:https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-tag
func (s *ProjectsService) DeleteTag(ctx context.Context, projectName, tagName string) (*Response, error) {
	u := fmt.Sprintf("projects/%s/tags/%s", url.QueryEscape(projectName), url.QueryEscape(tagName))

	req, err := s.client.NewRequest(ctx, "DELETE", u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}

// DeleteTags delete tags of a project
//
// Gerrit API docs:https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-tags
func (s *ProjectsService) DeleteTags(ctx context.Context, projectName string, input *DeleteTagsInput) (*Response, error) {
	u := fmt.Sprintf("projects/%s/tags:delete", url.QueryEscape(projectName))

	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)

	return resp, err
}
