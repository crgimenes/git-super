package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func folderExists(path string) (fs bool) {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}

		log.Println(err)

		return
	}

	fs = true

	return
}

func execHelper(path, name string, arg ...string) (err error) {
	cmd := exec.Command(name, arg...)
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()

	return
}

func gitExec(path string) (err error) {
	err = execHelper(path, "git", "stash")
	if err != nil {
		return
	}

	err = execHelper(path, "git", "checkout", "master")
	if err != nil {
		return
	}

	err = execHelper(path, "git", "pull")
	if err != nil {
		return
	}

	return
}

func visit(path string, file os.FileInfo, perr error) error {
	if perr != nil {
		return perr
	}

	if !file.IsDir() {
		return nil
	}

	if file.Name() == "vendor" {
		return filepath.SkipDir
	}

	fs := folderExists(path + "/.git")
	if !fs {
		return nil
	}

	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	log.Println(path)

	err = gitExec(path)
	if err != nil {
		return err
	}

	return filepath.SkipDir
}

func run() (err error) {
	err = filepath.Walk("./", visit)
	if errors.Is(err, io.EOF) {
		err = nil
	}

	return
}

func main() {
	err := run()
	if err != nil {
		fmt.Println(err)
	}
}
