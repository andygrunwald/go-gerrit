package gerrit_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

const (
	// testGerritInstanceURL is a test instance url that won`t be called
	testGerritInstanceURL = "https://go-review.googlesource.com/"
)

var (
	// testMux is the HTTP request multiplexer used with the test server.
	testMux *http.ServeMux

	// testClient is the gerrit client being tested.
	testClient *gerrit.Client

	// testServer is a test HTTP server used to provide mock API responses.
	testServer *httptest.Server
)

type testValues map[string]string

// setup sets up a test HTTP server along with a gerrit.Client that is configured to talk to that test server.
// Tests should register handlers on mux which provide mock responses for the API method being tested.
func setup() {
	// Test server
	testMux = http.NewServeMux()
	testServer = httptest.NewServer(testMux)

	// gerrit client configured to use test server
	testClient, _ = gerrit.NewClient(testServer.URL, nil)
}

// teardown closes the test HTTP server.
func teardown() {
	testServer.Close()
}

// makedigestheader takes the incoming request and produces a string
// which can be used for the WWW-Authenticate header.
func makedigestheader(request *http.Request) string {
	return fmt.Sprintf(
		`Digest realm="Gerrit Code Review", domain="http://%s/", qop="auth", nonce="fakevaluefortesting"`,
		request.Host)
}

// writeresponse writes the requested value to the provided response writer and sets
// the http code
func writeresponse(t *testing.T, writer http.ResponseWriter, value interface{}, code int) {
	writer.WriteHeader(code)

	unmarshalled, err := json.Marshal(value)
	if err != nil {
		t.Error(err.Error())
		return
	}

	data := []byte(`)]}'` + "\n" + string(unmarshalled))
	if _, err := writer.Write(data); err != nil {
		t.Error(err.Error())
		return
	}
}

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) { // nolint: unparam
	if got := r.URL.String(); got != want {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}

func testFormValues(t *testing.T, r *http.Request, values testValues) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	if err := r.ParseForm(); err != nil {
		t.Error(err)
	}
	if got := r.Form; !reflect.DeepEqual(got, want) {
		t.Errorf("Request parameters: %v, want %v", got, want)
	}
}

func TestNewClient_NoGerritInstance(t *testing.T) {
	mockData := []string{"", "://not-existing"}
	for _, data := range mockData {
		c, err := gerrit.NewClient(data, nil)
		if c != nil {
			t.Errorf("NewClient return is not nil. Expected no client. Go %+v", c)
		}
		if err == nil {
			t.Error("No error occured by empty Gerrit Instance. Expected one.")
		}
	}
}

func TestNewClient_Services(t *testing.T) {
	c, err := gerrit.NewClient("https://gerrit-review.googlesource.com/", nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}

	if c.Authentication == nil {
		t.Error("No AuthenticationService found.")
	}
	if c.Access == nil {
		t.Error("No AccessService found.")
	}
	if c.Accounts == nil {
		t.Error("No AccountsService found.")
	}
	if c.Changes == nil {
		t.Error("No ChangesService found.")
	}
	if c.Config == nil {
		t.Error("No ConfigService found.")
	}
	if c.Groups == nil {
		t.Error("No GroupsService found.")
	}
	if c.Plugins == nil {
		t.Error("No PluginsService found.")
	}
	if c.Projects == nil {
		t.Error("No ProjectsService found.")
	}
}

func TestNewClient_TestErrNoInstanceGiven(t *testing.T) {
	_, err := gerrit.NewClient("", nil)
	if err != gerrit.ErrNoInstanceGiven {
		t.Error("Expected `ErrNoInstanceGiven`")
	}
}

func TestNewClient_NoCredentials(t *testing.T) {
	client, err := gerrit.NewClient("http://localhost/", nil)
	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}
	if client.Authentication.HasAuth() {
		t.Error("Expected HasAuth() to return false")
	}
}

func TestNewClient_UsernameWithoutPassword(t *testing.T) {
	_, err := gerrit.NewClient("http://foo@localhost/", nil)
	if err != gerrit.ErrUserProvidedWithoutPassword {
		t.Error("Expected ErrUserProvidedWithoutPassword")
	}
}

