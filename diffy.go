package diffy

// A Client manages communication with the Gerrit API.
type Client struct {
}

// NewClient returns a new Gerrit API client.
func NewClient() *Client {
	// TODO Use http client?
	// Like https://github.com/google/go-github/blob/master/github/github.go#L128
	instance := &Client{}

	return instance
}
