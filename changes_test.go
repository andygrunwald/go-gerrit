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

func ExampleChangesService_QueryChanges() {
	ctx := context.Background()
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244",
	}
	opt.Limit = 2
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(ctx, opt)
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
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244",
	}
	opt.AdditionalFields = []string{"SUBMITTABLE"}

	changes, _, err := client.Changes.QueryChanges(ctx, opt)
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
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"change:249244 status:merged",
	}
	opt.Limit = 2
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(ctx, opt)
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
		_, err := fmt.Fprintf(w, "ok")
		if err != nil {
			panic(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		panic(err)
	}

	_, err = client.Changes.PublishChangeEdit(ctx, "123", "NONE")
	if err != nil {
		panic(err)
	}
}

func disallowEmptyFields(t *testing.T, payload map[string]interface{}, path string) {
	for field, generic := range payload {
		curPath := field
		if len(path) > 0 {
			curPath = path + "." + field
		}
		switch value := generic.(type) {
		case string:
			if len(value) == 0 {
				t.Errorf("Empty value for field %q", curPath)
			}
		case map[string]interface{}:
			if len(value) == 0 {
				t.Errorf("Empty value for field %q", curPath)
			}
			disallowEmptyFields(t, value, curPath)
		}
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
		if len(jsonStr) == 0 {
			t.Error("Empty request payload")
		}

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
		if merge, ok := payload["merge"]; ok {
			if _, ok := merge.(map[string]interface{})["source"]; !ok {
				t.Error(`Missing required field "merge.source"`)
			}
		}

		disallowEmptyFields(t, payload, "")

		if r.URL.Path != "/changes/" {
			t.Errorf("%s != /changes/", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Errorf("%s != POST", r.Method)
		}
		_, err = fmt.Fprintf(w, `{ "id": "abc1234", "project": "%s", "branch": "%s", "subject": "%s"}`, project, branch, subject)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}

	cases := map[string]gerrit.ChangeInput{
		"RequiredOnly": {
			Project: "myProject",
			Branch:  "main",
			Subject: "test change",
		},
		"WithMerge": {
			Project: "myProject",
			Branch:  "main",
			Subject: "test change",
			Merge: &gerrit.MergeInput{
				Source: "45/3/1",
			},
		},
		"WithAppend": {
			Project: "myProject",
			Branch:  "main",
			Subject: "test change",
			Author: &gerrit.AccountInput{
				Username: "roboto",
				Name:     "Rob Oto",
			},
		},
	}
	for name, input := range cases {
		t.Run(name, func(t *testing.T) {
			info, _, err := client.Changes.CreateChange(ctx, &input)
			if err != nil {
				t.Error(err)
			}

			if info.ID != "abc1234" {
				t.Error("Invalid id")
			}
		})
	}
}

func TestChangesService_SubmitChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/submit" {
			t.Errorf("%s != /changes/123/submit", r.URL.Path)
		}
		_, err := fmt.Fprint(w, `{"id": "123"}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.SubmitChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.SubmitChange(ctx, "123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_AbandonChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/abandon" {
			t.Errorf("%s != /changes/123/abandon", r.URL.Path)
		}
		_, err := fmt.Fprint(w, `{"id": "123"}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.AbandonChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.AbandonChange(ctx, "123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RebaseChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/rebase" {
			t.Errorf("%s != /changes/123/rebase", r.URL.Path)
		}
		_, err := fmt.Fprint(w, `{"id": "123"}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RebaseChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RebaseChange(ctx, "123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RestoreChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/restore" {
			t.Errorf("%s != /changes/123/restore", r.URL.Path)
		}
		_, err := fmt.Fprint(w, `{"id": "123"}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RestoreChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RestoreChange(ctx, "123", nil)
	if response.StatusCode != http.StatusConflict {
		t.Error("Expected 409 code")
	}
}

func TestChangesService_RevertChange(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/revert" {
			t.Errorf("%s != /changes/123/revert", r.URL.Path)
		}
		_, err := fmt.Fprint(w, `{"id": "123"}`)
		if err != nil {
			t.Error(err)
		}
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	info, _, err := client.Changes.RevertChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	_, response, _ := client.Changes.RevertChange(ctx, "123", nil)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.CommitMessageInput{Message: "New commit message"}
	_, err = client.Changes.SetCommitMessage(ctx, "123", cm)
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

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.CommitMessageInput{Message: "New commit message"}
	resp, err := client.Changes.SetCommitMessage(ctx, "123", cm)
	if err == nil {
		t.Error("Expected error, instead nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected 404 code")
	}
}

func TestChangesService_SetReadyForReview(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/ready" {
			t.Errorf("%s != /changes/123/ready", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Error("Method != POST")
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.ReadyForReviewInput{Message: "Now ready for review"}
	_, err = client.Changes.SetReadyForReview(ctx, "123", cm)
	if err != nil {
		t.Error(err)
	}
}

func TestChangesService_SetReadyForReview_NotFound(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/123/ready" {
			t.Errorf("%s != /changes/123/ready", r.URL.Path)
		}
		if r.Method != "POST" {
			t.Error("Method != POST")
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer ts.Close()

	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.ReadyForReviewInput{Message: "Now ready for review"}
	resp, err := client.Changes.SetReadyForReview(ctx, "123", cm)
	if err == nil {
		t.Error("Expected error, instead nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected 404 code")
	}
}
