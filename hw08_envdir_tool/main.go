package main

import (
	"fmt"
	"os"
)

func main() {
	// 0 путь к exe
	// 1 путь к Env
	// 2.. команда и (возможно) аргументы
	args := os.Args
	if len(args) < 3 {
		return
	}

	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	RunCmd(args[2:], env)
}
