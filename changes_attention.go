package gerrit

import (
	"context"
	"fmt"
)

// AttentionSetInfo entity contains details of users that are in the attention set.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#attention-set-info
type AttentionSetInfo struct {
	// AccountInfo entity.
	Account AccountInfo `json:"account"`
	// The timestamp of the last update.
	LastUpdate Timestamp `json:"last_update"`
	// The reason of for adding or removing the user.
	Reason string `json:"reason"`
}

// Doc: https://gerrit-review.googlesource.com/Documentation/user-notify.html#recipient-types
type RecipientType string

// AttentionSetInput entity contains details for adding users to the attention
// set and removing them from it.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#attention-set-input
type AttentionSetInput struct {
	User          string                       `json:"user,omitempty"`
	Reason        string                       `json:"reason"`
	Notify        string                       `json:"notify,omitempty"`
	NotifyDetails map[RecipientType]NotifyInfo `json:"notify_details,omitempty"`
}

// RemoveAttention deletes a single user from the attention set of a change.
// AttentionSetInput.Input must be provided
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#remove-from-attention-set
func (s *ChangesService) RemoveAttention(ctx context.Context, changeID, accountID string, input *AttentionSetInput) (*Response, error) {
	u := fmt.Sprintf("changes/%s/attention/%s", changeID, accountID)

	return s.client.DeleteRequest(ctx, u, input)
}
