package gerrit

import "fmt"

// GetHashtags gets the hashtags associated with a change.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#get-hashtags
func (c *ChangesService) GetHashtags(changeID string) ([]string, *Response, error) {
	u := fmt.Sprintf("changes/%s/hashtags", changeID)

	req, err := c.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var hashtags []string
	resp, err := c.client.Do(req, &hashtags)

	return hashtags, resp, err
}

// HashtagsInput entity contains information about hashtags to add to, and/or remove from, a change.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#hashtags-input
type HashtagsInput struct {
	// The list of hashtags to be added to the change.
	Add []string `json:"add,omitempty"`

	// The list of hashtags to be removed from the change.
	Remove []string `json:"remove,omitempty"`
}

// SetHashtags adds and/or removes hashtags from a change.
//
// As response the change’s hashtags are returned as a list of strings.
//
// https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#set-hashtags
func (c *ChangesService) SetHashtags(changeID string, input *HashtagsInput) ([]string, *Response, error) {
	u := fmt.Sprintf("changes/%s/hashtags", changeID)

	req, err := c.client.NewRequest("POST", u, input)
	if err != nil {
		return nil, nil, err
	}

	var hashtags []string
	resp, err := c.client.Do(req, &hashtags)

	return hashtags, resp, err
}
