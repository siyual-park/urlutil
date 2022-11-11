package util

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		whenSource   any
		whenKey      string
		expectResult any
		expectOk     bool
	}{
		{
			whenSource:   map[string]any{"k1": map[string]any{"k2": 1}},
			whenKey:      "k1.k2",
			expectResult: 1,
			expectOk:     true,
		},
		{
			whenSource:   map[string]any{"k1": []map[string]any{{"k2": 1}}},
			whenKey:      "k1[0].k2",
			expectResult: 1,
			expectOk:     true,
		},
		{
			whenSource: map[string]any{"k1": func() *sync.Map {
				m := sync.Map{}
				m.Store("k2", 1)
				return &m
			}()},
			whenKey:      "k1.k2",
			expectResult: 1,
			expectOk:     true,
		},
		{
			whenSource: map[string]any{"k1": struct {
				K2 int
			}{
				K2: 1,
			}},
			whenKey:      "k1.k2",
			expectResult: 1,
			expectOk:     true,
		},
	}

	for _, tc := range testCases {
		res, ok := Get[any](tc.whenSource, tc.whenKey)
		assert.Equal(t, tc.expectOk, ok)
		if ok {
			assert.Equal(t, tc.expectResult, res)
		}
	}
}

func TestSet(t *testing.T) {
	testCases := []struct {
		whenSource any
		whenKey    string
		whenValue  any
		expectOk   bool
	}{
		{
			whenSource: map[string]any{"k1": map[string]any{"k2": 1}},
			whenKey:    "k1.k2",
			whenValue:  2,
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": map[string]any{}},
			whenKey:    "k1.k2",
			whenValue:  2,
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": &map[string]any{}},
			whenKey:    "k1.k2",
			whenValue:  2,
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": []map[string]any{{"k2": 1}}},
			whenKey:    "k1[0].k2",
			whenValue:  2,
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": []map[string]any{{"k2": 0}}},
			whenKey:    "k1[0]",
			whenValue:  map[string]any{"k2": 1},
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": &sync.Map{}},
			whenKey:    "k1.k2",
			whenValue:  2,
			expectOk:   true,
		},
		{
			whenSource: map[string]any{"k1": &struct {
				K2 int
			}{
				K2: 1,
			}},
			whenKey:   "k1.k2",
			whenValue: 2,
			expectOk:  true,
		},
	}

	for _, tc := range testCases {
		ok := Set(&tc.whenSource, tc.whenKey, tc.whenValue)
		assert.Equal(t, tc.expectOk, ok)
		if ok {
			res, ok := Get[any](tc.whenSource, tc.whenKey)
			assert.Equal(t, tc.expectOk, ok)
			assert.Equal(t, tc.whenValue, res)
		}
	}
}
