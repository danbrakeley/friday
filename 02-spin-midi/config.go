package main

import (
	"encoding/json"
	"os"
)

type Config struct {
	MidiDevice      string `json:"midi_device"`
	Knob1Chan       int    `json:"knob1_chan"`
	Knob1Controller int    `json:"knob1_controller"`
}

func LoadConfig(path string) (Config, error) {
	var config Config

	fp, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer fp.Close()
	err = json.NewDecoder(fp).Decode(&config)
	return config, err
}

func SaveConfig(path string, config Config) error {
	fp, err := os.Create(path)
	if err != nil {
		return err
	}
	defer fp.Close()
	return json.NewEncoder(fp).Encode(config)
}
