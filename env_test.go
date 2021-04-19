package env2config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEnv(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		description string
		envPairs    []string
		expect      map[string]string
		expectPanic string
	}{
		{
			description: "no env",
			expect:      map[string]string{},
		},
		{
			description: "one pair",
			envPairs: []string{
				"A=B",
			},
			expect: map[string]string{
				"A": "B",
			},
		},
		{
			description: "permits empty",
			envPairs: []string{
				"A=",
			},
			expect: map[string]string{
				"A": "",
			},
		},
		{
			description: "permits non-alphanum",
			envPairs: []string{
				"A.B0=C",
			},
			expect: map[string]string{
				"A.B0": "C",
			},
		},
		{
			description: "duplicate pairs keeps first",
			envPairs: []string{
				"A=B",
				"A=D",
			},
			expect: map[string]string{
				"A": "B",
			},
		},
		{
			description: "many pairs",
			envPairs: []string{
				"A=B",
				"C=D",
			},
			expect: map[string]string{
				"A": "B",
				"C": "D",
			},
		},
		{
			description: "invalid env pair",
			envPairs: []string{
				"A",
			},
			expectPanic: "Invalid key value pair from os.Environ()",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			if tc.expectPanic != "" {
				assert.PanicsWithValue(t, tc.expectPanic, func() {
					parseEnv(tc.envPairs)
				})
				return
			}
			assert.Equal(t, tc.expect, parseEnv(tc.envPairs))
		})
	}
}

func TestConfigValues(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		description string
		configName  string
		env         map[string]string
		expect      map[string]string
	}{
		{
			description: "no env",
			configName:  "prefix",
			expect:      map[string]string{},
		},
		{
			description: "prefixed config vars",
			configName:  "prefix",
			env: map[string]string{
				"PREFIX_A":     "B",
				"PREFIX_C_D_E": "F",
			},
			expect: map[string]string{
				"A":     "B",
				"C_D_E": "F",
			},
		},
		{
			description: "ignore other vars",
			configName:  "prefix",
			env: map[string]string{
				"PREFIX_A": "B",
				"OTHER_C":  "D",
			},
			expect: map[string]string{
				"A": "B",
			},
		},
		{
			description: "ignore config opts",
			configName:  "prefix",
			env: map[string]string{
				"PREFIX_A":      "B",
				"PREFIX_OPTS_X": "D",
			},
			expect: map[string]string{
				"A": "B",
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			assert.Equal(t, tc.expect, configValues(tc.configName, tc.env))
		})
	}
}
