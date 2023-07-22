package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envSet, err := ReadDir("testdata/env")
	require.NoError(t, err)

	t.Run("unset", func(t *testing.T) {
		require.True(t, envSet["UNSET"].NeedRemove)
	})

	t.Run("terminal zero", func(t *testing.T) {
		require.Equal(t, "   foo", envSet["FOO"].Value)
	})

	t.Run("empty", func(t *testing.T) {
		require.Empty(t, envSet["EMPTY"].Value)
	})

	t.Run("only first line", func(t *testing.T) {
		require.Equal(t, "bar", envSet["BAR"].Value)
	})
}
