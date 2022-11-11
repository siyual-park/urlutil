package util

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestDeepCopy(t *testing.T) {
	testCases := []struct {
		when any
	}{
		{
			when: map[string]any{"k1": map[string]any{"k2": 1}},
		},
		{
			when: map[string]any{"k1": []map[string]any{{"k2": 1}}},
		},
		{
			when: map[string]any{"k1": func() *sync.Map {
				m := sync.Map{}
				m.Store("k2", 1)
				return &m
			}()},
		},
		{
			when: map[string]any{"k1": struct {
				K2 int
			}{
				K2: 1,
			}},
		},
	}

	for _, tc := range testCases {
		res := DeepCopy(tc.when)
		assert.Equal(t, tc.when, res)
	}
}
