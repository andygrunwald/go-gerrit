package gerrit_test

import (
	"context"
	"crypto/md5" // HTTP Digest's MD5 challenge is part of the test protocol.
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestDigestAuthenticationUsesRequestURLAndQuery(t *testing.T) {
	t.Parallel()

	for _, basePath := range []string{"/", "/r/", "/review/"} {
		for _, method := range []string{http.MethodGet, http.MethodPost} {
			t.Run(method+basePath, func(t *testing.T) {
				t.Parallel()
				const username = "fixture-user"
				const password = "fixture-password"
				const realm = "fixture-realm"
				const nonce = "fixture-nonce"
				const query = "q=status%3Aopen&n=2"
				wantURI := basePath + "a/changes/?" + query
				var requests atomic.Int32
				md5Hex := func(value string) string {
					return fmt.Sprintf("%x", md5.Sum([]byte(value)))
				}
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					requests.Add(1)
					if r.Method != method || r.RequestURI != wantURI {
						t.Errorf("request = %s %s, want %s %s", r.Method, r.RequestURI, method, wantURI)
						w.WriteHeader(http.StatusNotFound)
						return
					}
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Error(err)
						w.WriteHeader(http.StatusBadRequest)
						return
					}
					header := r.Header.Get("Authorization")
					if header == "" {
						if len(body) != 0 {
							t.Errorf("challenge sent the request body: %q", body)
						}
						// Omit uri so the client must use the complete request URI.
						w.Header().Set("WWW-Authenticate", fmt.Sprintf(
							`Digest realm="%s", nonce="%s", qop="auth", algorithm=MD5`,
							realm, nonce,
						))
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					fields := make(map[string]string)
					for _, field := range strings.Split(strings.TrimPrefix(header, "Digest "), ",") {
						key, value, ok := strings.Cut(strings.TrimSpace(field), "=")
						if ok {
							fields[key] = strings.Trim(value, `"`)
						}
					}
					ha1 := md5Hex(username + ":" + realm + ":" + password)
					ha2 := md5Hex(method + ":" + wantURI)
					wantResponse := md5Hex(strings.Join([]string{
						ha1, nonce, fields["nc"], fields["cnonce"], "auth", ha2,
					}, ":"))
					if !strings.HasPrefix(header, "Digest ") ||
						fields["username"] != username || fields["realm"] != realm ||
						fields["nonce"] != nonce || fields["uri"] != wantURI ||
						fields["qop"] != "auth" || fields["nc"] != "00000001" ||
						fields["cnonce"] == "" || fields["response"] != wantResponse {
						t.Errorf("invalid Digest header: %s", header)
						w.WriteHeader(http.StatusUnauthorized)
						return
					}
					wantBody := ""
					if method == http.MethodPost {
						wantBody = "{\"message\":\"fixture\"}\n"
					}
					if string(body) != wantBody {
						t.Errorf("authenticated request body = %q, want %q", body, wantBody)
					}
					_, err = io.WriteString(w, "{}")
					if err != nil {
						t.Error(err)
					}
				}))
				t.Cleanup(server.Close)
				client, err := gerrit.NewClient(context.Background(), server.URL+basePath, server.Client())
				if err != nil {
					t.Fatal(err)
				}
				client.Authentication.SetDigestAuth(username, password)
				var body interface{}
				if method == http.MethodPost {
					body = map[string]string{"message": "fixture"}
				}
				req, err := client.NewRequest(context.Background(), method, "changes/?"+query, body)
				if err != nil {
					t.Fatal(err)
				}
				var response map[string]interface{}
				_, err = client.Do(req, &response)
				if err != nil {
					t.Fatal(err)
				}
				if requests.Load() != 2 {
					t.Errorf("requests = %d, want 2", requests.Load())
				}
			})
		}
	}
}
