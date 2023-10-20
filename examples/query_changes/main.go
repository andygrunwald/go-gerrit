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
	start := 0
	for i < 3 {
		// Showcasing
		//	- Pagination
		//	- How to query for a project
		opt := &gerrit.QueryChangeOptions{
			QueryOptions: gerrit.QueryOptions{
				Query: []string{
					"project:gerrit",
				},
				Limit: limit,
			},
			Start: start,
		}
		fmt.Printf("QueryChanges with start %d and limit %d\n", start, limit)
		changes, _, err := client.Changes.QueryChanges(opt)
		if err != nil {
			panic(err)
		}

		for _, changeInfo := range *changes {
			fmt.Printf("%+v -> %+v\n", changeInfo.ID, changeInfo.Subject)
		}

		// Raising pagination pointer
		i++
		start += limit
	}

	// QueryChanges with start 0 and limit 2
	// gerrit~stable-3.3~I48896badc4a14927b98fef0311c4a63ba5b2251d -> Allow loading change notes from an existing Repository
	// gerrit~master~I644cdc74b679844b21054445e68b3e171e836834 -> Add support for 'is:attention' and 'has:attention'
	// QueryChanges with start 2 and limit 2
	// gerrit~master~Id3b9e395c977821f8957180a70de5379402d3621 -> Fix gr-identities "Link Another Identity" button test
	// gerrit~stable-3.2~I85db585f9f4799f35dfc0913ec10f473fe08c25a -> Fix serialization of AllUsersName and AllProjectsName
	// QueryChanges with start 4 and limit 2
	// gerrit~stable-3.2~Ic6520c1fc44ef51de8bb5d33e307e60ece630163 -> Log the result of git-upload-pack command in httpd_log
	// gerrit~stable-3.3~I668e0a322feb09a8dad47233baa6fd9585b4a8a9 -> Do not avertise ALL refs when HEAD is requested
}
