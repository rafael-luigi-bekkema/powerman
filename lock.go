package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

// ErrLockFileExists indicates another power instance is already active
var ErrLockFileExists = errors.New("powerman is already running")

func findExe(pid int) (string, error) {
	path, err := os.Readlink(fmt.Sprintf("/proc/%d/exe", pid))
	if err != nil {
		return "", fmt.Errorf("could not read link for pid %d: %w", pid, err)
	}
	return path, nil
}

func lockFile() error {
	lockf, err := os.Open(lockFilePath)
	notExists := os.IsNotExist(err)
	if err != nil && !notExists {
		return fmt.Errorf("failed to stat: %w", err)
	}

	if !notExists {
		defer lockf.Close()
		var pid int
		if _, err := fmt.Fscanf(lockf, "%d", &pid); err != nil {
			return fmt.Errorf("could not read lockfile: %w", err)
		}

		p1, err := findExe(pid)
		if err == nil {
			p2, _ := findExe(os.Getpid())
			if p1 == p2 {
				return ErrLockFileExists
			}
		}
		if err := os.Remove(lockFilePath); err != nil {
			return fmt.Errorf("failed to remove lockfile: %w", err)
		}
	}

	if err := ioutil.WriteFile(lockFilePath, []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		return fmt.Errorf("error writing lockfile: %w", err)
	}
	return nil
}
