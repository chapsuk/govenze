package main

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

type result struct {
	manager string
	size    float64
	elapsed time.Duration
	build   bool
}

var (
	config  = flag.String("config", "config.yml", "config file name")
	tmpPath = fmt.Sprintf("%s/%s", getHomeDir(), ".govenze")
)

func main() {
	flag.Parse()

	if *config == "" {
		handleError(fmt.Errorf("config flag is required"))
	}
	cfgs, err := LoadConfig(*config)
	handleError(err)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancel()
	}()

	exitFuncs := []func(){}
	go func() {
		<-ctx.Done()
		for _, f := range exitFuncs {
			f()
		}
		log.Print("process was stopped")
		os.Exit(1)
	}()

	managers := []manager{&dep{}, &glide{}, &godep{}, &govendor{}}
	for _, item := range cfgs {
		rgopath, err := detectGopath(item.Repo)
		handleError(err)
		log.Printf("Detected GOPATH: %s", rgopath)

		projPath := fmt.Sprintf("%s/src/%s", rgopath, item.Repo)

		log.Printf("Create tmp dir: %s", tmpPath)
		err = os.Mkdir(tmpPath, os.ModePerm)
		handleError(err)

		tmpProjPath := fmt.Sprintf("%s/proj/%s", tmpPath, item.Repo)
		err = copyDir(projPath, tmpProjPath)
		handleError(err)
		clearFunc := func() {
			log.Print("restore original")
			copyDir(tmpProjPath, projPath)
			log.Print("remove tmp")
			os.RemoveAll(tmpPath)
		}
		exitFuncs = append(exitFuncs, clearFunc)

		// remove original dir, if packages duplicates in GOPATH
		// govendor throw fatal error: stack overflow
		err = os.RemoveAll(projPath)
		handleError(err)

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
				mpath := fmt.Sprintf("%s/src/%s/", mgopath, item.Repo)

				err = copyDir(tmpProjPath, mpath)
				handleError(err)

				preparePath(mpath)

				gopath := fmt.Sprintf("%s:%s", rgopath, mgopath)

				_, errs, err := m.DoVendor(mpath, gopath)
				if err != nil {
					log.Printf("%s errors: %s error: %s", m.String(), errs.Bytes(), err)
				}

				r := &result{
					manager: m.String(),
					size:    dirSizeMB(fmt.Sprintf("%s/vendor/", mpath)),
					elapsed: time.Since(start),
				}
				if item.Build != "" {
					_, errs, err := runBuild(item.Build, mpath, gopath)
					if err != nil {
						log.Printf("Run build error, manager: %s errors: %s error: %s", m.String(), errs.Bytes(), err)
					} else {
						r.build = true
					}
				}
				res = append(res, r)
			}(m)
		}
		wg.Wait()
		prettyPrint(item.Repo, res)
		clearFunc()
		exitFuncs = []func(){}
	}
}

func dirSizeMB(path string) float64 {
	var dirSize int64
	readSize := func(path string, file os.FileInfo, err error) error {
		if file == nil {
			return errors.New("empty result vendor path")
		}
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

func runBuild(buildCMD, dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error) {
	if buildCMD == "" {
		err = errors.New("wrong build command")
		return
	}
	args := strings.Split(buildCMD, " ")
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errs
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	err = cmd.Run()
	return
}
