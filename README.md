# go-gerrit

[![GoDoc](https://pkg.go.dev/badge/github.com/andygrunwald/go-gerrit?utm_source=godoc)](https://pkg.go.dev/github.com/andygrunwald/go-gerrit)

go-gerrit is a [Go](https://golang.org/) client library for the [Gerrit Code Review](https://www.gerritcodereview.com/) system.

![go-gerrit - Go client/library for Gerrit Code Review](./img/logo.png "go-gerrit - Go client/library for Gerrit Code Review")

## Features

* [Authentication](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#AuthenticationService) (HTTP Basic, HTTP Digest, HTTP Cookie)
* Every API Endpoint like Gerrit
    * [/access/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#AccessService)
    * [/accounts/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#AccountsService)
    * [/changes/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#ChangesService)
    * [/config/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#ConfigService)
    * [/groups/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#GroupsService)
    * [/plugins/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#PluginsService)
    * [/projects/](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#ProjectsService)
* Supports optional plugin APIs such as
    * events-log - [About](https://gerrit.googlesource.com/plugins/events-log/+/master/src/main/resources/Documentation/about.md), [REST API](https://gerrit.googlesource.com/plugins/events-log/+/master/src/main/resources/Documentation/rest-api-events.md)

## Installation

_go-gerrit_ follows the [Go Release Policy](https://golang.org/doc/devel/release.html#policy).
This means we support the current + 2 previous Go versions.

It is go gettable ...

```sh
$ go get github.com/andygrunwald/go-gerrit
```

## API / Usage

Have a look at the [GoDoc documentation](https://pkg.go.dev/github.com/andygrunwald/go-gerrit) for a detailed API description.

The [Gerrit Code Review - REST API](https://gerrit-review.googlesource.com/Documentation/rest-api.html) was the foundation document.

### Authentication

Gerrit supports multiple ways for [authentication](https://gerrit-review.googlesource.com/Documentation/rest-api.html#authentication).

#### HTTP Basic

Some Gerrit instances (like [TYPO3](https://review.typo3.org/)) has [auth.gitBasicAuth](https://gerrit-review.googlesource.com/Documentation/config-gerrit.html#auth.gitBasicAuth) activated.
With this, you can authenticate with HTTP Basic like this:

```go
instance := "https://review.typo3.org/"
client, _ := gerrit.NewClient(instance, nil)
client.Authentication.SetBasicAuth("andy.grunwald", "my secrect password")

self, _, _ := client.Accounts.GetAccount("self")

fmt.Printf("Username: %s", self.Name)

// Username: Andy Grunwald
```

If you get a `401 Unauthorized`, check your Account Settings and have a look at the `HTTP Password` configuration.

#### HTTP Digest

Some Gerrit instances (like [Wikimedia](https://gerrit.wikimedia.org/)) has [Digest access authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) activated.

```go
instance := "https://gerrit.wikimedia.org/r/"
client, _ := gerrit.NewClient(instance, nil)
client.Authentication.SetDigestAuth("andy.grunwald", "my secrect http password")

self, resp, err := client.Accounts.GetAccount("self")

fmt.Printf("Username: %s", self.Name)

// Username: Andy Grunwald
```

If the chosen Gerrit instance does not support digest auth, an error like `WWW-Authenticate header type is not Digest` is thrown.

If you get a `401 Unauthorized`, check your Account Settings and have a look at the `HTTP Password` configuration.

#### HTTP Cookie

Some Gerrit instances hosted like the one hosted googlesource.com (e.g. [Go](https://go-review.googlesource.com/), [Android](https://android-review.googlesource.com/) or [Gerrit](https://gerrit-review.googlesource.com/)) support HTTP Cookie authentication.

You need the cookie name and the cookie value.
You can get them by click on "Settings > HTTP Password > Obtain Password" in your Gerrit instance.

There you can receive your values.
The cookie name will be (mostly) `o` (if hosted on googlesource.com).
Your cookie secret will be something like `git-your@email.com=SomeHash...`.

```go
instance := "https://gerrit-review.googlesource.com/"
client, _ := gerrit.NewClient(instance, nil)
client.Authentication.SetCookieAuth("o", "my-cookie-secret")

self, _, _ := client.Accounts.GetAccount("self")

fmt.Printf("Username: %s", self.Name)

// Username: Andy G.
```

## Examples

More examples are available

* in the [GoDoc examples section](https://pkg.go.dev/github.com/andygrunwald/go-gerrit#pkg-examples).
* in the [examples folder](./examples)

### Get version of Gerrit instance

Receive the version of the [Gerrit instance used by the Gerrit team](https://gerrit-review.googlesource.com/) for development:

```go
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
```

### Get all public projects

List all projects from [Chromium](https://chromium-review.googlesource.com/):

```go
package main

import (
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := "https://chromium-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.ProjectOptions{
		Description: true,
	}
	projects, _, err := client.Projects.ListProjects(opt)
	if err != nil {
		panic(err)
	}

	for name, p := range *projects {
		fmt.Printf("%s - State: %s\n", name, p.State)
	}

	// chromiumos/third_party/bluez - State: ACTIVE
	// external/github.com/Polymer/ShadowDOM - State: ACTIVE
	// external/github.com/domokit/mojo_sdk - State: ACTIVE
	// ...
}
```

### Query changes

Get some changes of the [kernel/common project](https://android-review.googlesource.com/#/q/project:kernel/common) from the [Android](http://source.android.com/)[Gerrit Review System](https://android-review.googlesource.com/).

```go
package main

import (
	"fmt"

	"github.com/andygrunwald/go-gerrit"
)

func main() {
	instance := "https://android-review.googlesource.com/"
	client, err := gerrit.NewClient(instance, nil)
	if err != nil {
		panic(err)
	}

	opt := &gerrit.QueryChangeOptions{}
	opt.Query = []string{"project:kernel/common"}
	opt.AdditionalFields = []string{"LABELS"}
	changes, _, err := client.Changes.QueryChanges(opt)
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
```

## Development

### Running tests and linters

Tests only:

```sh
$ make test
```

Checks, tests and linters

```sh
$ make vet staticcheck test
```

### Local Gerrit setup

For local development, we suggest the usage of the [official Gerrit Code Review docker image](https://hub.docker.com/r/gerritcodereview/gerrit):

```
$ docker run -ti -p 8080:8080 -p 29418:29418 gerritcodereview/gerrit:3.4.1
```

Wait a few minutes until the ```Gerrit Code Review NNN ready``` message appears,
where NNN is your current Gerrit version, then open your browser to http://localhost:8080
and you will be in Gerrit Code Review.

#### Authentication

For local development setups, go to http://localhost:8080/settings/#HTTPCredentials and click `GENERATE NEW PASSWORD`.
Now you can use (only for development purposes):

```go
client.Authentication.SetBasicAuth("admin", "secret")
```

Replace `secret` with your new value.

## Frequently Asked Questions (FAQ)

### How is the source code organized?

The source code organization is inspired by [go-github by Google](https://github.com/google/go-github).

Every REST API Endpoint (e.g. [`/access/`](https://gerrit-review.googlesource.com/Documentation/rest-api-access.html), [`/changes/`](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html)) is coupled in a service (e.g. [`AccessService` in access.go](./access.go), [`ChangesService` in changes.go](./changes.go)).
Every service is part of [`gerrit.Client`](./gerrit.go) as a member variable.

`gerrit.Client` can provide essential helper functions to avoid unnecessary code duplications, such as building a new request or parse responses.

Based on this structure, implementing a new API functionality is straight forward.
Here is an example of `*ChangeService.DeleteTopic*` / [DELETE /changes/{change-id}/topic](https://gerrit-review.googlesource.com/Documentation/rest-api-changes.html#delete-topic):

```go
func (s *ChangesService) DeleteTopic(changeID string) (*Response, error) {
    u := fmt.Sprintf("changes/%s/topic", changeID)
    return s.client.DeleteRequest(u, nil)
}
```

### What about the version compatibility with Gerrit?

The library was implemented based on the REST API of Gerrit version 2.11.3-1230-gb8336f1 and tested against this version.

This library might be working with older versions as well.
If you notice an incompatibility [open a new issue](https://github.com/andygrunwald/go-gerrit/issues/new).
We also appreciate your Pull Requests to improve this library.
We welcome contributions!

### What about adding code to support the REST API of an (optional) plugin?

It will depend on the plugin, and you are welcome to [open a new issue](https://github.com/andygrunwald/go-gerrit/issues/new) first to propose the idea and use-case.
As an example, the addition of support for `events-log` plugin was supported because the plugin itself is fairly
popular.
The structures that the REST API uses could also be used by `gerrit stream-events`.

## License

This project is released under the terms of the [MIT license](https://choosealicense.com/licenses/mit/).
