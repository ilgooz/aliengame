package aliengamecmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const testmapPath = "../../../mapdata/0.aliengame"

func TestAlienCmd(t *testing.T) {
	cmd := New()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"-m", testmapPath, "-a", "3"})
	require.NoError(t, cmd.Execute())
	output := buf.String()
	for _, word := range []string{
		"Foo",
		"Bar",
		"Baz",
		"Qu-ux",
		"Yee",
		"MAP STATE",
	} {
		require.True(t, strings.Contains(output, word), word)
	}
}
