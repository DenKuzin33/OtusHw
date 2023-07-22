package main

import (
	"bufio"
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
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()

		val := strings.TrimRight(scanner.Text(), " ")
		runes := []rune(val)
		for i, rune := range val {
			if rune == 0 {
				val = string(runes[:i])
				break
			}
		}

		env[fileInfo.Name()] = EnvValue{Value: val}
	}
	return env, nil
}
