package diffy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

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

	want := &ProjectInfo{
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

// +func (s *ProjectsService) CreateProject(name string, input *ProjectInput) (*ProjectInfo, *Response, error) {
func TestProjectsService_CreateProject(t *testing.T) {
	setup()
	defer teardown()

	input := &ProjectInput{
		Description: "The Go Programming Language",
	}

	testMux.HandleFunc("/projects/go/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		v := new(ProjectInput)
		json.NewDecoder(r.Body).Decode(v)

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `)]}'`+"\n"+`{"id":"go","name":"go","parent":"All-Projects","description":"The Go Programming Language"}`)
	})

	project, _, err := testClient.Projects.CreateProject("go", input)
	if err != nil {
		t.Errorf("Projects.CreateProject returned error: %v", err)
	}

	want := &ProjectInfo{
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

	if !reflect.DeepEqual(*description, want) {
		t.Errorf("Projects.GetProjectDescription returned %+v, want %+v", *description, want)
	}
}
