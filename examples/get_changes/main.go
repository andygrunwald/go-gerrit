package main

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := "https://android-review.googlesource.com/"
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{"project:kernel/common"}
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(ctx, opt)
	if err != nil {
		panic(err)
	}

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d\n", change.Project, change.Subject, instance, change.Number)
	}

	// Project: kernel/common -> ANDROID: GKI: Update symbols to symbol list -> https://android-review.googlesource.com/1830553
	// Project: kernel/common -> ANDROID: db845c_gki.fragment: Remove CONFIG_USB_NET_AX8817X from fragment -> https://android-review.googlesource.com/1830439
	// Project: kernel/common -> ANDROID: Update the ABI representation -> https://android-review.googlesource.com/1830469
	// ...
}
