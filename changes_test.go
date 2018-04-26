package gerrit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func ExampleChangesService_QueryChanges() {
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244",
	}
	opt.Limit = 2
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(opt)
	if err != nil {
		panic(err)
	}

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d\n", change.Project, change.Subject, instance, change.Number)
	}

	// Output:
	// Project: platform/art -> ART: Change return types of field access entrypoints -> https://android-review.googlesource.com/249244
}

// Prior to fixing #18 this test would fail.
func ExampleChangesService_QueryChanges_withSymbols() {
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244+status:merged",
	}
	opt.Limit = 2
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(opt)
	if err != nil {
		panic(err)
	}

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d\n", change.Project, change.Subject, instance, change.Number)
	}

	// Output:
	// Project: platform/art -> ART: Change return types of field access entrypoints -> https://android-review.googlesource.com/249244
}

func ExampleChangesService_PublishChangeEdit() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "ok")
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		panic(err)
	}

	_, err = client.Changes.PublishChangeEdit("123", "NONE")
	if err != nil {
		panic(err)
	}
}

func TestChangesService_SubmitChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/submit" {
			t.Errorf("%s != /changes/123/submit", r.URL.Path)
		}
		fmt.Fprint(w, `{"id": "123"}`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.SubmitChange("123", nil)
	if err != nil {
		t.Error(err)
	}
	if info.ID != "123" {
		t.Error("Invalid id")
	}
}

func TestChangesService_SubmitChange_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.SubmitChange("123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_AbandonChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/abandon" {
			t.Errorf("%s != /changes/123/abandon", r.URL.Path)
		}
		fmt.Fprint(w, `{"id": "123"}`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.AbandonChange("123", nil)
	if err != nil {
		t.Error(err)
	}
	if info.ID != "123" {
		t.Error("Invalid id")
	}
}

func TestChangesService_AbandonChange_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.AbandonChange("123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RebaseChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/rebase" {
			t.Errorf("%s != /changes/123/rebase", r.URL.Path)
		}
		fmt.Fprint(w, `{"id": "123"}`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RebaseChange("123", nil)
	if err != nil {
		t.Error(err)
	}
	if info.ID != "123" {
		t.Error("Invalid id")
	}
}

func TestChangesService_RebaseChange_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RebaseChange("123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RestoreChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/restore" {
			t.Errorf("%s != /changes/123/restore", r.URL.Path)
		}
		fmt.Fprint(w, `{"id": "123"}`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RestoreChange("123", nil)
	if err != nil {
		t.Error(err)
	}
	if info.ID != "123" {
		t.Error("Invalid id")
	}
}

func TestChangesService_RestoreChange_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RestoreChange("123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RevertChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/revert" {
			t.Errorf("%s != /changes/123/revert", r.URL.Path)
		}
		fmt.Fprint(w, `{"id": "123"}`)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RevertChange("123", nil)
	if err != nil {
		t.Error(err)
	}
	if info.ID != "123" {
		t.Error("Invalid id")
	}
}

func TestChangesService_RevertChange_Conflict(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusConflict)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RevertChange("123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}
