package main

import (
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := "https://chromium-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.ProjectOptions{
		Description: true,
	}
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		panic(err)
	}

	for name, p := range *projects {
		fmt.Printf("%s - State: %s\n", name, p.State)
	}

	// chromiumos/third_party/bluez - State: ACTIVE
	// external/github.com/Polymer/ShadowDOM - State: ACTIVE
	// external/github.com/domokit/mojo_sdk - State: ACTIVE
	// ...
}
