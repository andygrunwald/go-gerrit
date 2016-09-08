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

func testMethod(t *testing.T, r *http.Request, want string) {
	if got := r.Method; got != want {
		t.Errorf("Request method: %v, want %v", got, want)
	}
}

func testRequestURL(t *testing.T, r *http.Request, want string) {
	if got := r.URL.String(); got != want {
		t.Errorf("Request URL: %v, want %v", got, want)
	}
}

func testFormValues(t *testing.T, r *http.Request, values testValues) {
	want := url.Values{}
	for k, v := range values {
		want.Add(k, v)
	}

	r.ParseForm()
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

func TestNewRequest_InvalidJSON(t *testing.T) {
	c, err := gerrit.NewClient(testGerritInstanceURL, nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}

	type T struct {
		A map[int]interface{}
	}
	_, err = c.NewRequest("GET", "/", &T{})

	if err == nil {
		t.Error("Expected error to be returned.")
	}
	if err, ok := err.(*json.UnsupportedTypeError); !ok {
		t.Errorf("Expected a JSON error; got %#v.", err)
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

	req, _ := testClient.NewRequest("GET", "/", nil)
	responseBody := new(foo)
	testClient.Do(req, responseBody, nil)

	want := &foo{"a"}
	if !reflect.DeepEqual(responseBody, want) {
		t.Errorf("Response body = %v, want %v", responseBody, want)
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

	req, _ := testClient.NewRequest("GET", "/", nil)
	var buf []byte
	actual := bytes.NewBuffer(buf)
	testClient.Do(req, actual, nil)

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
	_, err := testClient.Do(req, nil, nil)

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
	_, err := testClient.Do(req, nil, nil)

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
