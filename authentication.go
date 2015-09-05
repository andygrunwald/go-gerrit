package diffy

const (
	// HTTP Basic Authentication
	authTypeBasic = 1
	// HTTP Digest Authentication
	authTypeDigest = 2
	// HTTP Cookie Authentication
	authTypeCookie = 3
)

// TODO Digest auth

// AuthenticationService contains Authentication related functions.
//
// Gerrit API docs: https://gerrit-review.googlesource.com/Documentation/rest-api.html#authentication
type AuthenticationService struct {
	client *Client

	// Storage for authentication
	// Username or name of cookie
	name string
	// Password or value of cookie
	secret   string
	authType int
}

// SetBasicAuth sets basic parameters for HTTP Basic auth
func (s *AuthenticationService) SetBasicAuth(username, password string) {
	s.name = username
	s.secret = password
	s.authType = authTypeBasic
}

// SetCookieAuth sets basic parameters for HTTP Cookie
func (s *AuthenticationService) SetCookieAuth(name, value string) {
	s.name = name
	s.secret = value
	s.authType = authTypeCookie
}

// HasBasicAuth checks if the auth type is HTTP Basic auth
func (s *AuthenticationService) HasBasicAuth() bool {
	return s.authType == authTypeBasic
}

// HasCookieAuth checks if the auth type is HTTP Cookie based
func (s *AuthenticationService) HasCookieAuth() bool {
	return s.authType == authTypeCookie
}

// HasAuth checks if an auth type is used
func (s *AuthenticationService) HasAuth() bool {
	return s.authType > 0
}

// ResetAuth resets all former authentification settings
func (s *AuthenticationService) ResetAuth() {
	s.name = ""
	s.secret = ""
	s.authType = 0
}