func TestNewClient_AuthenticationFailed(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/a/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		writeresponse(t, w, nil, http.StatusUnauthorized)
	})

	serverURL := fmt.Sprintf("http://admin:secret@%s/", testServer.Listener.Addr().String())
	client, err := gerrit.NewClient(serverURL, nil)
	if err != gerrit.ErrAuthenticationFailed {
		t.Error(err)
	}
	if client.Authentication.HasAuth() {
		t.Error("Expected HasAuth() == false")
	}
}

func TestNewClient_DigestAuth(t *testing.T) {
	setup()
	defer teardown()

	account := gerrit.AccountInfo{
		AccountID: 100000,
		Name:      "test",
		Email:     "test@localhost",
		Username:  "test"}
	hits := 0

	testMux.HandleFunc("/a/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch hits {
		case 1:
			w.Header().Set("WWW-Authenticate", makedigestheader(r))
			writeresponse(t, w, nil, http.StatusUnauthorized)
		case 2:
			// go-gerrit should set Authorization in response to a `WWW-Authenticate` header
			if !strings.Contains(r.Header.Get("Authorization"), `username="admin"`) {
				t.Error(`Missing username="admin"`)
			}
			writeresponse(t, w, account, http.StatusOK)
		case 3:
			t.Error("Did not expect another request")
		}
	})

	serverURL := fmt.Sprintf("http://admin:secret@%s/", testServer.Listener.Addr().String())
	client, err := gerrit.NewClient(serverURL, nil)
	if err != nil {
		t.Error(err)
	}
	if !client.Authentication.HasDigestAuth() {
		t.Error("Expected HasDigestAuth() == true")
	}
}

func TestNewClient_BasicAuth(t *testing.T) {
	setup()
	defer teardown()

	account := gerrit.AccountInfo{
		AccountID: 100000,
		Name:      "test",
		Email:     "test@localhost",
		Username:  "test"}
	hits := 0

	testMux.HandleFunc("/a/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch hits {
		case 1:
			writeresponse(t, w, nil, http.StatusUnauthorized)
		case 2:
			// The second request should be a basic auth request if the first request, which is for
			// digest based auth, fails.
			if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
				t.Error("Missing 'Basic ' prefix")
			}
			writeresponse(t, w, account, http.StatusOK)
		case 3:
			t.Error("Did not expect another request")
		}
	})

	serverURL := fmt.Sprintf("http://admin:secret@%s/", testServer.Listener.Addr().String())
	client, err := gerrit.NewClient(serverURL, nil)
	if err != nil {
		t.Error(err)
	}
	if !client.Authentication.HasBasicAuth() {
		t.Error("Expected HasBasicAuth() == true")
	}
}

func TestNewClient_ReParseURL(t *testing.T) {
	urls := map[string][]string{
		"http://admin:ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg@127.0.0.1:5000/": {
			"http://127.0.0.1:5000/", "admin", "ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg",
		},
		"http://admin:ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg@127.0.0.1:5000/foo": {
			"http://127.0.0.1:5000/foo", "admin", "ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg",
		},
		"http://admin:ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg@127.0.0.1:5000": {
			"http://127.0.0.1:5000", "admin", "ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg",
		},
		"https://admin:foo/bar@localhost:5": {
			"https://localhost:5", "admin", "foo/bar",
		},
	}
	for input, expectations := range urls {
		submatches := gerrit.ReParseURL.FindAllStringSubmatch(input, -1)
		submatch := submatches[0]
		username := submatch[2]
		password := submatch[3]
		endpoint := fmt.Sprintf(
			"%s://%s:%s%s", submatch[1], submatch[4], submatch[5], submatch[6])
		if endpoint != expectations[0] {
			t.Errorf("%s != %s", expectations[0], endpoint)
		}
		if username != expectations[1] {
			t.Errorf("%s != %s", expectations[1], username)
		}
		if password != expectations[2] {
			t.Errorf("%s != %s", expectations[2], password)
		}

	}
}

