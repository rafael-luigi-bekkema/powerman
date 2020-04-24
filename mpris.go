package main

import (
	"bufio"
	"fmt"
	"os/exec"
	"sync"
	"time"
)

type mpris struct {
	sync.Mutex
	value bool
	err   error
}

func newMpris() *mpris {
	m := mpris{}
	go func() {
		for {
			if err := m.follow(); err != nil {
				m.Lock()
				m.err = err
				m.Unlock()
			}
			time.Sleep(time.Second * 5)
		}
	}()
	return &m
}

func (m *mpris) follow() error {
	cmd := exec.Command("playerctl", "-F", "status")
	out, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to get stdout: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to get playerctl status: %w", err)
	}
	defer cmd.Process.Kill()
	rdr := bufio.NewReader(out)
	for {
		line, err := rdr.ReadString('\n')
		if err != nil {
			return fmt.Errorf("error reading playerctl status: %w", err)
		}
		m.Lock()
		m.value = line == "Playing\n"
		m.Unlock()
	}
}

func (m *mpris) Inhibit() (bool, error) {
	m.Lock()
	defer m.Unlock()
	if m.err != nil {
		err := m.err
		m.err = nil
		return false, err
	}
	return m.value, nil
}
