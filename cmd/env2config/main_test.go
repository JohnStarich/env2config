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

func TestRun(t *testing.T) {
	dir := t.TempDir()
	tempYaml := filepath.Join(dir, "some.yaml")
	require.NoError(t, os.Setenv("E2C_CONFIGS", "myprefix"))
	require.NoError(t, os.Setenv("MYPREFIX_OPTS_FILE", tempYaml))
	require.NoError(t, os.Setenv("MYPREFIX_OPTS_FORMAT", "yaml"))
	require.NoError(t, os.Setenv("MYPREFIX_FOO", "bar"))
	require.NoError(t, os.Setenv("MYPREFIX_bAz0", "bit"))
	require.NoError(t, os.Setenv("MYPREFIX_bar.nested.key", "value"))
	require.NoError(t, os.Setenv("MYPREFIX_bar.array.0", "value"))
	require.NoError(t, os.Setenv("MYPREFIX_bar.array.1.key", "value"))

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
