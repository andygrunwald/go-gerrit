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
		opt := &gerrit.BranchOptions{
			Limit: limit,
			Skip:  skip,
		}
		fmt.Printf("ListBranches with skip %d and limit %d\n", skip, limit)
		branches, _, err := client.Projects.ListBranches(projectName, opt)
		if err != nil {
			panic(err)
		}

		for _, branch := range *branches {
			fmt.Printf("%s -> %s\n", branch.Ref, branch.Revision)
		}

		// Raising pagination pointer
		i++
		skip += limit
	}

	// ListBranches with skip 0 and limit 2
	// HEAD -> master
	// refs/meta/config -> 35fc56a537db06bbb9e5bf92f3a3e2e096d4f4c9
	// ListBranches with skip 2 and limit 2
	// refs/heads/infra/config -> 54bf8ddf9bcb4ffc748dd6fead98ececa95c98e5
	// refs/heads/master -> 5f8903effc18583a4796f2adfba13d416f874444
	// ListBranches with skip 4 and limit 2
	// refs/heads/stable-2.10 -> ed2d5cedd8d79ee224fcc2280a6f11e8175fc2b0
	// refs/heads/stable-2.11 -> 854e55ec22dadfc76b1112eb086f0d20dd4a977c
}
