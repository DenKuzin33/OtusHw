package main

import (
	"bufio"
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	outFilePath := "testdata/outFile.txt"

	expected1, _ := os.OpenFile("testdata/out_offset0_limit0.txt", os.O_RDONLY, 0o444)
	expected2, _ := os.OpenFile("testdata/out_offset0_limit10.txt", os.O_RDONLY, 0o444)
	expected3, _ := os.OpenFile("testdata/out_offset0_limit1000.txt", os.O_RDONLY, 0o444)
	expected4, _ := os.OpenFile("testdata/out_offset0_limit10000.txt", os.O_RDONLY, 0o444)
	expected5, _ := os.OpenFile("testdata/out_offset100_limit1000.txt", os.O_RDONLY, 0o444)
	expected6, _ := os.OpenFile("testdata/out_offset6000_limit1000.txt", os.O_RDONLY, 0o444)

	testData := []struct {
		offset, limit int
		expected      *os.File
	}{
		{0, 0, expected1},
		{0, 10, expected2},
		{0, 1000, expected3},
		{0, 10000, expected4},
		{100, 1000, expected5},
		{6000, 1000, expected6},
	}
	defer os.Remove(outFilePath)

	for _, data := range testData {
		t.Run(data.expected.Name(), func(t *testing.T) {
			Copy("testdata/input.txt", outFilePath, int64(data.offset), int64(data.limit))
			require.True(t, deepCompare(data.expected.Name(), outFilePath))
		})
	}

	t.Run("offset exceeds file size", func(t *testing.T) {
		err := Copy("testdata/input.txt", outFilePath, 99999, 0)
		require.Error(t, err)
	})
}

func deepCompare(file1, file2 string) bool {
	source, err := os.Open(file1)
	if err != nil {
		return false
	}
	defer source.Close()

	target, err := os.Open(file2)
	if err != nil {
		return false
	}
	defer target.Close()

	sscan := bufio.NewScanner(source)
	tscan := bufio.NewScanner(target)

	for sscan.Scan() {
		tscan.Scan()
		if !bytes.Equal(sscan.Bytes(), tscan.Bytes()) {
			return false
		}
	}

	return true
}
