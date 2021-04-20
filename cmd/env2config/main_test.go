package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

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

func TestRun(t *testing.T) {
	dir := t.TempDir()
	tempYaml := filepath.Join(dir, "some.yaml")
	setEnv(t, "E2C_CONFIGS", "myprefix")
	setEnv(t, "MYPREFIX_OPTS_FILE", tempYaml)
	setEnv(t, "MYPREFIX_OPTS_FORMAT", "yaml")
	setEnv(t, "MYPREFIX_FOO", "bar")
	setEnv(t, "MYPREFIX_bAz0", "bit")
	setEnv(t, "MYPREFIX_bar.nested.key", "value")
	setEnv(t, "MYPREFIX_bar.array.0", "value")
	setEnv(t, "MYPREFIX_bar.array.1.key", "value")

	assert.NoError(t, run(nil))
	buf, err := ioutil.ReadFile(tempYaml)
	require.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(`
FOO: bar
bAz0: bit
bar:
    array:
        - value
        - key: value
    nested:
        key: value
`)+"\n", string(buf))
}

func TestRunTemplate(t *testing.T) {
	dir := t.TempDir()
	tmpYaml := filepath.Join(dir, "some.yaml")
	templateYaml := filepath.Join(dir, "template.yaml")
	setEnv(t, "E2C_CONFIGS", "myprefix")
	setEnv(t, "MYPREFIX_OPTS_FILE", tmpYaml)
	setEnv(t, "MYPREFIX_OPTS_FORMAT", "yaml")
	setEnv(t, "MYPREFIX_OPTS_TEMPLATE_FILE", templateYaml)
	setEnv(t, "MYPREFIX_FOO", "bar")
	setEnv(t, "MYPREFIX_bAz0", "bit")
	setEnv(t, "MYPREFIX_bar.nested.key", "value")
	setEnv(t, "MYPREFIX_bar.array.0.other_key", "value")
	setEnv(t, "MYPREFIX_bar.array.1.key", "value")
	setEnv(t, "MYPREFIX_bar.array.2", "value")

	require.NoError(t, ioutil.WriteFile(templateYaml, []byte(strings.TrimSpace(`
FOO: not bar
bar:
    array:
        - key: some default
    default_key: default_value
no_value:
`)), 0600))

	assert.NoError(t, run(nil))
	buf, err := ioutil.ReadFile(tmpYaml)
	require.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(`
FOO: bar
bAz0: bit
bar:
    array:
        - key: some default
          other_key: value
        - key: value
        - value
    default_key: default_value
    nested:
        key: value
no_value:
`)+"\n", string(buf))
}

func TestRunInputs(t *testing.T) {
	t.Run("input missing", func(t *testing.T) {
		dir := t.TempDir()
		tempYaml := filepath.Join(dir, "some.yaml")
		setEnv(t, "E2C_CONFIGS", "myprefix")
		setEnv(t, "MYPREFIX_OPTS_FILE", tempYaml)
		setEnv(t, "MYPREFIX_OPTS_FORMAT", "yaml")
		setEnv(t, "MYPREFIX_OPTS_IN_url", "REQUIRED_URL")
		assert.EqualError(t, run(nil), `Environment variable "REQUIRED_URL" is required`)
	})

	t.Run("input provided", func(t *testing.T) {
		dir := t.TempDir()
		tempYaml := filepath.Join(dir, "some.yaml")
		setEnv(t, "E2C_CONFIGS", "myprefix")
		setEnv(t, "MYPREFIX_OPTS_FILE", tempYaml)
		setEnv(t, "MYPREFIX_OPTS_FORMAT", "yaml")
		setEnv(t, "MYPREFIX_OPTS_IN_url", "REQUIRED_URL")
		setEnv(t, "REQUIRED_URL", "http://localhost/")

		assert.NoError(t, run(nil))
		buf, err := ioutil.ReadFile(tempYaml)
		require.NoError(t, err)
		assert.Equal(t, strings.TrimSpace(`
url: http://localhost/
`)+"\n", string(buf))
	})
}
