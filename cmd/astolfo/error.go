package main

import (
	"errors"
	"fmt"
	"os"
)

var (
	errEmptyArg         = errors.New("none of the arguments must have an empty string")
	errInsufficientArgs = errors.New("needs only 2 arguments")
	errPassLength       = fmt.Errorf("password length must be within the range of %d to %d", minGenPassLength, maxGenPassLength)
)

func die(err error, c *int, exitCode int) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err.Error())
	*c = exitCode
}

func verbose(str string) {
	if isVerbose {
		fmt.Printf(str)
	}
}

func warn(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err.Error())
}
