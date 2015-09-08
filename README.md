# diffy

[![GoDoc](https://godoc.org/github.com/andygrunwald/diffy?status.svg)](https://godoc.org/github.com/andygrunwald/diffy)
[![Build Status](https://travis-ci.org/andygrunwald/diffy.svg?branch=master)](https://travis-ci.org/andygrunwald/diffy)
[![Coverage Status](https://coveralls.io/repos/andygrunwald/diffy/badge.svg?branch=master&service=github)](https://coveralls.io/github/andygrunwald/diffy?branch=master)

diffy is a [Go(lang)](https://golang.org/) client library for accessing the [Gerrit Code Review](https://www.gerritcodereview.com/) API.

![Diffy - Go(lang) client/library for Gerrit Code Review](./img/diffy.png "Diffy - Go(lang) client/library for Gerrit Code Review")

## Features

* Authentication (HTTP Basic, HTTP Cookie)
* API complete
* TODO more features

## Installation

It is go gettable

    $ go get github.com/andygrunwald/diffy

(optional) to run unit / example tests:

    $ cd $GOPATH/src/github.com/andygrunwald/diffy
    $ go test -v

## API / Usage

Please have a look at the [GoDoc documentation](https://godoc.org/github.com/andygrunwald/diffy) for a detailed API description.

The [Gerrit Code Review - REST API](https://gerrit-review.googlesource.com/Documentation/rest-api.html) was the base document.

### Authentication

TODO

## Examples

Further a few examples how the API can be used.
A few more examples are available in the [GoDoc examples section](https://godoc.org/github.com/andygrunwald/diffy#pkg-examples).

### Get version of Gerrit instance

Receive the version of the [Gerrit instance used by the Gerrit team](https://gerrit-review.googlesource.com/) for development:

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/diffy"
)

func main() {
	instance := "https://gerrit-review.googlesource.com/"
	client, err := diffy.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	v, _, err := client.Config.GetVersion()
	fmt.Printf("Version: %s", *v)
	// Version: 2.11.3-1230-gb8336f1
}
```

### Get all public projects

List all projects from [cyanogenmod](http://review.cyanogenmod.org/):

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/diffy"
)

func main() {
	instance := "http://review.cyanogenmod.org/"
	client, err := diffy.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &diffy.ProjectOptions{
		Description: true,
	}
	projects, _, err := client.Projects.ListProjects(opt)
	for name, p := range *projects {
		fmt.Printf("%s - State: %s\n", name, p.State)
	}

	// CyanogenMod/android_external_drm - State: ACTIVE
	// CyanogenMod/android_external_jhead - State: ACTIVE
	// CyanogenMod/android_external_libppp - State: ACTIVE
	// ...
}
```

### Query changes

Get some changes of the [kernel/common project](https://android-review.googlesource.com/#/q/project:kernel/common) from the [Android](http://source.android.com/) [Gerrit Review System](https://android-review.googlesource.com/).

```go
package main

import (
	"fmt"
	"github.com/andygrunwald/diffy"
)

func main() {
	instance := "https://android-review.googlesource.com/"
	client, err := diffy.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &diffy.QueryChangeOptions{}
	opt.Query = []string{"project:kernel/common"}
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(opt)

	for _, change := range *changes {
		fmt.Printf("Project: %s -> %s -> %s%d\n", change.Project, change.Subject, instance, change.Number)
	}

	// Project: kernel/common -> android: binder: Fix BR_ERROR usage and change LSM denials to use it. -> https://android-review.googlesource.com/150839
	// Project: kernel/common -> android: binder: fix duplicate error return. -> https://android-review.googlesource.com/155031
	// Project: kernel/common -> dm-verity: Add modes and emit uevent on corrupted blocks -> https://android-review.googlesource.com/169572
}
```

## FAQ

### Where does the name come from?

*Diffy* is "The Kung Fu Review Cuckoo" by Gerrit itself.
As far as i know Diffy is the (official) logo of the Gerrit project.

All credits for name and logo goes to the Gerrit team.
Thank you!

### How is the source code organized?

The source code organisation was inspired by [go-github by Google](https://github.com/google/go-github).

Every REST API Endpoint (e.g. [/access/](https://gerrit-review.googlesource.com/Documentation/rest-api-access.html) or [/changes/](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html)) is coupled in a service (e.g. [AccessService in access.go](./access.go) or [ChangesService in changes.go](./changes.go)).
Every service is part of [diffy.Client](./diffy.go) as a member variable.

diffy.Client can provide basic helper functions to avoid unnecessary code duplications such as building a new request, parse responses and so on.

Based on this structure implementing a new API functionality is straight forwarded. Here is an example of *ChangeService.DeleteTopic* / [DELETE /changes/{change-id}/topic](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-topic):

```go
func (s *ChangesService) DeleteTopic(changeID string) (*Response, error) {
	u := fmt.Sprintf("changes/%s/topic", changeID)
	return s.client.DeleteRequest(u, nil)
}
```

### What about the version compatibility with Gerrit?

The library was implemented based on the REST API of Gerrit version 2.11.3-1230-gb8336f1 and tested against this version.

This library might be working with older versions as well.
If you notice an incompatibility [open a new issue](https://github.com/andygrunwald/diffy/issues/new) or try to fix it.
We welcome contribution!

## License

This project is released under the terms of the [MIT license](http://en.wikipedia.org/wiki/MIT_License).

