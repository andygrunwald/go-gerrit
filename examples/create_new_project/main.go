package main

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := fmt.Sprintf("http://%s:8080", "localhost")
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	// Get your credentials
	// For local development setups, go to
	// http://localhost:8080/settings/#HTTPCredentials
	// and click `GENERATE NEW PASSWORD`.
	// Replace `secret` with your new value.
	client.Authentication.SetBasicAuth("admin", "secret")

	gerritProject := "Example project to the moon"
	gerritBranch := "main"

	data := &gerrit.ProjectInput{
		Name:              gerritProject,
		Branches:          []string{gerritBranch},
		CreateEmptyCommit: true,
	}
	projectInfo, _, err := client.Projects.CreateProject(ctx, gerritProject, data)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Project '%s' created with ID '%+v'", projectInfo.Name, projectInfo.ID)

	// Project 'Example project to the moon' created with ID 'Example+project+to+the+moon'
}
