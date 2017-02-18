package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
)

const (
	dep = iota
	gb
	glide
	godep
)

type result struct {
	items []*resultItem
}

type resultItem struct {
	manager uint
	size    int64
}

var (
	manager    = []uint{dep, gb, glide, godep}
	targetPath = flag.String("path", "", "target project path")
)

func main() {
	flag.Parse()

	res := result{
		items: make([]*resultItem, len(manager)),
	}

	for _, m := range manager {
		go func(m uint) {
			gopath, err := createTmpGopath(*targetPath, m)
			handleError(err)
			defer cleare(gopath)

			newProjectPath, err := copyProject(*targetPath, gopath)
			handleError(err)

			err = goGet(gopath, newProjectPath)
			handleError(err)

			err = initVendor(gopath, newProjectPath)
			handleError(err)

			size, err := determineVendorSize(newProjectPath)
			handleError(err)

			res.items = append(res.items, &resultItem{
				manager: m,
				size:    size,
			})
		}(m)
	}

	prettyPrint(res)
}

func prettyPrint(res result) {
}

func cleare(gopath string) error {
	return nil
}

func determineVendorSize(projectPath string) (int64, error) {
	return 0, nil
}

func initVendor(gopath, projectPath string) error {
	return nil
}

func goGet(gopath, projectPath string) error {
	return nil
}

func copyProject(from, toGopath string) (string, error) {
	return "", nil
}

func createTmpGopath(basePath string, vendorManager uint) (string, error) {
	switch vendorManager {
	case dep:
	case gb:
	case glide:
	case godep:
	default:
		return "", errors.New("undefined vendor manager")
	}
	path := fmt.Sprintf("%s/%s/", basePath, managerString(vendorManager))
	err := os.Mkdir(path, os.ModePerm)
	return path, err
}

func managerString(managerType uint) string {
	switch managerType {
	case dep:
		return "dep"
	case gb:
		return "gb"
	case glide:
		return "glide"
	case godep:
		return "godep"
	default:
		return "undefined"
	}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
