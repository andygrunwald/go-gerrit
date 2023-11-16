package gerrit

import (
	"context"
	"fmt"
)

// ListIncludedGroups lists the directly included groups of a group.
// The entries in the list are sorted by group name and UUID.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#included-groups
func (s *GroupsService) ListIncludedGroups(ctx context.Context, groupID string) (*[]GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups/", groupID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]GroupInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetIncludedGroup retrieves an included group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-included-group
func (s *GroupsService) GetIncludedGroup(ctx context.Context, groupID, includeGroupID string) (*GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups/%s", groupID, includeGroupID)

	req, err := s.client.NewRequest(ctx, "GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(GroupInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// IncludeGroup includes an internal or external group into a Gerrit internal group.
// External groups must be specified using the UUID.
//
// As response a GroupInfo entity is returned that describes the included group.
// The request also succeeds if the group is already included in this group, but then the HTTP response code is 200 OK.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#include-group
func (s *GroupsService) IncludeGroup(ctx context.Context, groupID, includeGroupID string) (*GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups/%s", groupID, includeGroupID)

	req, err := s.client.NewRequest(ctx, "PUT", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(GroupInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// IncludeGroups includes one or several groups into a Gerrit internal group.
// The groups to be included into the group must be provided in the request body as a GroupsInput entity.
//
// As response a list of GroupInfo entities is returned that describes the groups that were specified in the GroupsInput.
// A GroupInfo entity is returned for each group specified in the input, independently of whether the group was newly included into the group or whether the group was already included in the group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#include-groups
func (s *GroupsService) IncludeGroups(ctx context.Context, groupID string, input *GroupsInput) (*[]GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups", groupID)

	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return nil, nil, err
	}

	v := new([]GroupInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// DeleteIncludedGroup deletes an included group from a Gerrit internal group.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#include-group
func (s *GroupsService) DeleteIncludedGroup(ctx context.Context, groupID, includeGroupID string) (*Response, error) {
	u := fmt.Sprintf("groups/%s/groups/%s", groupID, includeGroupID)
	return s.client.DeleteRequest(ctx, u, nil)
}

// DeleteIncludedGroups delete one or several included groups from a Gerrit internal group.
// The groups to be deleted from the group must be provided in the request body as a GroupsInput entity.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#delete-included-groups
func (s *GroupsService) DeleteIncludedGroups(ctx context.Context, groupID string, input *GroupsInput) (*Response, error) {
	u := fmt.Sprintf("groups/%s/groups.delete", groupID)

	req, err := s.client.NewRequest(ctx, "POST", u, input)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
