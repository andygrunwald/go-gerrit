package gerrit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func newClient(t *testing.T, server *httptest.Server) *gerrit.Client {
	client, err := gerrit.NewClient(server.URL, nil)
	if err != nil {
		t.Error(err)
	}
	return client
}

func TestChangesService_ListReviewers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers/"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}

		fmt.Fprint(w, `[{"_account_id": 1}]`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	data, _, err := client.Changes.ListReviewers("123")
	if err != nil {
		t.Error(err)
	}

	if len(*data) != 1 {
		t.Error("Length of data !=1 ")
	}

	if (*data)[0].AccountID != 1 {
		t.Error("AccountID != 1")
	}
}

func TestChangesService_SuggestReviewers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/suggest_reviewers"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}

		fmt.Fprint(w, `[{"account": {"_account_id": 1}}]`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	data, _, err := client.Changes.SuggestReviewers("123", nil)
	if err != nil {
		t.Error(err)
	}

	if len(*data) != 1 {
		t.Error("Length of data !=1 ")
	}

	if (*data)[0].Account.AccountID != 1 {
		t.Error("AccountID != 1")
	}
}

func TestChangesService_GetReviewer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers/1"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}

		fmt.Fprint(w, `{"_account_id": 1}`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	data, _, err := client.Changes.GetReviewer("123", "1")
	if err != nil {
		t.Error(err)
	}
	if data.AccountID != 1 {
		t.Error("AccountID != 1")
	}
}

func TestChangesService_AddReviewer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}
		if r.Method != "POST" {
			t.Error("Method != POST")
		}

		fmt.Fprint(w, `{"confirm": true}`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	data, _, err := client.Changes.AddReviewer("123", &gerrit.ReviewerInput{})
	if err != nil {
		t.Error(err)
	}
	if !data.Confirm {
		t.Error("Confirm != true")
	}
}

func TestChangesService_DeleteReviewer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers/1"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}
		if r.Method != "DELETE" {
			t.Error("Method != DELETE")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	_, err := client.Changes.DeleteReviewer("123", "1")
	if err != nil {
		t.Error(err)
	}
}

func TestChangesService_ListVotes(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers/1/votes/"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}
		fmt.Fprint(w, `{"Code-Review": 2, "Verified": 1}`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	votes, _, err := client.Changes.ListVotes("123", "1")
	if err != nil {
		t.Error(err)
	}
	if votes["Code-Review"] != 2 {
		t.Error("Code-Review != 2")
	}
	if votes["Verified"] != 1 {
		t.Error("Verified != 1")
	}
}

func TestChangesService_DeleteVote(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := "/changes/123/reviewers/1/votes/Code-Review"
		if r.URL.Path != expected {
			t.Errorf("%s != %s", r.URL.Path, expected)
		}

		if r.Method != "DELETE" {
			t.Error("Method != DELETE")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	_, err := client.Changes.DeleteVote("123", "1", "Code-Review", nil)
	if err != nil {
		t.Error(err)
	}
}
