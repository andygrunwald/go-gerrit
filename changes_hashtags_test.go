package gerrit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestChangesService_GetHashtags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			method = "GET"
			path   = "/changes/123/hashtags"
		)
		if r.Method != method {
			t.Errorf("Method %q != %q", r.Method, method)
		}
		if r.URL.Path != path {
			t.Errorf("Path %q != %q", r.URL.Path, path)
		}
		fmt.Fprintf(w, `)]}'
[
  "hashtag1",
  "hashtag2"
]
`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	hashtags, _, err := client.Changes.GetHashtags("123")
	if err != nil {
		t.Fatal(err)
	}

	if len(hashtags) != 2 || hashtags[0] != "hashtag1" || hashtags[1] != "hashtag2" {
		t.Errorf("Unexpected hashtags %+v", hashtags)
	}
}

func TestChangesService_SetHashtags(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const (
			method = "POST"
			path   = "/changes/123/hashtags"
		)
		if r.Method != method {
			t.Errorf("Method %q != %q", r.Method, method)
		}
		if r.URL.Path != path {
			t.Errorf("Path %q != %q", r.URL.Path, path)
		}
		fmt.Fprintf(w, `)]}'
[
  "hashtag1",
  "hashtag3"
]
`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	hashtags, _, err := client.Changes.SetHashtags("123", &gerrit.HashtagsInput{
		Add:    []string{"hashtag3"},
		Remove: []string{"hashtag2"},
	})
	if err != nil {
		t.Fatal(err)
	}

	if len(hashtags) != 2 || hashtags[0] != "hashtag1" || hashtags[1] != "hashtag3" {
		t.Errorf("Unexpected hashtags %+v", hashtags)
	}
}
