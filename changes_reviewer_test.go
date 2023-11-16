package gerrit_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func newClient(ctx context.Context, t *testing.T, server *httptest.Server) *gerrit.Client {
	client, err := gerrit.NewClient(ctx, server.URL, nil)
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.ListReviewers(ctx, "123")
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.SuggestReviewers(ctx, "123", nil)
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.GetReviewer(ctx, "123", "1")
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.AddReviewer(ctx, "123", &gerrit.ReviewerInput{})
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	_, err := client.Changes.DeleteReviewer(ctx, "123", "1")
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	votes, _, err := client.Changes.ListVotes(ctx, "123", "1")
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

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	_, err := client.Changes.DeleteVote(ctx, "123", "1", "Code-Review", nil)
	if err != nil {
		t.Error(err)
	}
}
