# Changelog

This is a high level log of changes, bugfixes, enhancements, etc
that have taken place between releases. Later versions are shown
first. For more complete details see
[the releases on GitHub.](https://github.com/andygrunwald/go-gerrit/releases)

## Versions

### 1.0.0 (2024-10-20)

This is the first release in 7 years (since 2017-11-04).
The project itself was alive and received new features and bug fixes during this time.
Only the creation of an "official" release was absent.

If you were/are using the latest head branch, it is safe to switch to v1.0.0.
No new breaking changes have been introduced.

If you are upgrading from v0.5.2 (2017-11-04), please check out the breaking changes.

**WARNING**: This release includes breaking changes. These changes are prefixed by `[BREAKING CHANGE]`.

#### Features

18c7753 Add support for submit_requirement fields in ChangeInfo. (#169)
0983e87 Update CherryPickInput to support new fields (#165)
bb1cabe Update SubmitInput entity to support new fields (#163)
444a268 Add missing fields `Error` and `ChangeInfo` to ReviewResult (#160)
235fa61 Implement parents_data in RevisionInfo (#157)
40e4f30 [BREAKING CHANGE] Add golang context support (#153)
f01a53b Add `Strategy`, `AllowConflicts`, `OnBehalfOfUploader`, `CommitterEmail` and `ValidationOptions` fields to RebaseInput (#155)
4f1401b AccountExternalIdInfo/Accounts.GetAccountExternalIDs: Adjusted docs according the official documentation
8a06083 Add AccountService.QueryAccounts
1d229fd Add AccountService.GetAccountExternalIDs
999473b Add fields `SecondaryEmails`, `Status`, `Inactive` and `Tags` in AccountInfo
c44fe2f Add create project's tag and delete project's tag and tags (#146)
2dcb9fb Add support for project deletion (with plugin) (#147)
423d372 add Comment{Input,Info} fields (#145)
8cd0d63 Add GetHashtags/SetHashtags (#141)
9ea5350 Add revision kind to RevisionInfo (#135)
5de64ee changes: Add RemoveAttention (#124)
2e881e2 Add fields `Reviewers`, `Ready`,  `WorkInProgress`, `AddToAttentionSet`, `RemoveFromAttentionSet` and `IgnoreAutomaticAttentionSetRules` to ReviewInput (#123)
f645b08 add: project commits included in (#122)
d3e91fb Add support for Change API's "Set Ready-For-Review" endpoint (#110)
ff14d06 Groups: Add _more_groups attribute to GroupInfo
67c9345 [BREAKING CHANGE] Fix #92: Add pagination support for SuggestAccount
67e5874 changes: fill in the rest of Changes.ChangeInput (#97)
101051e Fix #92: Add _more_accounts attribute and start parameter to paginate through SuggestAccount results
dcedff1 Add more supported fields to the ChangeInfo struct (#90)
4e82ec7 Add Changes.MoveChange (#81)
590067b Added go modules (#83)
ed2419a Add support for Change API's "Set Commit Message" endpoint (#80)
5c7c90c feat: add access method for project
b4b4fdb Add fields `Groups` and `ConfigWebLinks` to `ProjectAccessInfo`
5959a9b add ReviewerUpdates field to ChangeInfo
759bda1 feat: add submittable field for change info
64931d2 add Hashtags field to ChangeInfo
19ef3e9 add tests for project.GetBranch
197fe0d Add types for robot comments
5ea6031 Add Reviewers field to ChangeInfo
5632c7f Expose the client's BaseURL.
90fea2d Improve consistency of baseURL trailing slash handling.
70bbb05 [BREAKING CHANGE] Use Timestamp type for all time fields.
5416038 Add Avatars field to AccountInfo.
2e8da2e Add FilesOptions to ListFiles, ListFilesReviewed.
5ff0cbc Add missing fields to ChangeMessageInfo, FileInfo.
3418ea4 Add PatchOptions.Path field.
0655566 Added comment to DiffIntralineInfo from gerrit rest api docs (as suggested by @shurcooL).

#### Bugfixes

7c5be02 Fix MoveInput struct definition (keep_all_labels -> keep_all_votes) (#161)
b85f9f5 [BREAKING CHANGE] Change ReviewInput to use map[string]int for labels (#159)
b24a961 Escape % character in fm.Fprint (#106)
918f939 fix: remove a extra characters in url path (#103)
668ecf2 [BREAKING CHANGE] changes: rename Changes.DeleteDraftChange to DeleteChange (#95)
34f3353 fix: omit optional empty fields in ProjectInput (#94)
cc4e14e fix: no need to send content-type header if no header sent (#91)
eab0848 omitempty for AbandonInput.Notify and AbandonInput.NotifyDetails fields (#89)
3f5e365 fix: fix url error in `GroupsService.AddGroupMembers`, `GroupsService.DeleteGroupMember`, `ProjectsService.GetConfig`, `ProjectsService.SetProjectDescription`, `ProjectsService.DeleteProjectDescription`, `ProjectsService.BanCommit`, `ProjectsService.SetConfig`, `ProjectsService.SetHEAD`, `ProjectsService.SetProjectParent` and `ProjectsService.RunGC`
adf825f Fix JSON tags of `CheckAccessOptions` and added unit test `TestProjectsService_ListAccessRights`
6761407 Renamed projects_test.go to projects_access_test.go
ea89ae5 fix: Rename `ExampleChangesService_QueryChanges_WithSubmittable` to `ExampleChangesService_QueryChanges_withSubmittable`
7b0d1f8 escape branch names in the query as well
43cfd7a Fix GetCommitContent to actually get commit content
30ce279 Make code more readable and consistent regards return values
5f93656 Fixed handling of unhandled error
aab0406 [BREAKING CHANGE] Rename struct fields according go rules
3c4bc0f Fixed #44. DiffContent.A/B/AB should be []string instead of string. DiffIntralineInfo should be [][2]int.

#### Removal and deprecations

294d14b [BREAKING CHANGE] Remove CreateChangeDeprecated
7b2c737 api: deprecate CreateChange using ChangeInfo (use ChangeInput instead)

#### Chore and automation

d21ca62 Introduce goreleaser to make the release of a new version easier
08ff0c0 Github Actions: Upgrade supported Go versions from 1.21, 1.22 to 1.22, 1.23 (#172)
d093e01 Github Actions: Fix Actions trigger on branch name (#171)
b410a34 Github Actions: Upgrade `runs-on` image from ubuntu-22.04 to ubtuntu-24.04 (#170)
c725a0f Update Go version (add v1.22, remove v1.20) (#164)
e8bb8e4 Bump dominikh/staticcheck-action from 1.3.0 to 1.3.1 (#162)
6478d30 Bump actions/setup-go from 4 to 5 (#156)
89fb5cf Bump actions/checkout from 3 to 4 (#154)
d2361a1 Update Go version for testing: Remove 1.19, add 1.21 (#151)
e9d8f54 Update dependabot.yml: Set interval to monthly
da63a5c GitHub Actions: Upgrade Ubuntu + Go version and remove Caching
34adb1f Bump actions/setup-go from 3 to 4 (#144)
cf782c5 Bump actions/cache from 3.3.0 to 3.3.1 (#143)
3ee7c89 Bump actions/cache from 3.2.6 to 3.3.0 (#142)
d865180 Bump actions/cache from 3.2.5 to 3.2.6 (#140)
04e01d7 Upgrade static check to v2023.1 (#139)
7e1bcf3 go mod tidy (#137)
7e847d2 Bump actions/cache from 3.2.4 to 3.2.5 (#136)
15d9b44 Bump dominikh/staticcheck-action from 1.2.0 to 1.3.0 (#132)
14e2304 Bump actions/cache from 3.0.11 to 3.2.4 (#134)
952a13d Upgrade Go versions: Removing v1.18, adding v1.20
f726227 Bump actions/cache from 3.0.8 to 3.0.11 (#128)
41749a6 GitHub Actions: Raise Go versions to v1.18 and v1.19 (#125)
4fc9999 Bump actions/cache from 3.0.7 to 3.0.8 (#121)
b678d1c Bump actions/cache from 3.0.5 to 3.0.7 (#120)
4861c8c Bump actions/cache from 3.0.4 to 3.0.5 (#118)
28cf26f Bump actions/cache from 3.0.3 to 3.0.4 (#114)
2d5a93b Bump actions/cache from 3.0.2 to 3.0.3 (#113)
4bec308 Bump actions/cache from 3.0.1 to 3.0.2 (#108)
5f47610 Bump actions/setup-go from 2 to 3 (#109)
525eecd Bump dominikh/staticcheck-action from 1.1.0 to 1.2.0 (#107)
79adba8 Bump actions/checkout from 2 to 3 (#104)
6928c4f Bump actions/cache from 2 to 3.0.1 (#105)
90f99d2 Fix staticcheck in GitHub CI
5105eaf Switch CI to only support v1.17 und v1.18
916b31a Configured dependabot
58f949a go mod tidy
d813ef1 Run testing on a new pull request
1154c96 Testing: Use go version '1.17', '1.16', '1.15'
9d38b0b Bump github.com/google/go-querystring from 1.0.0 to 1.1.0 (#87)
58fd2fe Upgrade to GitHub-native Dependabot (#88)
e56a0c4 Rework CI: Removed travis + gometalinter, added GitHub Actions + staticcheck (#82)
9181c5d Fix order of package imports
78ea334 drop support for Go 1.9, add Go 1.12
174420e Remove "unused" from gometalinter
f48c3d1 Drop support for go v1.8
eefac78 gometalinter: Renamed linter gas to gosec
9e4f624 TravisCI: Added Go 1.11 to TravisCI build matrix
201f25d gometalinter: Rename gas security checker to gosec
c3ce3c2 Fix vet issue.
a9c12d7 Require Go 1.8 or newer.

#### Documentation

1fe64e5 Changes: Add documentation link to ChangeInput
027139b Improve documentation consistency.

#### README and examples

5a2c9c2 New example: Creating a project
024fc51 Add section about "Development" into README
3fc80f9 Smaller README rework and example directory
46815e4 README: Follow the Go Release Policy (#85)
8fc5ccb Rework README: Changed godoc links, Renamed Go(lang) to Go, added make commands (#84)

#### Other

c9ff84f Shorten paypal link for sponsoring
5dc4910 Add sponsoring information

### 0.5.2 (2017-11-04)

* Fix panic in checkAuth() if Gerrit is down #42
* Implement ListVotes(), DeleteVotes() and add missing tests

### 0.5.1 (2017-09-14)

* Added the `AbandonChange`, `RebaseChange`, `RestoreChange` and 
  `RevertChange` functions.

### 0.5.0 (2017-09-11)

**WARNING**: This release includes breaking changes.

* [BREAKING CHANGE] The SetReview function was returning the wrong
  entity type. (#40)

### 0.4.0 (2017-09-05)

**WARNING**: This release includes breaking changes.

* [BREAKING CHANGE] - Added gometalinter to the build and fixed problems 
  discovered by the linters.
    * Comment and error string fixes.
    * Numerous lint and styling fixes.
    * Ensured error values are being properly checked where appropriate.
    * Addition of missing documentation
    * Removed filePath parameter from DeleteChangeEdit which was unused and 
      unnecessary for the request.
    * Fixed CherryPickRevision and IncludeGroups functions which didn't pass
      along the provided input structs into the request.
* Go 1.5 has been removed from testing on Travis. The linters introduced in 
  0.4.0 do not support this version, Go 1.5 is lacking security updates and
  most Linux distros have moved beyond Go 1.5 now.
* Add Go 1.9 to the Travis matrix.
* Fixed an issue where urls containing certain characters in the credentials
  could cause NewClient() to use an invalid url. Something like `/`, which
  Gerrit could use for generated passwords, for example would break url.Parse's
  expectations.

### 0.3.0 (2017-06-05)

**WARNING**: This release includes breaking changes.

* [BREAKING CHANGE] Fix Changes.PublishDraftChange to accept a notify parameter.
* [BREAKING CHANGE] Fix PublishChangeEdit to accept a notify parameter.
* [BREAKING CHANGE] Fix ChangeFileContentInChangeEdit to allow the file content
  to be included in the request.
* Fix the url being used by CreateChange
* Fix type serialization of EventInfo.PatchSet.Number so it's consistent.
* Fix Changes.AddReviewer so it passes along the reviewer to the request.
* Simplify and optimize RemoveMagicPrefixLine

### 0.2.0 (2016-11-15)

**WARNING**: This release includes breaking changes.

* [BREAKING CHANGE] Several bugfixes to GetEvents:
  * Update EventInfo to handle the changeKey field and apply
    the proper type for the Project field
  * Provide a means to ignore marshaling errors
  * Update GetEvents() to return the failed lines and remove
    the pointer to the return value because it's unnecessary.
* [BREAKING CHANGE] In ec28f77 `ChangeInfo.Labels` has been changed to map
  to fix #21.

### 0.1.1 (2016-11-05)

* Minor fix to SubmitChange to use the `http.StatusConflict` constant
  instead of a hard coded value when comparing response codes.
* Updated AccountInfo.AccountID to be omitted of empty (such as when 
  used in ApprovalInfo).
* + and : in url parameters for queries are no longer escaped. This was
  causing `400 Bad Request` to be returned when the + symbol was
  included as part of the query. To match behavior with Gerrit's search
  handling, the : symbol was also excluded.
* Fixed documentation for NewClient and moved fmt.Errorf call from
  inside the function to a `ErrNoInstanceGiven` variable so it's
  easier to compare against.
* Updated internal function digestAuthHeader to return exported errors
  (ErrWWWAuthenticateHeader*) rather than calling fmt.Errorf. This makes
  it easier to test against externally and also fixes a lint issue too.
* Updated NewClient function to handle credentials in the url.
* Added the missing `Submitted` field to `ChangeInfo`.
* Added the missing `URL` field to `ChangeInfo` which is usually included
  as part of an event from the events-log plugin.

### 0.1.0 (2016-10-08)

* The first official release
* Implemented digest auth and several fixes for it.
* Ensured Content-Type is included in all requests
* Fixed several internal bugs as well as a few documentation issues
