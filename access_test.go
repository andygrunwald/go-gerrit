package gerrit_test

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/andygrunwald/go-gerrit"
)

func TestAccessService_ListAccessRights(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/access/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, testValues{
			"project": "go",
		})

		fmt.Fprint(w, `)]}'`+"\n"+`{"go":{"revision":"08f45ba74baef9699b650f42022df6467389c1f0","inherits_from":{"id":"All-Projects","name":"All-Projects","description":"Access inherited by all other projects.","state":"ACTIVE"},"local":{},"owner_of":[],"config_visible":false}}`)
	})

	opt := &gerrit.ListAccessRightsOptions{
		Project: []string{"go"},
	}
	access, _, err := testClient.Access.ListAccessRights(opt)
	if err != nil {
		t.Errorf("Access.ListAccessRights returned error: %v", err)
	}

	want := &map[string]gerrit.ProjectAccessInfo{
		"go": {
			Revision: "08f45ba74baef9699b650f42022df6467389c1f0",
			InheritsFrom: gerrit.ProjectInfo{
				ID:          "All-Projects",
				Name:        "All-Projects",
				Parent:      "",
				Description: "Access inherited by all other projects.",
				State:       "ACTIVE",
			},
			Local:         map[string]gerrit.AccessSectionInfo{},
			OwnerOf:       []string{},
			ConfigVisible: false,
		},
	}
	if !reflect.DeepEqual(access, want) {
		t.Errorf("Access.ListAccessRights returned %+v, want %+v", access, want)
	}
}

func TestAccessService_ListAccessRights_WithoutOpts(t *testing.T) {
	setup()
	defer teardown()

	testMux.HandleFunc("/access/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")

		fmt.Fprint(w, `)]}'`+"\n"+`{}`)
	})

	access, _, err := testClient.Access.ListAccessRights(nil)
	if err != nil {
		t.Errorf("Access.ListAccessRights returned error: %v", err)
	}

	want := &map[string]gerrit.ProjectAccessInfo{}
	if !reflect.DeepEqual(access, want) {
		t.Errorf("Access.ListAccessRights returned %+v, want %+v", access, want)
	}
}
