package gerrit_test

import (
	"context"
	"encoding/json"
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

		_, err := fmt.Fprint(w, `[{"_account_id": 1}]`)
		if err != nil {
			t.Error(err)
		}
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

		_, err := fmt.Fprint(w, `[{"account": {"_account_id": 1}}]`)
		if err != nil {
			t.Error(err)
		}
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

		_, err := fmt.Fprint(w, `{"_account_id": 1}`)
		if err != nil {
			t.Error(err)
		}
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

		_, err := fmt.Fprint(w, `{"confirm": true}`)
		if err != nil {
			t.Error(err)
		}
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
		_, err := fmt.Fprint(w, `{"Code-Review": 2, "Verified": 1}`)
		if err != nil {
			t.Error(err)
		}
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

func TestChangesService_AddReviewer_WithState(t *testing.T) {
	testCases := []struct {
		name          string
		input         *gerrit.ReviewerInput
		expectedState string
	}{
		{
			name: "AddAsCC",
			input: &gerrit.ReviewerInput{
				Reviewer: "user@example.com",
				State:    "CC",
			},
			expectedState: "CC",
		},
		{
			name: "AddAsReviewer",
			input: &gerrit.ReviewerInput{
				Reviewer: "user@example.com",
				State:    "REVIEWER",
			},
			expectedState: "REVIEWER",
		},
		{
			name: "RemoveReviewer",
			input: &gerrit.ReviewerInput{
				Reviewer: "user@example.com",
				State:    "REMOVED",
			},
			expectedState: "REMOVED",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "POST" {
					t.Error("Method != POST")
				}

				// Verify the request body contains the expected state
				var input gerrit.ReviewerInput
				if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
					t.Errorf("Failed to decode request body: %v", err)
				}

				if input.State != tc.expectedState {
					t.Errorf("State = %s, want %s", input.State, tc.expectedState)
				}

				if input.Reviewer != "user@example.com" {
					t.Errorf("Reviewer = %s, want user@example.com", input.Reviewer)
				}

				_, err := fmt.Fprint(w, `{"input": "user@example.com"}`)
				if err != nil {
					t.Error(err)
				}
			}))
			defer ts.Close()

			ctx := context.Background()
			client := newClient(ctx, t, ts)
			data, _, err := client.Changes.AddReviewer(ctx, "123", tc.input)
			if err != nil {
				t.Error(err)
			}
			if data.Input != "user@example.com" {
				t.Errorf("Input = %s, want user@example.com", data.Input)
			}
		})
	}
}

func TestChangesService_AddReviewer_WithNotify(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Method != POST")
		}

		var input gerrit.ReviewerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if input.Notify != "OWNER" {
			t.Errorf("Notify = %s, want OWNER", input.Notify)
		}

		if input.State != "CC" {
			t.Errorf("State = %s, want CC", input.State)
		}

		_, err := fmt.Fprint(w, `{"ccs": [{"_account_id": 1}]}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.AddReviewer(ctx, "123", &gerrit.ReviewerInput{
		Reviewer: "user@example.com",
		State:    "CC",
		Notify:   "OWNER",
	})
	if err != nil {
		t.Error(err)
	}
	if len(data.CCS) != 1 {
		t.Errorf("CCS length = %d, want 1", len(data.CCS))
	}
}

func TestChangesService_AddReviewer_WithConfirmed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Error("Method != POST")
		}

		var input gerrit.ReviewerInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}

		if !input.Confirmed {
			t.Error("Confirmed should be true")
		}

		if input.State != "REVIEWER" {
			t.Errorf("State = %s, want REVIEWER", input.State)
		}

		_, err := fmt.Fprint(w, `{"reviewers": [{"_account_id": 1}, {"_account_id": 2}]}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client := newClient(ctx, t, ts)
	data, _, err := client.Changes.AddReviewer(ctx, "123", &gerrit.ReviewerInput{
		Reviewer:  "mygroup",
		State:     "REVIEWER",
		Confirmed: true,
	})
	if err != nil {
		t.Error(err)
	}
	if len(data.Reviewers) != 2 {
		t.Errorf("Reviewers length = %d, want 2", len(data.Reviewers))
	}
}
