package diffy

// ConfigService contains Config related REST endpoints
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api-config.html
type ConfigService struct {
	client *Client
}
