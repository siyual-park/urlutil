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

### Pointer
Helps convert the value of Pointer.

```go
ptr := util.Ptr("any_string")
assert.Equal(t, "any_string", *ptr)
``` 

```go
ptr := util.Ptr("any_string")
assert.Equal(t, "any_string", util.UnPtr(ptr))

assert.Equal(t, "", util.UnPtr[string](nil))
```

### Access
It gives you smart access to any type of value.

#### rules
- find getter in interface. (when key name is "name", getter method name is "Name")
- find map like access method. (Get(name string), Load(name string))
- find public struct property.
- find slice or array property.
- find map property.
- un-pointer and research.

#### Get
```go
v := map[string]any{"k1": map[string]any{"k2": 1}}
r, ok := util.Get[int](v, "k1.k2")
assert.True(t, ok)
assert.Equal(t, 1, r)
``` 

```go
v := map[string]any{"k1": func() *sync.Map {
    m := sync.Map{}
    m.Store("k2", 1)
    return &m
}()}
r, ok := util.Get[int](v, "k1.k2")
assert.True(t, ok)
assert.Equal(t, 1, r)
```

#### Set
```go
v := map[string]any{"k1": map[string]any{}}
ok := util.Set[int](v, "k1.k2", 1)
assert.True(t, ok)
``` 

### Iterator
#### KeyTo
```go
v := map[string]int{"a": 1, "b": 2}
r := util.KeyTo(v, func(key any) any {
    if k, ok := key.(string); ok {
        return "_" + k
    }
    return key
})
assert.Equal(t, map[string]int{"_a": 1, "_b": 2}, r)
``` 

#### ValueTo
```go
v := map[string]int{"a": 1, "b": 2}
r := util.ValueTo(v, func(value any) any {
    if v, ok := value.(int); ok {
        return v + 1
    }
    return value
})
assert.Equal(t, map[string]int{"a": 2, "b": 3}, r)
``` 

### Copy
#### DeepCopy
```go
v := map[string]int{"a": 1, "b": 2}
r := util.DeepCopy(v)
assert.Equal(t, v, r)
``` 
