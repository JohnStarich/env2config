package env2config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseKeyPath(t *testing.T) {
	for _, tc := range []struct {
		input  string
		expect []string
	}{
		{
			input:  `x`,
			expect: []string{`x`},
		},
		{
			input:  `x.y.z`,
			expect: []string{`x`, `y`, `z`},
		},
		{
			input:  `x\.y.z`,
			expect: []string{`x.y`, `z`},
		},
	} {
		t.Run(tc.input, func(t *testing.T) {
			assert.Equal(t, tc.expect, parseKeyPath(tc.input))
		})
	}
}

func TestDeleteKeyPath(t *testing.T) {
	for _, tc := range []struct {
		description  string
		key          string
		input        interface{}
		expect       interface{}
		expectDelete bool
	}{
		{
			description:  "base case",
			key:          "",
			input:        42,
			expect:       nil,
			expectDelete: true,
		},
		{
			description: "delete non existent key",
			key:         "foo",
			input:       map[string]interface{}{"bar": 1},
			expect:      map[string]interface{}{"bar": 1},
		},
		{
			description: "delete map key",
			key:         "foo",
			input:       map[string]interface{}{"bar": 1, "foo": 2},
			expect:      map[string]interface{}{"bar": 1},
		},
		{
			description: "delete array index",
			key:         "1",
			input:       []interface{}{"bar", "foo"},
			expect:      []interface{}{"bar"},
		},
		{
			description: "delete non-existent array index",
			key:         "10",
			input:       []interface{}{"bar", "foo"},
			expect:      []interface{}{"bar", "foo"},
		},
		{
			description: "delete array index equal to length",
			key:         "2",
			input:       []interface{}{"bar", "foo"},
			expect:      []interface{}{"bar", "foo"},
		},
		{
			description: "delete nested keys",
			key:         "foo.bar.1.baz",
			input: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": []interface{}{
						10,
						map[string]interface{}{
							"baz":  42,
							"biff": "six times seven",
						},
					},
				},
				"bag": 1,
			},
			expect: map[string]interface{}{
				"foo": map[string]interface{}{
					"bar": []interface{}{
						10,
						map[string]interface{}{
							"biff": "six times seven",
						},
					},
				},
				"bag": 1,
			},
		},
		{
			description: "delete unexpected type",
			key:         "foo.bar",
			input:       "baz",
			expect:      "baz",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			keys := parseKeyPath(tc.key)
			output, deleteMe := deleteKeyPath(tc.input, keys)
			assert.Equal(t, tc.expectDelete, deleteMe)
			assert.Equal(t, tc.expect, output)
		})
	}
}

func TestSortTemplateDeleteKeys(t *testing.T) {
	sample1 := func() interface{} {
		return map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": []interface{}{
					"baz",
					"biff",
				},
			},
		}
	}
	sample1DeleteBazBiff := func() interface{} {
		return map[string]interface{}{
			"foo": map[string]interface{}{
				"bar": []interface{}{},
			},
		}
	}
	sample2 := func() interface{} {
		return map[string]interface{}{
			"foo": []interface{}{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
		}
	}
	sample2Delete10thItems := func() interface{} {
		return map[string]interface{}{
			"foo": []interface{}{0, 1, 3, 4, 5, 6, 7, 8, 9},
		}
	}
	sample3 := func() interface{} {
		return map[string]interface{}{
			"foo": []interface{}{
				"biff",
				map[string]interface{}{
					"bar": "baz",
				},
				"beep",
			},
		}
	}
	sample3DeleteDifferentLengthPaths := func() interface{} {
		return map[string]interface{}{
			"foo": []interface{}{
				"biff",
				"beep",
			},
		}
	}
	for _, tc := range []struct {
		description   string
		keys          []string
		testContent   interface{}
		expectKeys    []string
		expectContent interface{}
	}{
		{
			description:   "no keys",
			keys:          nil,
			testContent:   sample1(),
			expectKeys:    nil,
			expectContent: sample1(),
		},
		{
			description:   "delete both array items in order",
			keys:          []string{"foo.bar.0", "foo.bar.1"},
			testContent:   sample1(),
			expectKeys:    []string{"foo.bar.1", "foo.bar.0"},
			expectContent: sample1DeleteBazBiff(),
		},
		{
			description:   "delete both array items out of order",
			keys:          []string{"foo.bar.1", "foo.bar.0"},
			testContent:   sample1(),
			expectKeys:    []string{"foo.bar.1", "foo.bar.0"},
			expectContent: sample1DeleteBazBiff(),
		},
		{
			description:   "delete multi-digit arrays indexes",
			keys:          []string{"foo.2", "foo.10"},
			testContent:   sample2(),
			expectKeys:    []string{"foo.10", "foo.2"},
			expectContent: sample2Delete10thItems(),
		},
		{
			description: "delete array indices via keys with different path lengths",
			// if foo.1 is deleted first, then foo.1.bar will delete a sub-field of a totally different array item
			keys:          []string{"foo.1", "foo.1.bar"},
			testContent:   sample3(),
			expectKeys:    []string{"foo.1.bar", "foo.1"},
			expectContent: sample3DeleteDifferentLengthPaths(),
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			content := tc.testContent
			sortTemplateDeleteKeys(tc.keys)
			assert.Equal(t, tc.expectKeys, tc.keys)
			for _, key := range tc.keys {
				keyPath := parseKeyPath(key)
				var deleteMe bool
				content, deleteMe = deleteKeyPath(content, keyPath)
				assert.False(t, deleteMe)
			}
			assert.Equal(t, tc.expectContent, content)
		})
	}
}
