package gerrit

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/url"
	"time"
)

// PatchSet contains detailed information about a specific patch set.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/json.html#patchSet
type PatchSet struct {
	Number    string      `json:"number"`
	Revision  string      `json:"revision"`
	Parents   []string    `json:"parents"`
	Ref       string      `json:"ref"`
	Uploader  AccountInfo `json:"uploader"`
	Author    AccountInfo `json:"author"`
	CreatedOn int         `json:"createdOn"`
	IsDraft   bool        `json:"isDraft"`
	Kind      string      `json:"kind"`
}

// RefUpdate contains data about a reference update.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/json.html#refUpdate
type RefUpdate struct {
	OldRev  string `json:"oldRev"`
	NewRev  string `json:"newRev"`
	RefName string `json:"refName"`
	Project string `json:"project"`
}

// EventInfo contains information about an event emitted by Gerrit.  This
// structure can be used either when parsing streamed events or when reading
// the output of the events-log plugin.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/cmd-stream-events.html#events
type EventInfo struct {
	Type           string        `json:"type"`
	Change         ChangeInfo    `json:"change,omitempty"`
	PatchSet       PatchSet      `json:"patchSet,omitempty"`
	EventCreatedOn int           `json:"eventCreatedOn,omitempty"`
	Reason         string        `json:"reason,omitempty"`
	Abandoner      AccountInfo   `json:"abandoner,omitempty"`
	Restorer       AccountInfo   `json:"restorer,omitempty"`
	Submitter      AccountInfo   `json:"submitter,omitempty"`
	Author         AccountInfo   `json:"author,omitempty"`
	Uploader       AccountInfo   `json:"uploader,omitempty"`
	Approvals      []AccountInfo `json:"approvals,omitempty"`
	Comment        string        `json:"comment,omitempty"`
	Editor         AccountInfo   `json:"editor,omitempty"`
	Added          []string      `json:"added,omitempty"`
	Removed        []string      `json:"removed,omitempty"`
	Hashtags       []string      `json:"hashtags,omitempty"`
	RefUpdate      RefUpdate     `json:"refUpdate,omitempty"`
	Project        string        `json:"project,omitempty"`
	Reviewer       AccountInfo   `json:"reviewer,omitempty"`
	OldTopic       string        `json:"oldTopic,omitempty"`
	Changer        AccountInfo   `json:"changer,omitempty"`
}

// EventsLogService contains functions for querying the API provided
// by the optional events-log plugin.
type EventsLogService struct {
	client *Client
}

// EventsLogOptions contains options for querying events from the events-logs
// plugin.
type EventsLogOptions struct {
	From time.Time
	To   time.Time
}

// getURL returns the url that should be used in the request.  This will vary
// depending on the options provided to GetEvents.
func (events *EventsLogService) getURL(options *EventsLogOptions) (string, error) {
	parsed, err := url.Parse("/plugins/events-log/events/")
	if err != nil {
		return "", err
	}

	query := parsed.Query()

	if !options.From.IsZero() {
		query.Set("t1", options.From.Format("2006-01-02 15:04:05"))
	}

	if !options.To.IsZero() {
		query.Set("t2", options.To.Format("2006-01-02 15:04:05"))
	}

	return parsed.String(), nil
}

// GetEvents returns a list of events for the given input options.  Use of this
// function an authenticated user.
//
// Gerrit API docs: https://<yourserver>/plugins/events-log/Documentation/rest-api-events.html
func (events *EventsLogService) GetEvents(options *EventsLogOptions) (*[]EventInfo, *Response, error) {
	requestURL, err := events.getURL(options)
	if err != nil {
		return nil, nil, err
	}

	request, err := events.client.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, nil, err
	}

	// Perform the request but do not pass in a structure to unpack
	// the response into.  The format of the response is one EventInfo
	// object per line so we need to manually handle the response here.
	response, err := events.client.Do(request, nil)
	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(response.Body)

	defer response.Body.Close()
	if err != nil {
		return nil, nil, err
	}

	eventInfo := new([]EventInfo)
	for _, line := range bytes.Split(body, []byte("\n")) {
		if len(line) > 0 {
			event := EventInfo{}
			err := json.Unmarshal(line, &event)
			if err != nil {
				return nil, nil, err
			}
			*eventInfo = append(*eventInfo, event)
		}
	}

	return eventInfo, response, err
}
