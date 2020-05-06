package main

import (
	"fmt"
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type GeneralConfig struct {
	DevicePattern string `toml:"device_pattern"`
	SuspendAfter  int    `toml:"suspend_after"`
	CheckInterval int    `toml:"check_interval"`
}

type Config struct {
	General GeneralConfig `toml:"general"`
}

func ReadConfig() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("could not find user config dir: %w", err)
	}
	configDir = path.Join(configDir, "powerman")
	os.MkdirAll(configDir, 0755)
	configPath := path.Join(configDir, "config")

	var cfg = Config{
		General: GeneralConfig{
			DevicePattern: "usb",
			SuspendAfter:  20,
			CheckInterval: 1,
		},
	}

	if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
		if os.IsNotExist(err) {
			return &cfg, nil
		}
		return nil, fmt.Errorf("could not read config: %w", err)
	}

	return &cfg, nil
}
