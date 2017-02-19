package main

import (
	"fmt"
	"os"
	"strings"
)

func detectGopath(proj string) (string, error) {
	gp := os.Getenv("GOPATH")
	paths := strings.Split(gp, ":")
	for _, path := range paths {
		t := fmt.Sprintf("%s/src/%s", path, proj)
		if f, err := os.Stat(t); err == nil {
			if f.IsDir() {
				return path, nil
			}
		}
	}
	return "", fmt.Errorf("project %s not found, GOPATH: %s", proj, gp)
}
