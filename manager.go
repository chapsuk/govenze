package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

type manager interface {
	String() string
	DoVendor(dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error)
}

type gb struct{}
type glide struct{}
type godep struct{}
type dep struct{}

func (d *dep) String() string { return "dep" }
func (d *dep) DoVendor(dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error) {
	cmd := exec.Command("dep", "init")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errs
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	err = cmd.Run()
	if err == nil {
		cmdUpdate := exec.Command("dep", "ensure", "-update")
		cmdUpdate.Dir = dir
		cmdUpdate.Stdout = &out
		cmdUpdate.Stderr = &errs
		cmdUpdate.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
		cmdUpdate.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		err = cmdUpdate.Run()
	}
	return
}

func (gb *gb) String() string { return "gb" }
func (gb *gb) DoVendor(dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error) {
	cmd := exec.Command("gb", "vendor")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errs
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	err = cmd.Run()
	return
}

func (g *glide) String() string { return "glide" }
func (g *glide) DoVendor(dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error) {
	cmd := exec.Command("glide", "create", "--non-interactive")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errs
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	err = cmd.Run()
	if err == nil {
		cmdInstall := exec.Command("glide", "install", "--skip-test")
		cmdInstall.Dir = dir
		cmdInstall.Stdout = &out
		cmdInstall.Stderr = &errs
		cmdInstall.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
		cmdInstall.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
		err = cmdInstall.Run()
	}
	return
}

func (g *godep) String() string { return "godep" }
func (g *godep) DoVendor(dir, gopath string) (out bytes.Buffer, errs bytes.Buffer, err error) {
	cmd := exec.Command("godep", "save", "./...")
	cmd.Dir = dir
	cmd.Stdout = &out
	cmd.Stderr = &errs
	cmd.Env = append(cmd.Env, fmt.Sprintf("GOPATH=%s", gopath))
	cmd.Env = append(cmd.Env, fmt.Sprintf("PATH=%s", os.Getenv("PATH")))
	err = cmd.Run()
	return
}
