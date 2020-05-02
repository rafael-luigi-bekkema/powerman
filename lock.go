package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"syscall"
)

var errAlreadyRunning = errors.New("powerman is already running")

func getLockPath() string {
	cacheDir := path.Join(os.TempDir(), "go_powerman.lock")
	return cacheDir
}

func lockFile(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return nil, err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return nil, errAlreadyRunning
	}

	if _, err := file.WriteString(fmt.Sprint(os.Getpid())); err != nil {
		file.Close()
		return nil, err
	}

	return file, nil
}
