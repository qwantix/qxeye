package config

import (
	"encoding/json"
	"os"

	"github.com/qwantix/qxeye/util"
)

type Config struct {
	Cascades []string                 `json:"cascades"`
	Cameras  []CameraConfig           `json:"cameras"`
	Triggers []TriggerConfig          `json:"triggers"`
	Matchers []MatcherConfig          `json:"matchers"`
	Services map[string]ServiceConfig `json:"services"`
}

type CameraZoneConfig struct {
	Name   string   `json:"name"`
	Ignore bool     `json:"ignore"`
	Color  string   `json:"color"`
	Mask   []string `json:"mask"`
}

type CameraConfig struct {
	Name        string             `json:"name"`
	Enabled     bool               `json:"enabled"`
	Endpoint    string             `json:"endpoint"`
	Matcher     string             `json:"matcher"`
	Persistance float32            `json:"persistence"`
	Zones       []CameraZoneConfig `json:"zones"`
}

type MatcherConfig struct {
	Name   string `json:"name"`
	Type   string `json:"type"`
	Params hmap   `json:"params"`
}

type TriggerConfig struct {
	On         string   `json:"on"`
	Confidence float32  `json:"confidence"`
	Zones      []string `json:"zones"`
	Delay      int      `json:"delay"`
	Service   string   `json:"service"`
	Params     hmap     `json:"params"`
}

type ServiceConfig = hmap


func (t *TriggerConfig) HasZone(zone string) bool {
	for _, z := range t.Zones {
		if z == zone {
			return true
		}
	}
	return false
}

func Load(filename string) Config {
	var cfg Config
	// Init defaults
	util.Log("config: Load config file ", filename)
	file, err := os.Open(filename)
	util.CheckErrPanic(err, "Unable to load config")
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&cfg)
	util.CheckErrPanic(err, "Unable to parse config file")
	return cfg
}
