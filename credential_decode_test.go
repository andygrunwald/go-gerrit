package gerrit

import (
	"net/url"
	"testing"
)

// TestCredentialPercentDecoding verifies that percent-escaped credentials are
// decoded consistently, whether or not the URL contains an explicit port.
// URLs with a port take the ReParseURL path inside NewClient(); URLs without
// one are handled by url.Parse. Both must yield the same credentials.
//
// See https://github.com/andygrunwald/go-gerrit/issues/224
func TestCredentialPercentDecoding(t *testing.T) {
	const (
		rawWithPort    = "http://ad%6Din:se%2Fcret@localhost:8080/"
		rawWithoutPort = "http://ad%6Din:se%2Fcret@gerrit.example.com/"
	)

	withPort, err := credentialsFromRegex(rawWithPort)
	if err != nil {
		t.Fatalf("regex path: %v", err)
	}
	withoutPort, err := credentialsFromURLParse(rawWithoutPort)
	if err != nil {
		t.Fatalf("url.Parse path: %v", err)
	}

	if withPort != withoutPort {
		t.Fatalf("inconsistent decoding: %+v vs %+v", withPort, withoutPort)
	}
	if withPort.username != "admin" || withPort.password != "se/cret" {
		t.Fatalf("unexpected credentials: %+v", withPort)
	}
}

type credentials struct {
	username string
	password string
}

// credentialsFromRegex mirrors the ReParseURL branch of NewClient().
func credentialsFromRegex(rawurl string) (credentials, error) {
	m := ReParseURL.FindAllStringSubmatch(rawurl, -1)
	if len(m) == 0 || len(m[0]) <= 5 {
		return credentials{}, ErrNoInstanceGiven
	}
	username, err := url.PathUnescape(m[0][2])
	if err != nil {
		return credentials{}, err
	}
	password, err := url.PathUnescape(m[0][3])
	if err != nil {
		return credentials{}, err
	}
	return credentials{username, password}, nil
}

// credentialsFromURLParse mirrors the url.Parse branch of NewClient().
func credentialsFromURLParse(rawurl string) (credentials, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return credentials{}, err
	}
	password, _ := u.User.Password()
	return credentials{u.User.Username(), password}, nil
}
