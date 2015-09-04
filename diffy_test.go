package diffy

import (
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestNewClient_NoGerritInstance(t *testing.T) {
	mockData := []string{"", "://not-existing"}
	for _, data := range mockData {
		c, err := NewClient(data, nil)
		if c != nil {
			t.Errorf("NewClient return is not nil. Expected no client. Go %+v", c)
		}
		if err == nil {
			t.Error("No error occured by empty Gerrit Instance. Expected one.")
		}
	}
}

func TestNewClient_HttpClient(t *testing.T) {
	customHTTPClient := &http.Client{
		Timeout: 30 * time.Second,
	}
	mockData := []struct {
		HTTPClient         *http.Client
		ExpectedHTTPClient *http.Client
	}{
		{nil, http.DefaultClient},
		{customHTTPClient, customHTTPClient},
	}

	for _, mock := range mockData {
		c, err := NewClient("https://gerrit-review.googlesource.com/", mock.HTTPClient)
		if err != nil {
			t.Errorf("An error occured. Expected nil. Got %+v.", err)
		}
		if reflect.DeepEqual(c.client, mock.ExpectedHTTPClient) == false {
			t.Errorf("Wrong HTTP Client. Expected %+v. Got %+v", mock.ExpectedHTTPClient, c.client)
		}
	}
}

func TestNewClient_Services(t *testing.T) {
	c, err := NewClient("https://gerrit-review.googlesource.com/", nil)
	if err != nil {
		t.Errorf("An error occured. Expected nil. Got %+v.", err)
	}

	if c.Access == nil {
		t.Error("No AccessService found.")
	}
	if c.Accounts == nil {
		t.Error("No AccountsService found.")
	}
	if c.Changes == nil {
		t.Error("No ChangesService found.")
	}
	if c.Config == nil {
		t.Error("No ConfigService found.")
	}
	if c.Groups == nil {
		t.Error("No GroupsService found.")
	}
	if c.Plugins == nil {
		t.Error("No PluginsService found.")
	}
	if c.Projects == nil {
		t.Error("No ProjectsService found.")
	}
}
