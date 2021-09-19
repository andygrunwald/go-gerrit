package main

import (
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := "https://gerrit-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	i := 0
	limit := 2
	skip := 0
	for i < 3 {
		// Showcasing pagination
		opt := &gerrit.ProjectOptions{
			ProjectBaseOptions: gerrit.ProjectBaseOptions{
				Limit: limit,
				Skip:  skip,
			},
			Description: true,
		}
		fmt.Printf("ListProjects with skip %d and limit %d\n", skip, limit)
		projects, _, err := client.Projects.ListProjects(opt)
		if err != nil {
			panic(err)
		}

		for name, p := range *projects {
			fmt.Printf("%s - State: %s\n", name, p.State)
		}

		// Raising pagination pointer
		i++
		skip += limit
	}

	// ListProjects with skip 0 and limit 2
	// Core-Plugins - State: ACTIVE
	// Public-Plugins - State: ACTIVE
	// ListProjects with skip 2 and limit 2
	// Public-Projects - State: ACTIVE
	// TestRepo - State: ACTIVE
	// ListProjects with skip 4 and limit 2
	// apps/analytics-etl - State: ACTIVE
	// apps/kibana-dashboard - State: ACTIVE
}
