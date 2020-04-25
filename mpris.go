package main

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type mpris struct{}

func (m *mpris) isAnythingPlaying() (bool, error) {
	cmd := exec.Command("playerctl", "-a", "status")
	var out, serr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &serr
	if err := cmd.Run(); err != nil {
		if serr.String() == "No players found\n" {
			return false, nil
		}
		return false, fmt.Errorf("could not run playerctl: %w", err)
	}
	for _, status := range strings.Split(strings.TrimSpace(out.String()), "\n") {
		if status == "Playing" {
			return true, nil
		}
	}
	return false, nil
}

func (m *mpris) Inhibit() (bool, error) {
	return m.isAnythingPlaying()
}
