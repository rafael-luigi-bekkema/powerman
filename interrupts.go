package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type interrupts struct {
	count  int
	device *regexp.Regexp
}

func newInterrupts(devicePattern string) *interrupts {
	return &interrupts{
		device: regexp.MustCompile(devicePattern),
	}
}

func (i *interrupts) updateInterruptCount() (int, error) {
	interrupts, err := os.Open("/proc/interrupts")
	if err != nil {
		return 0, fmt.Errorf("failed to open interrupts: %w", err)
	}
	defer interrupts.Close()

	rdr := bufio.NewReader(interrupts)
	var total int
	for {
		line, err := rdr.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("failed to read line: %s", err)
		}
		values := strings.Fields(line)
		if len(values) >= 8 && i.device.MatchString(line) {
			for i := 1; i <= 4; i++ {
				num, _ := strconv.Atoi(values[i])
				total += num
			}
		}
	}
	return total, nil
}

func (i *interrupts) Name() string {
	return "/proc/interrupts"
}

func (i *interrupts) Inhibit() (bool, error) {
	count, err := i.updateInterruptCount()
	if err != nil {
		return false, err
	}
	if count == 0 {
		return false, fmt.Errorf("could not find interrupt device: %q", i.device)
	}
	if count != i.count {
		i.count = count
		return true, nil
	}
	return false, nil
}
