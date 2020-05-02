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

func lockFile(filePath string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}

	err = syscall.Flock(int(file.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
	if err != nil {
		file.Close()
		return errAlreadyRunning
	}

	if _, err := file.WriteString(fmt.Sprint(os.Getpid())); err != nil {
		file.Close()
		return err
	}

	return nil
}
