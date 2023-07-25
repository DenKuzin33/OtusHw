package main

import (
	"bufio"
	"bytes"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := Environment{}
	dirItems, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}

	for _, item := range dirItems {
		fileInfo, err := item.Info()
		if err != nil {
			return env, err
		}

		if fileInfo.Size() == 0 {
			env[fileInfo.Name()] = EnvValue{NeedRemove: true}
			continue
		}

		file, err := os.Open(path.Join(dir, fileInfo.Name()))
		if err != nil {
			return env, err
		}

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		file.Close()

		val := strings.TrimRight(scanner.Text(), " ")
		bytesVal := []byte(val)
		val = string(bytes.ReplaceAll(bytesVal, []byte{0}, []byte{10}))

		env[fileInfo.Name()] = EnvValue{Value: val}
	}
	return env, nil
}
