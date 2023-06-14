package gerrit_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

// +func (s *ProjectsService) CreateTag(projectName, tagName string, input *TagInput) (*TagInfo, *Response, error)
func TestProjectsService_CreateTag(t *testing.T) {
	setup()
	defer teardown()

	input := &gerrit.TagInput{
		Ref:      "v1.0.0",
		Revision: "master",
		Message:  "v1.0.0 release",
	}

	testMux.HandleFunc("/projects/go/tags/v1.0.0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")

		v := new(gerrit.TagInput)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		fmt.Fprint(w, `)]}'`+"\n"+`{"ref":"v1.0.0","revision":"master","message":"v1.0.0 release"}`)
	})

	tag, _, err := testClient.Projects.CreateTag("go", "v1.0.0", input)
	if err != nil {
		t.Errorf("Projects.CreateTag returned error: %v", err)
	}

	want := &gerrit.TagInfo{
		Ref:      "v1.0.0",
		Revision: "master",
		Message:  "v1.0.0 release",
	}

	if !reflect.DeepEqual(tag, want) {
		t.Errorf("Projects.CreateTag returned %+v, want %+v", tag, want)
	}
}

func TestProjectsService_DeleteTag(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/go/tags/v1.0.0", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
		w.WriteHeader(http.StatusNoContent)
	})

	_, err := testClient.Projects.DeleteTag("go", "v1.0.0")
	if err != nil {
		t.Errorf("Projects.DeleteTag returned error: %v", err)
	}
}

func TestProjectsService_DeleteTags(t *testing.T) {
	setup()
	defer teardown()
	input := &gerrit.DeleteTagsInput{
		Tags: []string{"v1.0.0", "v1.1.0"},
	}
	testMux.HandleFunc("/projects/go/tags:delete", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		v := new(gerrit.DeleteTagsInput)
		if err := json.NewDecoder(r.Body).Decode(v); err != nil {
			t.Error(err)
		}

		if !reflect.DeepEqual(v, input) {
			t.Errorf("Request body = %+v, want %+v", v, input)
		}

		w.WriteHeader(http.StatusNoContent)
	})

	_, err := testClient.Projects.DeleteTags("go", input)
	if err != nil {
		t.Errorf("Projects.DeleteTags returned error: %v", err)
	}
}
