# WebLoop

Headless WebKit with a Go API. Inspired by [PhantomJS](http://phantomjs.org/).

* [Documentation on Sourcegraph](https://sourcegraph.com/github.com/sourcegraph/webloop/tree)

[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/webloop/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/webloop)
[![xrefs](https://sourcegraph.com/api/repos/github.com/sourcegraph/webloop/badges/xrefs.png)](https://sourcegraph.com/github.com/sourcegraph/webloop)
[![funcs](https://sourcegraph.com/api/repos/github.com/sourcegraph/webloop/badges/funcs.png)](https://sourcegraph.com/github.com/sourcegraph/webloop)
[![top func](https://sourcegraph.com/api/repos/github.com/sourcegraph/webloop/badges/top-func.png)](https://sourcegraph.com/github.com/sourcegraph/webloop)
[![library users](https://sourcegraph.com/api/repos/github.com/sourcegraph/webloop/badges/library-users.png)](https://sourcegraph.com/github.com/sourcegraph/webloop)


## Requirements

* [Go](http://golang.org) >= 1.2rc1 (due to [#3250](https://code.google.com/p/go/issues/detail?id=3250))
* [WebKitGTK+](http://webkitgtk.org/) >= 2.0.0
* [go-webkit2](https://sourcegraph.com/github.com/sourcegraph/go-webkit2/readme)

For instructions on installing these dependencies, see the [go-webkit2
README](https://sourcegraph.com/github.com/sourcegraph/go-webkit2/readme).


## Usage

### Example: rendering static HTML from a dynamic, single-page [AngularJS](http://angularjs.org) app

See the `examples/angular-static-seo/` directory for example code. Run the included binary with:

```
go run examples/angular-static-seo/server.go
```

Instructions will be printed for accessing the 2 local demo HTTP servers. Run
with `-h` to see more information.


## TODO

* [Set up CI testing.](https://github.com/sourcegraph/webloop/issues/1) This
  is difficult because all of the popular CI services run older versions of
  Ubuntu that make it difficult to install WebKitGTK+ >= 2.0.0.
* Add the ability for JavaScript code to send messages to WebLoop, similar to
  [PhantomJS's callPhantom]
  (https://github.com/ariya/phantomjs/wiki/API-Reference-WebPage#oncallback)
  mechanism.


## Contributors

See the AUTHORS file for a list of contributors.

Submit contributions via GitHub pull request. Patches should include tests and
should pass [golint](https://github.com/golang/lint).
