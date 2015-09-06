package diffy

import (
	"fmt"
)

// ListGroupMembersOptions specifies the different options for the ListGroupMembers call.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#group-members
type ListGroupMembersOptions struct {
	// To resolve the included groups of a group recursively and to list all members the parameter recursive can be set.
	// Members from included external groups and from included groups which are not visible to the calling user are ignored.
	Recursive bool `url:"recursive,omitempty"`
}

// ListGroupMembers lists the direct members of a Gerrit internal group.
// The entries in the list are sorted by full name, preferred email and id.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#group-members
func (s *GroupsService) ListGroupMembers(groupID string, opt *ListGroupMembersOptions) (*[]AccountInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/members/", groupID)

	u, err := addOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new([]AccountInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

// GetGroupMember retrieves a group member.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html#get-group-member
func (s *GroupsService) GetGroupMember(groupID, accountID string) (*AccountInfo, *Response, error) {
	u := fmt.Sprintf("groups/%s/members/%s", groupID, accountID)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(AccountInfo)
	resp, err := s.client.Do(req, v)
	if err != nil {
		return nil, resp, err
	}

	return v, resp, err
}

/*
Missing Group Member Endpoints
	Add Group Member
	Add Group Members
	Delete Group Member
	Delete Group Members
*/
