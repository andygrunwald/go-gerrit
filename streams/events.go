package streams

// Event represents events sent via the `gerrit stream-events` command, or sent as a webhook.
// It differs slightly from the REST API types, which is why this lives in a separate package.
// Docs: https://gerrit-review.googlesource.com/Documentation/cmd-stream-events.html#events
type Event struct {
	Type string `json:"type"`

	Change   Change   `json:"change"`
	PatchSet PatchSet `json:"patchSet"`

	Abandoner Account `json:"abandoner"`
	Author    Account `json:"author"`
	Adder     Account `json:"adder"`
	Changer   Account `json:"changer"`
	Deleter   Account `json:"deleter"`
	Editor    Account `json:"editor"`
	Restorer  Account `json:"restorer"`
	Remover   Account `json:"remover"`
	Reviewer  Account `json:"reviewer"`
	Submitter Account `json:"submitter"`
	Uploader  Account `json:"uploader"`

	Reason         string     `json:"reason,omitempty"`
	EventCreatedOn int        `json:"eventCreatedOn,omitempty"`
	NewRev         string     `json:"newRev,omitempty"`
	Approvals      []Approval `json:"approvals,omitempty"`
	Comment        string     `json:"comment,omitempty"`

	Added    []string `json:"added,omitempty"`
	Removed  []string `json:"removed,omitempty"`
	Hashtags []string `json:"hashtags,omitempty"`

	ProjectName string `json:"projectName,omitempty"`
	ProjectHead string `json:"projectHead,omitempty"`

	RefUpdate  RefUpdate   `json:"refUpdate"`
	RefUpdates []RefUpdate `json:"refUpdates,omitempty"`

	OldTopic string `json:"oldTopic,omitempty"`

	OldHead string `json:"oldHead,omitempty"`
	NewHead string `json:"newHead,omitempty"`
}
