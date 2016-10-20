# Changelog

This is a high level log of changes, bugfixes, enhancements, etc
that have taken place between releases. Later versions are shown
first. For more complete details see
[the releases on GitHub.](https://github.com/andygrunwald/go-gerrit/releases)

## Versions

### 0.1.1 (not yet released)

* Minor fix to SubmitChange to use the `http.StatusConflict` constant
  instead of a hard coded value when comparing response codes.
* Updated AccountInfo.AccountID to be omitted of empty (such as when 
  used in ApprovalInfo).
* Fixed documentation for NewClient and moved fmt.Errorf call from
  inside the function to a `ErrNoInstanceGiven` variable so it's
  easier to compare against.
* Updated internal function digestAuthHeader to return exported errors
  (ErrWWWAuthenticateHeader*) rather than calling fmt.Errorf. This makes
  it easier to test against externally and also fixes a lint issue too.
* Updated NewClient function to handle credentials in the url.

### 0.1.0

* The first official release
* Implemented digest auth and several fixes for it.
* Ensured Content-Type is included in all requests
* Fixed several internal bugs as well as a few documentation issues
