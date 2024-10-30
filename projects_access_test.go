package gerrit_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestProjectsService_ListAccessRights(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/MyProject/access", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// from: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-access
		resp := `{"can_add":true,"can_add_tags":true,"can_upload":true,"config_visible":true,"groups":{"c2ce4749a32ceb82cd6adcce65b8216e12afb41c":{"created_on":"2009-06-08 23:31:00.000000000","description":"Users who perform batch actions on Gerrit","group_id":2,"name":"Non-Interactive Users","options":{},"owner":"Administrators","owner_id":"d5b7124af4de52924ed397913e2c3b37bf186948","url":"#/admin/groups/uuid-c2ce4749a32ceb82cd6adcce65b8216e12afb41c"},"global:Anonymous-Users":{"name":"Anonymous Users","options":{}}},"inherits_from":{"description":"Access inherited by all other projects.","id":"All-Projects","name":"All-Projects"},"is_owner":true,"local":{"refs/*":{"permissions":{"read":{"rules":{"c2ce4749a32ceb82cd6adcce65b8216e12afb41c":{"action":"ALLOW","force":false},"global:Anonymous-Users":{"action":"ALLOW","force":false}}}}}},"owner_of":["refs/*"],"revision":"61157ed63e14d261b6dca40650472a9b0bd88474"}`
		_, _ = fmt.Fprint(w, `)]}'`+"\n"+resp)
	})

	projectAccessRight, _, err := testClient.Projects.ListAccessRights(context.Background(), "MyProject")
	if err != nil {
		t.Errorf("project: list access rights error: %s", err)
	}

	// Doing one deep check to verify the mapping
	if projectAccessRight.InheritsFrom.Name != "All-Projects" {
		t.Errorf("projectAccessRight.InheritsFrom.Name not matching. Expected '%s', got '%s'", "All-Projects", projectAccessRight.InheritsFrom.Name)
	}
}

func TestProjectsService_AddUpdateDeleteAccessRights(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/MyProject/access", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")

		// from: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-access
		resp := `{"revision":"61157ed63e14d261b6dca40650472a9b0bd88474","inherits_from":{"id":"All-Projects","name":"All-Projects","description":"Accessinheritedbyallotherprojects."},"local":{"refs/*":{"permissions":{"read":{"rules":{"global:Anonymous-Users":{"action":"ALLOW","force":false}}}}}},"is_owner":true,"owner_of":["refs/*"],"can_upload":true,"can_add":true,"can_add_tags":true,"config_visible":true,"groups":{"global:Anonymous-Users":{"options":{},"name":"AnonymousUsers"}}}`
		_, _ = fmt.Fprint(w, `)]}'`+"\n"+resp)
	})

	var data = `{"remove":{"refs/*":{"permissions":{"read":{"rules":{"c2ce4749a32ceb82cd6adcce65b8216e12afb41c":{"action":"ALLOW"}}}}}}}`
	var req = new(gerrit.ProjectAccessInput)
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		t.Errorf("project: add/update/delete access right request params error: %s", err)
	}

	projectAccessRight, _, err := testClient.Projects.AddUpdateDeleteAccessRights(context.Background(), "MyProject", req)
	if err != nil {
		t.Errorf("project: add/update/delete access right error: %s", err)
	}

	// Doing one deep check to verify the mapping
	if projectAccessRight.InheritsFrom.Name != "All-Projects" {
		t.Errorf("projectAccessRight.InheritsFrom.Name not matching. Expected '%s', got '%s'", "All-Projects", projectAccessRight.InheritsFrom.Name)
	}
}

func TestProjectsService_AccessCheck(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/MyProject/check.access", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		testQueryValues(t, r, testValues{
			"account": "1000098",
			"ref":     "refs/heads/secret/bla",
		})

		// from: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#set-access
		resp := `{"message": "user Kristen Burns \u003cKristen.Burns@gerritcodereview.com\u003e (1000098) cannot see ref refs/heads/secret/bla in project MyProject","status":403}`
		_, _ = fmt.Fprint(w, `)]}'`+"\n"+resp)
	})

	var data = `{"account":"1000098","ref":"refs/heads/secret/bla"}`
	var req = new(gerrit.CheckAccessOptions)
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		t.Errorf("project: access check request params error: %s", err)
	}

	accessCheckInfo, _, err := testClient.Projects.CheckAccess(context.Background(), "MyProject", req)
	if err != nil {
		t.Errorf("project: access check error: %s", err)
	}

	// Doing one deep check to verify the mapping
	if accessCheckInfo.Status != 403 {
		t.Errorf("accessCheckInfo.Status not matching. Expected '%d', got '%d'", 403, accessCheckInfo.Status)
	}
}

func TestProjectsService_CreateAccessChange(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/MyProject/access:review", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		// from: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#create-access-change
		resp := `{"id":"testproj~refs%2Fmeta%2Fconfig~Ieaf185bf90a1fc3b58461e399385e158a20b31a2","project":"testproj","branch":"refs/meta/config","hashtags":[],"change_id":"Ieaf185bf90a1fc3b58461e399385e158a20b31a2","subject":"Reviewaccesschange","status":"NEW","created":"2017-09-07 14:31:11.852000000","updated":"2017-09-07 14:31:11.852000000","submit_type":"CHERRY_PICK","mergeable":true,"insertions":2,"deletions":0,"unresolved_comment_count":0,"has_review_started":true,"_number":7,"owner":{"_account_id":1000000}}`
		_, _ = fmt.Fprint(w, `)]}'`+"\n"+resp)
	})

	var data = `{"add":{"refs/heads/*":{"permissions":{"read":{"rules":{"global:Anonymous-Users":{"action":"DENY","force":false}}}}}}}`
	var req = new(gerrit.ProjectAccessInput)
	if err := json.Unmarshal([]byte(data), &req); err != nil {
		t.Errorf("project: create access change request params error: %s", err)
	}

	changeInfo, _, err := testClient.Projects.CreateAccessRightChange(context.Background(), "MyProject", req)
	if err != nil {
		t.Errorf("project: create access change error: %s", err)
	}

	// Doing one deep check to verify the mapping
	if changeInfo.ChangeID != "Ieaf185bf90a1fc3b58461e399385e158a20b31a2" {
		t.Errorf("changeInfo.ChangeID not matching. Expected '%s', got '%s'", "Ieaf185bf90a1fc3b58461e399385e158a20b31a2", changeInfo.ChangeID)
	}
}
