package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"
)

var device = regexp.MustCompile(`usb`)
var suspendAfter = time.Minute * 15

func parseInterrupts(input io.Reader) (int, error) {
	rdr := bufio.NewReader(input)
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
		if len(values) >= 8 && device.MatchString(line) {
			for i := 1; i <= 4; i++ {
				num, _ := strconv.Atoi(values[i])
				total += num
			}
		}
	}
	return total, nil
}

func suspend() {
	cmd := exec.Command("systemctl", "suspend")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to suspend: %s", err)
	}
}

type lastCount struct {
	value int
	at    time.Time
}

type mprisBusy struct {
	sync.Mutex
	value bool
}

func (mb *mprisBusy) Set(value bool) {
	mb.Lock()
	defer mb.Unlock()
	mb.value = value
}

func (mb *mprisBusy) Get() bool {
	mb.Lock()
	defer mb.Unlock()
	return mb.value
}

var mprisB mprisBusy

func mpris() error {
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
		mprisB.Set(line == "Playing\n")
	}
}

func main() {
	log.SetFlags(0)
	go mpris()
	interrupts, err := os.Open("/proc/interrupts")
	if err != nil {
		log.Fatalf("Failed to open interrupts: %s", err)
	}
	defer interrupts.Close()

	var lc lastCount
	for {
		var reset bool
		var newCount int
		if mprisB.Get() {
			reset = true
		} else {
			count, err := parseInterrupts(interrupts)
			interrupts.Seek(0, 0)
			if err != nil {
				log.Fatalf("Failed to parse interrupts: %s", err)
			}
			reset = count != lc.value
			newCount = count
		}

		now := time.Now()
		if reset {
			lc.at = now
			if newCount > 0 {
				lc.value = newCount
			}
		}
		idleTime := now.Sub(lc.at)
		fmt.Printf("idle time: %s of %s\n", idleTime, suspendAfter)
		if idleTime >= suspendAfter {
			suspend()
		}
		time.Sleep(time.Second * 2)
	}
}
