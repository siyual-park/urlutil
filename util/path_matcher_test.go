package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatchPath(t *testing.T) {
	testCases := []struct {
		whenPath    string
		whenPattern string
		expectParam map[string]string
	}{
		{
			whenPath:    "/static",
			whenPattern: "/static",
			expectParam: map[string]string{},
		},
		{
			whenPath:    "/static/any",
			whenPattern: "/static/*",
			expectParam: map[string]string{"*": "any"},
		},
		{
			whenPath:    "/params/1",
			whenPattern: "/params/:foo",
			expectParam: map[string]string{"foo": "1"},
		},
		{
			whenPath:    "/params/1/bar/2",
			whenPattern: "/params/:foo/bar/:qux",
			expectParam: map[string]string{"foo": "1", "qux": "2"},
		},
		{
			whenPath:    "/params/1/bar/2/any",
			whenPattern: "/params/:foo/bar/:qux/*",
			expectParam: map[string]string{"foo": "1", "qux": "2", "*": "any"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.whenPath, func(t *testing.T) {
			ok, param := MatchPath(tc.whenPattern, tc.whenPath)
			assert.True(t, ok)
			assert.Equal(t, tc.expectParam, param)
		})
	}
}

func TestPathMatcher_Match(t *testing.T) {
	m := NewPathMatcher()

	testCases := []struct {
		whenPath    string
		expectPath  string
		expectParam map[string]string
	}{
		{
			whenPath:    "/static",
			expectPath:  "/static",
			expectParam: map[string]string{},
		},
		{
			whenPath:    "/static/any",
			expectPath:  "/static/*",
			expectParam: map[string]string{"*": "any"},
		},
		{
			whenPath:    "/params/1",
			expectPath:  "/params/:foo",
			expectParam: map[string]string{"foo": "1"},
		},
		{
			whenPath:    "/params/1/bar/2",
			expectPath:  "/params/:foo/bar/:qux",
			expectParam: map[string]string{"foo": "1", "qux": "2"},
		},
		{
			whenPath:    "/params/1/bar/2/any",
			expectPath:  "/params/:foo/bar/:qux/*",
			expectParam: map[string]string{"foo": "1", "qux": "2", "*": "any"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.whenPath, func(t *testing.T) {
			m.Add(tc.expectPath)
			path, param := m.Match(tc.whenPath)
			assert.Equal(t, tc.expectPath, path)
			assert.Equal(t, tc.expectParam, param)
		})
	}
}

func TestPathMatcher_Remove(t *testing.T) {
	m := NewPathMatcher()

	testCases := []struct {
		whenPath string
	}{
		{
			whenPath: "/static",
		},
		{
			whenPath: "/static/*",
		},
		{
			whenPath: "/params/:foo",
		},
		{
			whenPath: "/params/:foo/bar/:qux",
		},
		{
			whenPath: "/params/:foo/bar/:qux/*",
		},
		{
			whenPath: "/params/:foo/:bar/:qux",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.whenPath, func(t *testing.T) {
			m.Add(tc.whenPath)
			assert.True(t, m.Remove(tc.whenPath))
		})
	}
}
