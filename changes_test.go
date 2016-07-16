package gerrit_test

import (
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func ExampleChangesService_QueryChanges() {
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{
		"status:merged",
	}
	opt.Limit = 2
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(opt)

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d\n", change.Project, change.Subject, instance, change.Number)
	}

	// Output:
	// Project: platform/external/llvm -> [InstCombine] allow X + signbit --> X ^ signbit for vector splats am: 2e9433d42d -> https://android-review.googlesource.com/248616
	// Project: platform/external/llvm -> add vector test to show missing transform am: 713ceaf392 -> https://android-review.googlesource.com/248615
}
