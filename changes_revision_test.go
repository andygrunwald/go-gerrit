package gerrit_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/andygrunwald/go-gerrit"
)

func TestChangesService_ListFiles(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.String(), "/changes/123/revisions/456/files/?base=7"; got != want {
			t.Errorf("request URL:\ngot:  %q\nwant: %q", got, want)
		}
		fmt.Fprint(w, `{
		  "/COMMIT_MSG": {
		    "status": "A",
		    "lines_inserted": 7,
		    "size_delta": 551,
		    "size": 551
		  },
		  "gerrit-server/RefControl.java": {
		    "lines_inserted": 5,
		    "lines_deleted": 3,
		    "size_delta": 98,
		    "size": 23348
		  }
		}`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	got, _, err := client.Changes.ListFiles("123", "456", &gerrit.FilesOptions{
		Base: "7",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]gerrit.FileInfo{
		"/COMMIT_MSG": {
			Status:        "A",
			LinesInserted: 7,
			SizeDelta:     551,
			Size:          551,
		},
		"gerrit-server/RefControl.java": {
			LinesInserted: 5,
			LinesDeleted:  3,
			SizeDelta:     98,
			Size:          23348,
		},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Changes.ListFiles:\ngot:  %+v\nwant: %+v", got, want)
	}
}

func TestChangesService_ListFilesReviewed(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got, want := r.URL.String(), "/changes/123/revisions/456/files/?q=abc&reviewed=true"; got != want {
			t.Errorf("request URL:\ngot:  %q\nwant: %q", got, want)
		}
		fmt.Fprint(w, `["/COMMIT_MSG","gerrit-server/RefControl.java"]`)
	}))
	defer ts.Close()

	client := newClient(t, ts)
	got, _, err := client.Changes.ListFilesReviewed("123", "456", &gerrit.FilesOptions{
		Q: "abc",
	})
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"/COMMIT_MSG", "gerrit-server/RefControl.java"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("client.Changes.ListFilesReviewed:\ngot:  %q\nwant: %q", got, want)
	}
}

func TestChangesService_ListRevisionRobotComment(t *testing.T) {
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
		if r.URL.Path != "/changes/changeID/revisions/revisionID/robotcomments" {
			t.Errorf("%s != /changes/changeID/revisions/revisionID/robotcomments", r.URL.Path)
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

	got, _, err := client.Changes.ListRevisionRobotComments("changeID", "revisionID")
	if err != nil {
		t.Errorf("Changes.ListRevisionRobotComments returned error: %v", err)
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
		t.Errorf("Change.ListRevisionRobotComments:\ngot: %+v\nwant: %+v", got, want)
	}
}

func TestChangesService_GetRobotComment(t *testing.T) {
	getRobotCommentResponse := `)]}'
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
		}`

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/changes/changeID/revisions/revisionID/robotcomments/commentID" {
			t.Errorf("%s != /changes/changeID/revisions/revisionID/robotcomments/commentID", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Error("Method != GET")
		}
		fmt.Fprint(w, getRobotCommentResponse)
	}))
	defer ts.Close()

	client, err := gerrit.NewClient(ts.URL, nil)
	if err != nil {
		t.Error(err)
	}

	got, _, err := client.Changes.GetRobotComment("changeID", "revisionID", "commentID")
	if err != nil {
		t.Errorf("Changes.GetRobotComment returned error: %v", err)
	}

	want := &gerrit.RobotCommentInfo{
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
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Change.GetRobotComment:\ngot: %+v\nwant: %+v", got, want)
	}
}
