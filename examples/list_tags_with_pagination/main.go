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

	projectName := "gerrit"

	i := 0
	limit := 2
	skip := 0
	for i < 3 {
		// Showcasing pagination
		opt := &gerrit.ProjectBaseOptions{
			Limit: limit,
			Skip:  skip,
		}
		fmt.Printf("ListTags with skip %d and limit %d\n", skip, limit)
		tags, _, err := client.Projects.ListTags(projectName, opt)
		if err != nil {
			panic(err)
		}

		for _, tag := range *tags {
			fmt.Printf("%s -> %s\n", tag.Message, tag.Ref)
		}

		// Raising pagination pointer
		i++
		skip += limit
	}

	// ListTags with skip 0 and limit 2
	// gerrit 2.0 -> refs/tags/v2.0
	// gerrit 2.0-rc0 -> refs/tags/v2.0-rc0
	// ListTags with skip 2 and limit 2
	// gerrit 2.0.1 -> refs/tags/v2.0.1
	// gerrit 2.0.10 -> refs/tags/v2.0.10
	// ListTags with skip 4 and limit 2
	// gerrit 2.0.11 -> refs/tags/v2.0.11
	// gerrit 2.0.12 -> refs/tags/v2.0.12
}
