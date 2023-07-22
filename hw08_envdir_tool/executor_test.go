package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const checkTemplate = `
if [ "$%s" == "%s" ]; then
	exit 0;
else
	exit 1;
fi;
`

func TestRunCmd(t *testing.T) {
	testEnv, err := ReadDir("testdata/env")
	require.NoError(t, err)

	testData := []struct {
		name  string
		value string
	}{
		{"FOO", "   foo\nwith new line"},
		{"BAR", "bar"},
		{"EMPTY", ""},
		{"HELLO", "\\\"hello\\\""},
	}

	t.Run("set env", func(t *testing.T) {
		for _, data := range testData {
			checkCommand := fmt.Sprintf(checkTemplate, data.name, data.value)
			result := RunCmd([]string{"bash", "-c", checkCommand}, testEnv)
			require.Equal(t, 0, result)
		}
	})

	t.Run("unset env", func(t *testing.T) {
		os.Setenv("UNSET", "test")
		RunCmd([]string{"bash", "-c", "echo $UNSET"}, testEnv)
		require.Empty(t, os.Getenv("UNSET"))
	})
}
