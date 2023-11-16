package gerrit_test

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

// +func (s *ProjectsService) GetIncludeIn(projectName, commitID string) (*IncludedInInfo, *Response, error){
func TestProjectsService_GetIncludeIn(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/projects/swift/commits/a8a477efffbbf3b44169bb9a1d3a334cbbd9aa96/in", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `)]}'`+"\n"+`{"branches": ["master"],"tags": ["1.1.0"]}`)
	})

	includedInInfo, _, err := testClient.Projects.GetIncludeIn(context.Background(), "swift", "a8a477efffbbf3b44169bb9a1d3a334cbbd9aa96")
	if err != nil {
		t.Errorf("Projects.GetIncludeIn returned error: %v", err)
	}

	want := &gerrit.IncludedInInfo{
		Branches: []string{
			"master",
		},
		Tags: []string{
			"1.1.0",
		},
	}

	if !reflect.DeepEqual(includedInInfo, want) {
		t.Errorf("Projects.GetIncludeIn returned %+v, want %+v", includedInInfo, want)
	}
}