func TestNewClient_BasicAuth_PasswordWithSlashes(t *testing.T) {
	setup()
	defer teardown()

	account := gerrit.AccountInfo{
		AccountID: 100000,
		Name:      "test",
		Email:     "test@localhost",
		Username:  "test"}
	hits := 0

	testMux.HandleFunc("/a/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch hits {
		case 1:
			writeresponse(t, w, nil, http.StatusUnauthorized)
		case 2:
			// The second request should be a basic auth request if the first request, which is for
			// digest based auth, fails.
			if !strings.HasPrefix(r.Header.Get("Authorization"), "Basic ") {
				t.Error("Missing 'Basic ' prefix")
			}
			writeresponse(t, w, account, http.StatusOK)
		case 3:
			t.Error("Did not expect another request")
		}
	})

	serverURL := fmt.Sprintf(
		"http://admin:ZOSOKjgV/kgEkN0bzPJp+oGeJLqpXykqWFJpon/Ckg@%s",
		testServer.Listener.Addr().String())
	client, err := gerrit.NewClient(serverURL, nil)
	if err != nil {
		t.Error(err)
	}
	if !client.Authentication.HasAuth() {
		t.Error("Expected HasAuth() == true")
	}
}

func TestNewClient_CookieAuth(t *testing.T) {
	setup()
	defer teardown()

	account := gerrit.AccountInfo{
		AccountID: 100000,
		Name:      "test",
		Email:     "test@localhost",
		Username:  "test"}
	hits := 0

	testMux.HandleFunc("/a/accounts/self", func(w http.ResponseWriter, r *http.Request) {
		hits++
		switch hits {
		case 1:
			writeresponse(t, w, nil, http.StatusUnauthorized)
		case 2:
			writeresponse(t, w, nil, http.StatusUnauthorized)
		case 3:
			if r.Header.Get("Cookie") != "admin=secret" {
				t.Error("Expected cookie to equal 'admin=secret")
			}

			writeresponse(t, w, account, http.StatusOK)
		case 4:
			t.Error("Did not expect another request")
		}
	})

	serverURL := fmt.Sprintf("http://admin:secret@%s/", testServer.Listener.Addr().String())
	client, err := gerrit.NewClient(serverURL, nil)
	if err != nil {
		t.Error(err)
	}
	if !client.Authentication.HasCookieAuth() {
		t.Error("Expected HasCookieAuth() == true")
	}
}

func TestNewRequest(t *testing.T) {
	c, err := gerrit.NewClient(testGerritInstanceURL, nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "/foo", testGerritInstanceURL+"foo"
	inBody, outBody := &gerrit.PermissionRuleInfo{Action: "ALLOW", Force: true, Min: 0, Max: 0}, `{"action":"ALLOW","force":true,"min":0,"max":0}`+"\n"
	req, _ := c.NewRequest("GET", inURL, inBody)

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), outBody; got != want {
		t.Errorf("NewRequest Body is %v, want %v", got, want)
	}
}

func TestNewRawPutRequest(t *testing.T) {
	c, err := gerrit.NewClient(testGerritInstanceURL, nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}

	inURL, outURL := "/foo", testGerritInstanceURL+"foo"
	req, _ := c.NewRawPutRequest(inURL, "test raw PUT contents")

	// Test that relative URL was expanded
	if got, want := req.URL.String(), outURL; got != want {
		t.Errorf("NewRequest(%q) URL is %v, want %v", inURL, got, want)
	}

	// Test that body was JSON encoded
	body, _ := ioutil.ReadAll(req.Body)
	if got, want := string(body), "test raw PUT contents"; got != want {
		t.Errorf("NewRequest Body is %v, want %v", got, want)
	}
}

func testURLParseError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Expected error to be returned")
	}
	if err, ok := err.(*url.Error); !ok || err.Op != "parse" {
		t.Errorf("Expected URL parse error, got %+v", err)
	}
}

func TestNewRequest_BadURL(t *testing.T) {
	c, err := gerrit.NewClient(testGerritInstanceURL, nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}
	_, err = c.NewRequest("GET", ":", nil)
	testURLParseError(t, err)
}

