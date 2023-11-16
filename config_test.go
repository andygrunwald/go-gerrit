package gerrit_test

import (
	"context"
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func ExampleConfigService_GetVersion() {
	instance := "https://gerrit-review.googlesource.com/"
	ctx := context.Background()
	client, err := gerrit.NewClient(ctx, instance, nil)
	if err != nil {
		panic(err)
	}

	v, _, err := client.Config.GetVersion(ctx)
	if err != nil {
		panic(err)
	}
	// We can`t output the direct version here, because
	// the test would fail if gerrit-review.googlesource.com
	// will upgrade their instance.
	// To access the version just print variable v
	if len(v) > 0 {
		fmt.Println("Got version from Gerrit")
	}

	// Output: Got version from Gerrit
}
