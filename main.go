package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"time"
)

const defaultSuspendAfter = 15
const defaultDevice = "usb"
const defaultInterval = 1

var errorLog = log.New(os.Stderr, "[error] ", log.LstdFlags)
var infoLog = log.New(os.Stderr, "[info] ", log.LstdFlags)

func suspend() {
	log.Println("suspending")
	cmd := exec.Command("systemctl", "suspend")
	if err := cmd.Run(); err != nil {
		log.Printf("Failed to suspend: %s", err)
	}
}

type inhibitor interface {
	Name() string
	Inhibit() (bool, error)
}

func main() {
	after := flag.Int("after", defaultSuspendAfter, "suspend after this amount of inactivity (in minutes)")
	device := flag.String("device", defaultDevice, "device pattern to match against in /proc/interrupts")
	interval := flag.Int("interval", defaultInterval, "check activity every <interval> minutes")
	flag.Parse()

	if err := lockFile(getLockPath()); err != nil {
		errorLog.Fatal(err)
	}

	tinterval := time.Duration(*interval) * time.Minute
	suspendAfter := time.Duration(*after) * time.Minute

	inhibitors := []inhibitor{&mpris{}, newInterrupts(*device)}

	lastUpdate := time.Now()
	var inhibitedBy string
	for {
		now := time.Now()
		for _, inh := range inhibitors {
			inhibit, err := inh.Inhibit()
			if err != nil {
				errorLog.Println(err)
				continue
			}
			if inhibit {
				if by := inh.Name(); by != inhibitedBy {
					inhibitedBy = by
					infoLog.Printf("inhibited by: %s", by)
				}
				lastUpdate = now
				break
			}
		}

		if !lastUpdate.Equal(now) {
			inhibitedBy = ""
			idleTime := now.Sub(lastUpdate)
			// infoLog.Printf("idle time: %s of %s\n", idleTime, suspendAfter)
			if idleTime >= suspendAfter {
				suspend()
			}
		}

		time.Sleep(tinterval)
	}
}
