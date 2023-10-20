package gerrit_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

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
		fmt.Fprintf(w, `{ "id": "abc1234", "project": "%s", "branch": "%s", "subject": "%s"}`, project, branch, subject)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
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
			info, _, err := client.Changes.CreateChange(&input)
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

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.ReadyForReviewInput{Message: "Now ready for review"}
	_, err = client.Changes.SetReadyForReview("123", cm)
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

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}
	cm := &gerrit.ReadyForReviewInput{Message: "Now ready for review"}
	resp, err := client.Changes.SetReadyForReview("123", cm)
	if err == nil {
		t.Error("Expected error, instead nil")
	}
	if resp.StatusCode != http.StatusNotFound {
		t.Error("Expected 404 code")
	}
}

func TestChangesService_ListChangeRobotComment(t *testing.T) {
	listRobotCommentResponse := `)]}'
		{
		  "main.c": [
		    {
		      "robot_id": "robot",
		      "robot_run_id": "1",
		      "fix_suggestions": [
		        {
		          "fix_id": "c3302a6f_1578ee9e",
		          "description": "suggestion",
		          "replacements": [
		            {
		              "path": "main.c",
		              "range": {
		                "start_line": 3,
		                "start_character": 0,
		                "end_line": 5,
		                "end_character": 1
		              },
		              "replacement": "int main() { printf(\"Hello world!\"); }"
		            }
		          ]
		        }
		      ],
		      "author": {
		        "_account_id": 1000000,
		        "name": "Jhon Smith",
		        "username": "jhon"
		      },
		      "change_message_id": "517e6c92e4bd105f1e611294a3010ea177771551",
		      "patch_set": 1,
		      "id": "f7d3d07f_bf17b66e",
		      "line": 5,
		      "range": {
		        "start_line": 3,
		        "start_character": 0,
		        "end_line": 5,
		        "end_character": 1
		      },
		      "updated": "2022-07-05 13:44:34.000000000",
		      "message": "[clang-format] fix suggestion",
		      "commit_id": "be8ce493368f2ce0fa73b56f5e2bd0dc17ca4359"
		    }
		  ]
		}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/changeID/robotcomments" {
			t.Errorf("%s != /changes/changeID/robotcomments", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Error("Method != GET")
		}
		fmt.Fprint(w, listRobotCommentResponse)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}

	got, _, err := client.Changes.ListChangeRobotComments("changeID")
	if err != nil {
		t.Errorf("Changes.ListChangeRobotComments returned error: %v", err)
	}

	want := map[string][]gerrit.RobotCommentInfo{
		"main.c": {
			{
				CommentInfo: gerrit.CommentInfo{
					PatchSet: 1,
					ID:       "f7d3d07f_bf17b66e",
					Line:     5,
					Range: &gerrit.CommentRange{
						StartLine:      3,
						StartCharacter: 0,
						EndLine:        5,
						EndCharacter:   1,
					},
					Message: "[clang-format] fix suggestion",
					Updated: &gerrit.Timestamp{
						Time: time.Date(2022, 7, 5, 13, 44, 34, 0, time.UTC),
					},
					Author: gerrit.AccountInfo{
						AccountID: 1000000,
						Name:      "Jhon Smith",
						Username:  "jhon",
					},
				},
				RobotID:    "robot",
				RobotRunID: "1",
				FixSuggestions: []gerrit.FixSuggestionInfo{
					{
						FixID:       "c3302a6f_1578ee9e",
						Description: "suggestion",
						Replacements: []gerrit.FixReplacementInfo{
							{
								Path: "main.c",
								Range: gerrit.CommentRange{
									StartLine:      3,
									StartCharacter: 0,
									EndLine:        5,
									EndCharacter:   1,
								},
								Replacement: "int main() { printf(\"Hello world!\"); }",
							},
						},
					},
				},
			},
		},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Change.ListChangeRobotComments:\ngot: %+v\nwant: %+v", got, want)
	}
}
