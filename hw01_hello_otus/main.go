package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	helloOtus := stringutil.Reverse("Hello, OTUS!")
	fmt.Print(helloOtus)
}
