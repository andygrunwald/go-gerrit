package gerrit

import (
	"context"
	"fmt"
	"net/url"
)

type LabelDefinitionInfo struct {
	Name               string            `json:"name,omitempty"`
	Description        string            `json:"description,omitempty"`
	ProjectName        string            `json:"project_name,omitempty"`
	Function           string            `json:"function,omitempty"`
	Values             map[string]string `json:"values,omitempty"`
	DefaultValue       int               `json:"default_value,omitempty"`
	Branches           []string          `json:"branches,omitempty"`
	CanOverride        bool              `json:"can_override,omitempty"`
	CopyCondition      string            `json:"copy_condition,omitempty"`
	AllowPostSubmit    bool              `json:"allow_post_submit,omitempty"`
	IgnoreSelfApproval bool              `json:"ignore_self_approval,omitempty"`
}

type LabelDefinitionInput struct {
	CommitMessage      string            `json:"commit_message,omitempty"`
	Name               string            `json:"name,omitempty"`
	Description        string            `json:"descriptioan,omitempty"`
	Function           string            `json:"function,omitempty"`
	Values             map[string]string `json:"values,omitempty"`
	DefaultValue       int               `json:"default_value,omitempty"`
	Branches           []string          `json:"branches,omitempty"`
	CanOverride        bool              `json:"can_override,omitempty"`
	CopyCondition      string            `json:"copy_condition,omitempty"`
	UnsetCopyCondition bool              `json:"unset_copy_condition,omitempty"`
	AllowPostSubmit    bool              `json:"allow_post_submit,omitempty"`
	IgnoreSelfApproval bool              `json:"ignore_self_approval,omitempty"`
}

type DeleteLabelInput struct {
	CommitMessage string `json:"commit_message,omitempty"`
}

type BatchLabelInput struct {
	CommitMessage string                 `json:"commit_message,omitempty"`
	Delete        []string               `json:"delete,omitempty"`
	Create        []LabelDefinitionInput `json:"create,omitempty"`
	Update        []LabelDefinitionInput `json:"update,omitempty"`
}

// ListLabels lists the labels for a project
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#list-labels
func (s *ProjectsService) ListLabels(ctx context.Context, projectName string) (*[]LabelDefinitionInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/labels/", url.QueryEscape(projectName))

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]LabelDefinitionInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetLabel gets the definition of a label associated with a project
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-label
func (s *ProjectsService) GetLabel(ctx context.Context, projectName, labelName string) (*LabelDefinitionInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/labels/%s", url.QueryEscape(projectName), url.QueryEscape(labelName))

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(LabelDefinitionInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// CreateLabel creates a label type for an associated project
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-label
func (s *ProjectsService) CreateLabel(ctx context.Context, projectName string, input *LabelDefinitionInput) (*LabelDefinitionInfo, *Response, error) {
	// this uses SetLabel internally because it appears to use the same method
	// on the same URL as SetLabel
	return s.SetLabel(ctx, projectName, input)
}

// SetLabel sets the definition of a label for an associated project
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-label
func (s *ProjectsService) SetLabel(ctx context.Context, projectName string, input *LabelDefinitionInput) (*LabelDefinitionInfo, *Response, error) {
	u := fmt.Sprintf("projects/%s/labels/%s", url.QueryEscape(projectName), url.QueryEscape(input.Name))

	req, err := s.client.NewRequest(ctx, "PUT", u, input)
	if err != nil {
		return nil, nil, err
	}

	v := new(LabelDefinitionInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// DeleteLabel deletes a label definition for an associated project
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#delete-label
func (s *ProjectsService) DeleteLabel(ctx context.Context, projectName, labelName string, input *DeleteLabelInput) (*Response, error) {
	u := fmt.Sprintf("projects/%s/labels/%s", url.QueryEscape(projectName), url.QueryEscape(labelName))
	return s.client.DeleteRequest(ctx, u, input)
}

// Batch update labels
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#batch-update-labels
func (s *ProjectsService) BatchUpdateLabels(ctx context.Context, projectName string, input *BatchLabelInput) (*Response, error) {
	u := fmt.Sprintf("projects/%s/labels/", url.QueryEscape(projectName))

	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Create labels change for review
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-labels-change
func (s *ProjectsService) CreateLabelsChange(ctx context.Context, projectName string, input *BatchLabelInput) (*Response, error) {
	u := fmt.Sprintf("projects/%s/labels:review", url.QueryEscape(projectName))

	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
