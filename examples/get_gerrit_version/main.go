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

	v, _, err := client.Config.GetVersion()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Version: %s", v)

	// Version: 3.4.1-2066-g8db5605430
}
