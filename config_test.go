package env2config

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type gorpMarshaler struct {
	marshaledValue interface{}
	marshalErr     error
}

func (g *gorpMarshaler) Marshal(w io.Writer, value interface{}) error {
	g.marshaledValue = value
	return g.marshalErr
}

func TestNew(t *testing.T) {
	RegisterFormat("gorp", &gorpMarshaler{})
	require.NoError(t, os.Setenv("MYPREFIX_OPTS_FILE", "/some/path.gorp"))
	require.NoError(t, os.Setenv("MYPREFIX_OPTS_FORMAT", "gorp"))
	require.NoError(t, os.Setenv("MYPREFIX_FOO", "bar"))
	require.NoError(t, os.Setenv("MYPREFIX_bAz0", "bit"))

	t.Run("happy path", func(t *testing.T) {
		config, err := New("MYPREFIX")
		assert.Equal(t, Config{
			Name: "myprefix",
			Opts: Opts{
				File:   "/some/path.gorp",
				Format: "gorp",
			},
			Values: map[string]string{
				"FOO":  "bar",
				"bAz0": "bit",
			},
			registry: defaultRegistry,
		}, config)
		assert.NoError(t, err)
	})

	t.Run("missing config name", func(t *testing.T) {
		_, err := New("")
		assert.EqualError(t, err, "Config name is required")
	})

	t.Run("bad config name", func(t *testing.T) {
		_, err := New("not a valid name")
		assert.EqualError(t, err, `Config names must only use letters or numbers: "not a valid name"`)
	})
}

func TestWrite(t *testing.T) {
	dir := t.TempDir()
	tempFile := filepath.Join(dir, "temp.gorp")

	for _, tc := range []struct {
		description   string
		config        Config
		expectMarshal interface{}
		marshalErr    error
		expectErr     string
	}{
		{
			description: "one pair",
			config: Config{
				Opts: Opts{Format: "gorp", File: tempFile},
				Values: map[string]string{
					"A": "B",
				},
			},
			expectMarshal: map[string]interface{}{
				"A": "B",
			},
		},
		{
			description: "nested key paths",
			config: Config{
				Opts: Opts{Format: "gorp", File: tempFile},
				Values: map[string]string{
					"A":               "B",
					"c.d.some_subkey": "some value",
				},
			},
			expectMarshal: map[string]interface{}{
				"A": "B",
				"c": map[string]interface{}{
					"d": map[string]interface{}{
						"some_subkey": "some value",
					},
				},
			},
		},
		{
			description: "array paths",
			config: Config{
				Opts: Opts{Format: "gorp", File: tempFile},
				Values: map[string]string{
					"A.1.B": "D",
					"A.0.B": "C",
				},
			},
			expectMarshal: map[string]interface{}{
				"A": []interface{}{
					map[string]interface{}{"B": "C"},
					map[string]interface{}{"B": "D"},
				},
			},
		},
		{
			description: "marshal err",
			config: Config{
				Opts: Opts{Format: "gorp", File: tempFile},
			},
			marshalErr: errors.New("some error"),
			expectErr:  "some error",
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			marshaler := &gorpMarshaler{marshalErr: tc.marshalErr}
			tc.config.registry = &registry{
				marshalers: map[string]Marshaler{
					"gorp": marshaler,
				},
			}

			err := tc.config.Write()
			if tc.expectErr != "" {
				assert.EqualError(t, err, tc.expectErr)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tc.expectMarshal, marshaler.marshaledValue)
		})
	}
}
