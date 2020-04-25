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
	flag.Parse()

	if err := lockFile(getLockPath()); err != nil {
		errorLog.Fatal(err)
	}

	suspendAfter := time.Duration(*after) * time.Minute

	inhibitors := []inhibitor{&mpris{}, newInterrupts(*device)}

	lastUpdate := time.Now()
	for {
		now := time.Now()
		for _, inh := range inhibitors {
			inhibit, err := inh.Inhibit()
			if err != nil {
				errorLog.Println(err)
				continue
			}
			if inhibit {
				infoLog.Printf("inhibited by: %s", inh.Name())
				lastUpdate = now
				break
			}
		}

		idleTime := now.Sub(lastUpdate)
		infoLog.Printf("idle time: %s of %s\n", idleTime, suspendAfter)
		if idleTime >= suspendAfter {
			suspend()
		}

		time.Sleep(time.Minute)
	}
}