// If a nil body is passed to gerrit.NewRequest, make sure that nil is also passed to http.NewRequest.
// In most cases, passing an io.Reader that returns no content is fine,
// since there is no difference between an HTTP request body that is an empty string versus one that is not set at all.
// However in certain cases, intermediate systems may treat these differently resulting in subtle errors.
func TestNewRequest_EmptyBody(t *testing.T) {
	c, err := gerrit.NewClient(testGerritInstanceURL, nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}
	req, err := c.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatalf("NewRequest returned unexpected error: %v", err)
	}
	if req.Body != nil {
		t.Fatalf("constructed request contains a non-nil Body")
	}
}

func TestDo(t *testing.T) {
	setup()
	defer teardown()

	type foo struct {
		A string
	}

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, `)]}'`+"\n"+`{"A":"a"}`)
	})

	req, err := testClient.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	body := new(foo)
	if _, err := testClient.Do(req, body); err != nil {
		t.Error(err)
	}

	want := &foo{"a"}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("Response body = %v, want %v", body, want)
	}
}

func TestDo_ioWriter(t *testing.T) {
	setup()
	defer teardown()
	content := `)]}'` + "\n" + `{"A":"a"}`

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if m := "GET"; m != r.Method {
			t.Errorf("Request method = %v, want %v", r.Method, m)
		}
		fmt.Fprint(w, content)
	})

	req, err := testClient.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error(err)
	}
	var buf []byte
	actual := bytes.NewBuffer(buf)
	if _, err := testClient.Do(req, actual); err != nil {
		t.Error(err)
	}

	expected := []byte(content)
	if !reflect.DeepEqual(actual.Bytes(), expected) {
		t.Errorf("Response body = %v, want %v", actual, string(expected))
	}
}

func TestDo_HTTPError(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Bad Request", 400)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected HTTP 400 error.")
	}
}

// Test handling of an error caused by the internal http client's Do() function.
// A redirect loop is pretty unlikely to occur within the Gerrit API, but does allow us to exercise the right code path.
func TestDo_RedirectLoop(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/", http.StatusFound)
	})

	req, _ := testClient.NewRequest("GET", "/", nil)
	_, err := testClient.Do(req, nil)

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*url.Error); !ok {
		t.Errorf("Expected a URL error; got %#v.", err)
	}
}

func TestRemoveMagicPrefixLine(t *testing.T) {
	mockData := []struct {
		Current, Expected []byte
	}{
		{[]byte(`{"A":"a"}`), []byte(`{"A":"a"}`)},
		{[]byte(`)]}'` + "\n" + `{"A":"a"}`), []byte(`{"A":"a"}`)},
	}
	for _, mock := range mockData {
		body := gerrit.RemoveMagicPrefixLine(mock.Current)
		if !reflect.DeepEqual(body, mock.Expected) {
			t.Errorf("Response body = %v, want %v", body, mock.Expected)
		}
	}
}

func TestRemoveMagicPrefixLineDoesNothingWithoutPrefix(t *testing.T) {
	mockData := []struct {
		Current, Expected []byte
	}{
		{[]byte(`{"A":"a"}`), []byte(`{"A":"a"}`)},
		{[]byte(`{"A":"a"}`), []byte(`{"A":"a"}`)},
	}
	for _, mock := range mockData {
		body := gerrit.RemoveMagicPrefixLine(mock.Current)
		if !reflect.DeepEqual(body, mock.Expected) {
			t.Errorf("Response body = %v, want %v", body, mock.Expected)
		}
	}
}

func TestNewClientFailsOnDeadConnection(t *testing.T) {
	setup()
	serverURL := fmt.Sprintf("http://admin:secret@%s/", testServer.Listener.Addr().String())
	teardown() // Closes the server
	_, err := gerrit.NewClient(serverURL, nil)
	if err == nil {
		t.Fatal("Expected err to not be nil")
	}
	if !strings.Contains(err.Error(), "connection refused") {
		t.Fatalf("Unexpected error. 'connected refused' not found in %s", err.Error())
	}
}
