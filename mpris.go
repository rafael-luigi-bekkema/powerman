package main

import (
	"fmt"
	"strings"

	"github.com/godbus/dbus/v5"
)

type mpris struct {
	name string
}

func (m *mpris) isAnythingPlaying() (bool, error) {
	conn, err := dbus.SessionBus()
	if err != nil {
		return false, fmt.Errorf("dbus connection failed: %w", err)
	}

	var names []string
	if err := conn.BusObject().Call("org.freedesktop.DBus.ListNames", 0).Store(&names); err != nil {
		return false, fmt.Errorf("failed to list names: %w", err)
	}

	for _, name := range names {
		if strings.HasPrefix(name, "org.mpris.MediaPlayer2.") {
			obj := conn.Object(name, "/org/mpris/MediaPlayer2")
			ident, err := obj.GetProperty("org.mpris.MediaPlayer2.Identity")
			if err != nil {
				return false, fmt.Errorf("could not get identity: %w", err)
			}
			status, err := obj.GetProperty("org.mpris.MediaPlayer2.Player.PlaybackStatus")
			if err != nil {
				return false, fmt.Errorf("could not get status: %w", err)
			}
			if status.Value().(string) == "Playing" {
				m.name = ident.Value().(string)
				return true, nil
			}
		}
	}

	return false, nil
}

func (m *mpris) Name() string {
	if m.name == "" {
		return "Media (Unknown)"
	}
	return m.name
}

func (m *mpris) Inhibit() (bool, error) {
	return m.isAnythingPlaying()
}
