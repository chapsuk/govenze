package main

import (
	"os"
	"runtime"
)

// Source: https://github.com/golang/go/blob/go1.8rc2/src/go/build/build.go#L260-L277

func getHomeDir() string {
	env := "HOME"
	if runtime.GOOS == "windows" {
		env = "USERPROFILE"
	} else if runtime.GOOS == "plan9" {
		env = "home"
	}
	return os.Getenv(env)
}
