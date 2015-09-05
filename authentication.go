package diffy

import ()

const (
	// HTTP Basic Authentication
	authTypeBasic = 1
	// HTTP Digest Authentication
	authTypeDigest = 2
)

// TODO Digest auth

// AuthenticationService contains Authentication related functions.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api.html#authentication
type AuthenticationService struct {
	client *Client

	// Storage for authentication
	username string
	password string
	authType int
}

// SetBasicAuth sets basic parameters for HTTP Basic auth
func (s *AuthenticationService) SetBasicAuth(username, password string) {
	s.username = username
	s.password = password
	s.authType = authTypeBasic
}

// HasBasicAuth checks if the auth type is HTTP Basic auth
func (s *AuthenticationService) HasBasicAuth() bool {
	return s.authType == authTypeBasic
}

// HasAuth checks if an auth type is used
func (s *AuthenticationService) HasAuth() bool {
	return s.authType > 0 && len(s.username) > 0 && len(s.password) > 0
}
