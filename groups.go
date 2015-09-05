package diffy

// GroupsService contains Group related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-groups.html
type GroupsService struct {
	client *Client
}

// GroupAuditEventInfo entity contains information about an audit event of a group.
type GroupAuditEventInfo struct {
	// TODO Member AccountInfo OR GroupInfo `json:"member"`
	Type string      `json:"type"`
	User AccountInfo `json:"user"`
	Date string      `json:"date"`
}

// GroupInfo entity contains information about a group.
// This can be a Gerrit internal group, or an external group that is known to Gerrit.
type GroupInfo struct {
	ID          string           `json:"id"`
	Name        string           `json:"name,omitempty"`
	URL         string           `json:"url,omitempty"`
	Options     GroupOptionsInfo `json:"options"`
	Description string           `json:"description,omitempty"`
	GroupID     int              `json:"group_id,omitempty"`
	Owner       string           `json:"owner,omitempty"`
	OwnerID     string           `json:"owner_id,omitempty"`
	Members     []AccountInfo    `json:"members,omitempty"`
	Includes    []GroupInfo      `json:"includes,omitempty"`
}

// GroupInput entity contains information for the creation of a new internal group.
type GroupInput struct {
	Name         string `json:"name,omitempty"`
	Description  string `json:"description,omitempty"`
	VisibleToAll bool   `json:"visible_to_all,omitempty"`
	OwnerID      string `json:"owner_id,omitempty"`
}

// GroupOptionsInfo entity contains options of the group.
type GroupOptionsInfo struct {
	VisibleToAll bool `json:"visible_to_all,omitempty"`
}

// GroupOptionsInput entity contains new options for a group.
type GroupOptionsInput struct {
	VisibleToAll bool `json:"visible_to_all,omitempty"`
}

// GroupsInput entity contains information about groups that should be included into a group or that should be deleted from a group.
type GroupsInput struct {
	OneGroup string   `json:"_one_group,omitempty"`
	Groups   []string `json:"groups,omitempty"`
}

// MembersInput entity contains information about accounts that should be added as members to a group or that should be deleted from the group
type MembersInput struct {
	OneMember string   `json:"_one_member,omitempty"`
	Members   []string `json:"members,omitempty"`
}
