package gerrit_test

import (
	"encoding/json"
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

func ExampleChangesService_QueryChanges_withSubmittable() {
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244",
	}
	opt.AdditionalFields = []string{"SUBMITTABLE"}

	changes, _, err := client.Changes.QueryChanges(opt)
	if err != nil {
		panic(err)
	}

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d, Ready to submit: %t\n", change.Project, change.Subject, instance, change.Number, change.Submittable)
	}

	// Output:
	// Project: platform/art -> ART: Change return types of field access entrypoints -> https://android-review.googlesource.com/249244, Ready to submit: false
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

func TestChangesService_CreateChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		decoder := json.NewDecoder(r.Body)
		var payload map[string]interface{}
		if err := decoder.Decode(&payload); err != nil {
			t.Error(err)
		}

		jsonStr, err := json.MarshalIndent(payload, "", " ")
		if err != nil {
			t.Error(err)
		}
		t.Logf("Request payload:\n%s", jsonStr)

		required := func(field string) string {
			value, ok := payload[field]
			strVal := value.(string)
			if !ok {
				t.Errorf("Missing required field %q", field)
			}
			return strVal
		}
		project := required("project")
		branch := required("branch")
		subject := required("subject")

		for field, generic := range payload {
			switch value := generic.(type) {
			case string:
				if len(value) == 0 {
					t.Errorf("Empty value for field %q", field)
				}
			}
		}

		if r.URL.Path != "/changes/" {
			t.Errorf("%s != /changes/", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("%s != POST", r.Method)
		}
		fmt.Fprintf(w, `{ "id": "abc1234", "project": "%s", "branch": "%s", "subject": "%s"}`, project, branch, subject)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.CreateChange(&gerrit.ChangeInput{
		Project: "myProject",
		Branch:  "main",
		Subject: "test change",
	})
	if err != nil {
		t.Error(err)
	}
	if info.ID != "abc1234" {
		t.Error("Invalid id")
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

func TestChangesService_SetCommitMessage(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/message" {
			t.Errorf("%s != /changes/123/message", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Error("Method != PUT")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.CommitMessageInput{Message: "New commit message"}
	_, err = client.Changes.SetCommitMessage("123", cm)
	if err != nil {
		t.Error(err)
	}
}

func TestChangesService_SetCommitMessage_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/message" {
			t.Errorf("%s != /changes/123/message", r.URL.Path)
		}
		if r.Method != "PUT" {
			t.Error("Method != PUT")
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.CommitMessageInput{Message: "New commit message"}
	resp, err := client.Changes.SetCommitMessage("123", cm)
	if err == nil {
		t.Error("Expected error, instead nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected 404 code")
	}
}
