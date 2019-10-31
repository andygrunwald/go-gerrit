package gerrit_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestProjectsService_ListProjects(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, testValues{
			"r": "(arch|benchmarks)",
			"n": "2",
		})

		fmt.Fprint(w, `)]}'`+"\n"+`{"arch":{"id":"arch","state":"ACTIVE"},"benchmarks":{"id":"benchmarks","state":"ACTIVE"}}`)
	})

	opt := &gerrit.ProjectOptions{
		Regex: "(arch|benchmarks)",
	}
	opt.Limit = 2
	project, _, err := testClient.Projects.ListProjects(opt)
	if err != nil {
		t.Errorf("Projects.ListProjects returned error: %v", err)
	}

	want := &map[string]gerrit.ProjectInfo{
		"arch": {
			ID:    "arch",
			State: "ACTIVE",
		},
		"benchmarks": {
			ID:    "benchmarks",
			State: "ACTIVE",
		},
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.ListProjects returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/go/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"go","name":"go","parent":"All-Projects","description":"The Go Programming Language","state":"ACTIVE"}`)
	})

	project, _, err := testClient.Projects.GetProject("go")
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "go",
		Name:        "go",
		Parent:      "All-Projects",
		Description: "The Go Programming Language",
		State:       "ACTIVE",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

func TestProjectsService_GetProject_WithSlash(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/plugins/delete-project", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testRequestURL(t, r, "/projects/plugins%2Fdelete-project")

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"plugins%2Fdelete-project","name":"plugins/delete-project","parent":"Public-Plugins","description":"A plugin which allows projects to be deleted from Gerrit via an SSH command","state":"ACTIVE"}`)
	})
	project, _, err := testClient.Projects.GetProject("plugins/delete-project")
	if err != nil {
		t.Errorf("Projects.GetProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "plugins%2Fdelete-project",
		Name:        "plugins/delete-project",
		Parent:      "Public-Plugins",
		Description: "A plugin which allows projects to be deleted from Gerrit via an SSH command",
		State:       "ACTIVE",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.GetProject returned %+v, want %+v", project, want)
	}
}

// +func (s *ProjectsService) CreateProject(name string, input *ProjectInput) (*ProjectInfo, *Response, error) {
func TestProjectsService_CreateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &gerrit.ProjectInput{
		Description: "The Go Programming Language",
	}

	testMux.HandleFunc("/projects/go/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		v := new(gerrit.ProjectInput)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"go","name":"go","parent":"All-Projects","description":"The Go Programming Language"}`)
	})

	project, _, err := testClient.Projects.CreateProject("go", input)
	if err != nil {
		t.Errorf("Projects.CreateProject returned error: %v", err)
	}

	want := &gerrit.ProjectInfo{
		ID:          "go",
		Name:        "go",
		Parent:      "All-Projects",
		Description: "The Go Programming Language",
	}

	if !reflect.DeepEqual(project, want) {
		t.Errorf("Projects.CreateProject returned %+v, want %+v", project, want)
	}
}

// +func (s *ProjectsService) GetProjectDescription(name string) (*string, *Response, error) {
func TestProjectsService_GetProjectDescription(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/go/description/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `)]}'`+"\n"+`"The Go Programming Language"`)
	})

	description, _, err := testClient.Projects.GetProjectDescription("go")
	if err != nil {
		t.Errorf("Projects.GetProjectDescription returned error: %v", err)
	}

	want := "The Go Programming Language"

	if !reflect.DeepEqual(description, want) {
		t.Errorf("Projects.GetProjectDescription returned %+v, want %+v", description, want)
	}
}

func ExampleProjectsService_ListProjects() {
	instance := "https://chromium-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.ProjectOptions{
		Description: true,
		Prefix:      "infra/infra/infra_l",
	}
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		panic(err)
	}

	for name, p := range *projects {
		fmt.Printf("%s - State: %s\n", name, p.State)
	}

	// Output:
	// infra/infra/infra_libs - State: ACTIVE
}

