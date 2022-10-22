# urlutil
Package urlutil provides URL utility functions, complementing the more common ones in the url package.

## Install
```shell
go get github.com/siyual-park/urlutil
```

## Basic Example
### Matcher
Parse the path, and find the best candidate path.
```go
matcher := urlutil.NewMatcher()

matcher.Add("/static")
matcher.Add("/static/*")
matcher.Add("/params/:foo")

matcher.Match("/static") // "/static", map[string]string{}
matcher.Match("/static/any") // "/static/*", map[string]string{"*": "any"}
matcher.Match("/params/1") // "/params/:foo", map[string]string{"foo": "1"}
``` 

Match only one path.
```go
urlutil.Match("/params/:foo", "/params/1") // true, map[string]string{"foo": "1"}
``` 

## Special thanks
Some code for this package was taken from https://github.com/labstack/echo
