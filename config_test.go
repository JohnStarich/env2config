package env2config

import (
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setEnv(t *testing.T, key, value string) {
	prev, exists := os.LookupEnv(key)
	require.NoError(t, os.Setenv(key, value))
	t.Cleanup(func() {
		if exists {
			require.NoError(t, os.Setenv(key, prev))
		} else {
			require.NoError(t, os.Unsetenv(key))
		}
	})
}

type gorpMarshaler struct {
	marshaledValue  interface{}
	marshalErr      error
	unmarshalResult map[string]interface{}
	unmarshalErr    error
}

func (g *gorpMarshaler) Marshal(w io.Writer, value interface{}) error {
	g.marshaledValue = value
	return g.marshalErr
}

func (g *gorpMarshaler) Unmarshal(r io.Reader, value interface{}) error {
	*value.(*map[string]interface{}) = g.unmarshalResult
	return g.unmarshalErr
}

func TestNew(t *testing.T) {
	RegisterFormat("gorp", &gorpMarshaler{})
	setEnv(t, "MYPREFIX_OPTS_FILE", "/some/path.gorp")
	setEnv(t, "MYPREFIX_OPTS_FORMAT", "gorp")

	t.Run("happy path", func(t *testing.T) {
		setEnv(t, "MYPREFIX_FOO", "bar")
		setEnv(t, "MYPREFIX_bAz0", "bit")
		setEnv(t, "MYPREFIX_OPTS_IN_port", "BIND_PORT")
		setEnv(t, "BIND_PORT", "8080")

		config, err := New("MYPREFIX")
		assert.Equal(t, Config{
			Name: "myprefix",
			Opts: Opts{
				File:   "/some/path.gorp",
				Format: "gorp",
				Inputs: map[string]string{
					"port": "BIND_PORT",
				},
			},
			Values: map[string]string{
				"FOO":  "bar",
				"bAz0": "bit",
				"port": "8080",
			},
			registry: defaultRegistry,
		}, config)
		assert.NoError(t, err)
	})

	t.Run("missing required var", func(t *testing.T) {
		setEnv(t, "MYPREFIX_OPTS_IN_port", "BIND_PORT")
		_, err := New("MYPREFIX")
		assert.EqualError(t, err, `Missing required environment variables: BIND_PORT`)
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
	templateFile := filepath.Join(dir, "template.gorp")
	require.NoError(t, ioutil.WriteFile(templateFile, nil, 0600))

	for _, tc := range []struct {
		description     string
		config          Config
		unmarshalResult map[string]interface{}
		expectMarshal   interface{}
		marshalErr      error
		expectErr       string
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
		{
			description: "template file defaults",
			config: Config{
				Opts: Opts{
					Format:       "gorp",
					File:         tempFile,
					TemplateFile: templateFile,
				},
				Values: map[string]string{
					"A": "B",
				},
			},
			unmarshalResult: map[string]interface{}{
				"C": "D",
			},
			expectMarshal: map[string]interface{}{
				"A": "B",
				"C": "D",
			},
		},
		{
			description: "nested template fields",
			config: Config{
				Opts: Opts{
					Format:       "gorp",
					File:         tempFile,
					TemplateFile: templateFile,
				},
				Values: map[string]string{
					"A.B.C": "1",
					"F.0.I": "2",
				},
			},
			unmarshalResult: map[string]interface{}{
				"A": map[string]interface{}{
					"B": map[string]interface{}{
						"C": "D",
						"E": "F",
					},
				},
				"F": []interface{}{
					map[string]interface{}{
						"G": "H",
						"I": "J",
					},
					"K",
				},
			},
			expectMarshal: map[string]interface{}{
				"A": map[string]interface{}{
					"B": map[string]interface{}{
						"C": "1",
						"E": "F",
					},
				},
				"F": []interface{}{
					map[string]interface{}{
						"G": "H",
						"I": "2",
					},
					"K",
				},
			},
		},
		{
			description: "nested key array paths",
			config: Config{
				Opts: Opts{Format: "gorp", File: tempFile},
				Values: map[string]string{
					"A.B.0.C": "1",
				},
			},
			expectMarshal: map[string]interface{}{
				"A": map[string]interface{}{
					"B": []interface{}{
						map[string]interface{}{
							"C": "1",
						},
					},
				},
			},
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			marshaler := &gorpMarshaler{
				marshalErr:      tc.marshalErr,
				unmarshalResult: tc.unmarshalResult,
			}
			tc.config.registry = newRegistry()
			tc.config.registry.RegisterFormat("gorp", marshaler)

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
