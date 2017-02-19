package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type result struct {
	manager string
	size    float64
	elapsed time.Duration
}

var (
	target  = flag.String("target", "", "target project")
	clearm  = flag.Bool("clear", true, "clear tmp dirafter complete")
	tmpPath = fmt.Sprintf("%s/%s", getHomeDir(), ".govenze")
)

func main() {
	flag.Parse()

	if *target == "" {
		handleError(fmt.Errorf("target flag is required"))
	}

	rgopath, err := detectGopath(*target)
	handleError(err)
	log.Printf("Detected GOPATH: %s", rgopath)

	projPath := fmt.Sprintf("%s/src/%s", rgopath, *target)

	log.Printf("Create tmp dir: %s", tmpPath)
	err = os.Mkdir(tmpPath, os.ModePerm)
	handleError(err)

	if *clearm {
		defer func() {
			err := clear(tmpPath)
			if err != nil {
				log.Printf("Clear tmp dir error: %s", err)
			}
		}()
	}

	tmpProjPath := fmt.Sprintf("%s/proj/%s", tmpPath, *target)
	err = copyDir(projPath, tmpProjPath)
	handleError(err)

	managers := []manager{&dep{}, &glide{}, &godep{}}
	res := make([]*result, len(managers))
	wg := sync.WaitGroup{}
	wg.Add(len(managers))
	for _, m := range managers {
		go func(m manager) {
			start := time.Now()
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Manager %s stopped with error: %s", m.String(), r)
				}
				wg.Done()
			}()
			mgopath := fmt.Sprintf("%s/%s", tmpPath, m.String())
			mpath := fmt.Sprintf("%s/src/%s/", mgopath, *target)

			err = copyDir(tmpProjPath, mpath)
			handleError(err)

			preparePath(mpath)

			gopath := fmt.Sprintf("%s:%s", rgopath, mgopath)

			_, errs, err := m.DoVendor(mpath, gopath)
			if err != nil {
				log.Printf("%s error: %s", m.String(), errs.Bytes())
			}

			res = append(res, &result{
				manager: m.String(),
				size:    dirSizeMB(fmt.Sprintf("%s/vendor/", mpath)),
				elapsed: time.Since(start),
			})
		}(m)
	}
	wg.Wait()

	prettyPrint(res)
}

func prettyPrint(res []*result) {
	for _, r := range res {
		if r == nil {
			continue
		}
		fmt.Printf("\nVendor manager: %s\n", r.manager)
		fmt.Print("======\n")
		fmt.Printf("Size: %.2fMb\n", r.size)
		fmt.Printf("Time: %.2fs\n\n", r.elapsed.Seconds())
	}
}

func dirSizeMB(path string) float64 {
	var dirSize int64

	readSize := func(path string, file os.FileInfo, err error) error {
		if !file.IsDir() {
			dirSize += file.Size()
		}

		return nil
	}

	filepath.Walk(path, readSize)

	sizeMB := float64(dirSize) / 1024.0 / 1024.0

	return sizeMB
}

func preparePath(path string) {
	dirs := []string{"vendor", "Godeps"}
	for _, d := range dirs {
		p := fmt.Sprintf("%s/%s", path, d)
		if dir, err := os.Stat(p); err == nil {
			if dir.IsDir() {
				err = os.RemoveAll(p)
				if err != nil {
					log.Printf("remove excluded dir %s error %s", d, err)
				}
			}
		}
	}

	files := []string{"glide.lock", "glide.yaml", "glide.yml", "lock.json", "manifest.json"}
	for _, f := range files {
		p := fmt.Sprintf("%s/%s", path, f)
		if _, err := os.Stat(p); err == nil {
			err = os.Remove(p)
			if err != nil {
				log.Printf("remove excluded file %s error %s", p, err)
			}
		}
	}
}

func clear(path string) error {
	return os.RemoveAll(path)
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}
