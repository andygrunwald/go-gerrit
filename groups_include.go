package diffy

import (
	"fmt"
)

// ListIncludedGroups lists the directly included groups of a group.
// The entries in the list are sorted by group name and UUID.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#included-groups
func (s *GroupsService) ListIncludedGroups(groupID string) (*[]GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups/", groupID)

	req, err := s.client.NewRequest("GET", u, nil)
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
func (s *GroupsService) GetIncludedGroup(groupID, includeGroupID string) (*GroupInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/groups/%s", groupID, includeGroupID)

	req, err := s.client.NewRequest("GET", u, nil)
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

/*
Missing Group Include Endpoints
	Include Group
	Include Groups
	Delete Included Group
	Delete Included Groups
*/
