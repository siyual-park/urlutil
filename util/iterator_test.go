package util

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeyTo(t *testing.T) {
	convert := func(key any) any {
		if k, ok := key.(string); ok {
			return "_" + k
		}
		return key
	}

	var testCase = []struct {
		when   any
		expect any
	}{
		{
			when:   map[string]int{"a": 1, "b": 2},
			expect: map[string]int{"_a": 1, "_b": 2},
		},
		{
			when:   map[string]any{"a": 1, "b": 2, "c": map[string]any{"a": 1, "b": 2}},
			expect: map[string]any{"_a": 1, "_b": 2, "_c": map[string]any{"_a": 1, "_b": 2}},
		},
		{
			when:   []map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: []map[string]int{{"_a": 1, "_b": 2}, {"_c": 1, "_d": 2}},
		},
		{
			when:   [2]map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: [2]map[string]int{{"_a": 1, "_b": 2}, {"_c": 1, "_d": 2}},
		},
		{
			when:   &[]map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: &[]map[string]int{{"_a": 1, "_b": 2}, {"_c": 1, "_d": 2}},
		},
		{
			when:   &map[string]int{"a": 1, "b": 2},
			expect: &map[string]int{"_a": 1, "_b": 2},
		},
	}

	for _, tc := range testCase {
		assert.Equal(t, tc.expect, KeyTo(tc.when, convert))
	}
}

func TestValueTo(t *testing.T) {
	convert := func(value any) any {
		if v, ok := value.(int); ok {
			return v + 1
		}
		return value
	}

	var testCase = []struct {
		when   any
		expect any
	}{
		{
			when:   map[string]int{"a": 1, "b": 2},
			expect: map[string]int{"a": 2, "b": 3},
		},
		{
			when:   map[string]any{"a": 1, "b": 2, "c": map[string]any{"a": 1, "b": 2}},
			expect: map[string]any{"a": 2, "b": 3, "c": map[string]any{"a": 2, "b": 3}},
		},
		{
			when:   []map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: []map[string]int{{"a": 2, "b": 3}, {"c": 2, "d": 3}},
		},
		{
			when:   [2]map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: [2]map[string]int{{"a": 2, "b": 3}, {"c": 2, "d": 3}},
		},
		{
			when:   &[]map[string]int{{"a": 1, "b": 2}, {"c": 1, "d": 2}},
			expect: &[]map[string]int{{"a": 2, "b": 3}, {"c": 2, "d": 3}},
		},
		{
			when:   &map[string]int{"a": 1, "b": 2},
			expect: &map[string]int{"a": 2, "b": 3},
		},
	}

	for _, tc := range testCase {
		assert.Equal(t, tc.expect, ValueTo(tc.when, convert))
	}
}
