# Go Utilities
Reusable collection of golang utility.

## Install
```shell
go get github.com/siyual-park/go-util/util
```


## Basic Example
### Path Matcher
Parse the path, and find the best candidate path.
```go
matcher := util.NewPathMatcher()

matcher.Add("/static")
matcher.Add("/static/*")
matcher.Add("/params/:foo")

matcher.Match("/static") // "/static", map[string]string{}
matcher.Match("/static/any") // "/static/*", map[string]string{"*": "any"}
matcher.Match("/params/1") // "/params/:foo", map[string]string{"foo": "1"}
``` 

Match only one path.
```go
util.MatchPath("/params/:foo", "/params/1") // true, map[string]string{"foo": "1"}
``` 

#### Special thanks
Some code for this package was taken from https://github.com/labstack/echo
