package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

const defaultSuspendAfter = 15
const defaultDevice = "usb"
const defaultInterval = 1

var logger = log.New(os.Stderr, "", 0)

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
	cfg, err := ReadConfig()
	if err != nil {
		log.Fatalf("Could not read config: %s", err)
	}

	f, err := lockFile(getLockPath())
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	tinterval := time.Duration(cfg.General.CheckInterval) * time.Minute
	suspendAfter := time.Duration(cfg.General.SuspendAfter) * time.Minute

	inhibitors := []inhibitor{&mpris{}, newInterrupts(cfg.General.DevicePattern)}

	lastUpdate := time.Now()
	var inhibitedBy string
	for {
		now := time.Now()
		for _, inh := range inhibitors {
			inhibit, err := inh.Inhibit()
			if err != nil {
				logger.Println(err)
				continue
			}
			if inhibit {
				if by := inh.Name(); by != inhibitedBy {
					inhibitedBy = by
					logger.Printf("inhibited by: %s", by)
				}
				lastUpdate = now
				break
			}
		}

		if !lastUpdate.Equal(now) {
			inhibitedBy = ""
			idleTime := now.Sub(lastUpdate)
			if idleTime >= suspendAfter {
				suspend()
			}
		}

		time.Sleep(tinterval)
	}
}