func TestProjectsService_GetBranch(t *testing.T) {
	setup()
	defer teardown()

	existBranches := map[string]*gerrit.BranchInfo{
		"branch": {
			Ref:       "123",
			Revision:  "abcd1234",
			CanDelete: true,
		},
		"branch/foo": {
			Ref:       "456",
			Revision:  "deadbeef",
			CanDelete: false,
		},
	}

	testMux.HandleFunc("/projects/go/branches/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		branchName := r.URL.Path[len("/projects/go/branches/"):]

		branchInfo, ok := existBranches[branchName]
		if !ok {
			http.Error(w, branchName, http.StatusBadRequest)
		}

		branchInfoRaw, err := json.Marshal(&branchInfo)
		if err != nil {
			http.Error(w, branchName, http.StatusBadRequest)
		}

		fmt.Fprint(w, `)]}'`+"\n"+string(branchInfoRaw))
	})

	var tests = []struct {
		name     string
		branch   string
		expected *gerrit.BranchInfo
	}{
		{
			name:   "branch without slash",
			branch: "branch",
			expected: &gerrit.BranchInfo{
				Ref:       "123",
				Revision:  "abcd1234",
				CanDelete: true,
			},
		},
		{
			name:   "branch with slash",
			branch: "branch/foo",
			expected: &gerrit.BranchInfo{
				Ref:       "456",
				Revision:  "deadbeef",
				CanDelete: false,
			},
		},
	}

	for _, tc := range tests {
		branchInfo, _, err := testClient.Projects.GetBranch("go", tc.branch)
		if err != nil {
			t.Errorf("tc %s: Projects.GetProject returned error: %v", tc.name, err)
		}

		if !reflect.DeepEqual(branchInfo, tc.expected) {
			t.Errorf("tc %s: Projects.GetBranch returned %+v, want %+v", tc.name, branchInfo, tc.expected)
		}
	}
}

func TestProjectsService_ListAccessRights(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/MyProject/access", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		// from: https://gerrit-review.googlesource.com/Documentation/rest-api-projects.html#get-access
		resp := `{"can_add":true,"can_add_tags":true,"can_upload":true,"config_visible":true,"groups":{"c2ce4749a32ceb82cd6adcce65b8216e12afb41c":{"created_on":"2009-06-08 23:31:00.000000000","description":"Users who perform batch actions on Gerrit","group_id":2,"name":"Non-Interactive Users","options":{},"owner":"Administrators","owner_id":"d5b7124af4de52924ed397913e2c3b37bf186948","url":"#/admin/groups/uuid-c2ce4749a32ceb82cd6adcce65b8216e12afb41c"},"global:Anonymous-Users":{"name":"Anonymous Users","options":{}}},"inherits_from":{"description":"Access inherited by all other projects.","id":"All-Projects","name":"All-Projects"},"is_owner":true,"local":{"refs/*":{"permissions":{"read":{"rules":{"c2ce4749a32ceb82cd6adcce65b8216e12afb41c":{"action":"ALLOW","force":false},"global:Anonymous-Users":{"action":"ALLOW","force":false}}}}}},"owner_of":["refs/*"],"revision":"61157ed63e14d261b6dca40650472a9b0bd88474"}`
		_, _ = fmt.Fprint(w, `)]}'`+"\n"+resp)
	})

	projectAccessRight, _, err := testClient.Projects.ListAccessRights("MyProject")
	if err != nil {
		t.Errorf("project: list access rights error: %s", err)
	}

	fmt.Printf("project access rights: %v\n", projectAccessRight)
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

	projectAccessRight, _, err := testClient.Projects.AddUpdateDeleteAccessRights("MyProject", req)
	if err != nil {
		t.Errorf("project: add/update/delete access right error: %s", err)
	}

	fmt.Printf("project access rights: %v\n", projectAccessRight)
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

	accessCheckInfo, _, err := testClient.Projects.CheckAccess("MyProject", req)
	if err != nil {
		t.Errorf("project: access check error: %s", err)
	}

	fmt.Printf("project access check info: %v\n", accessCheckInfo)
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

	changeInfo, _, err := testClient.Projects.CreateAccessRightChange("MyProject", req)
	if err != nil {
		t.Errorf("project: create access change error: %s", err)
	}

	fmt.Printf("project create access change info: %v\n", changeInfo)
}
