package common

import (
	"os"

	"gopkg.in/yaml.v3"
)

type TerratestSettings struct {
	Package   string      `yaml:"package"`
	Functions []string    `yaml:"functions"`
	Options   map[any]any `yaml:"options"`
}

func GetTerratestSettings(file string) TerratestSettings {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var settings *TerratestSettings
	yaml.NewDecoder(f).Decode(&settings)

	return *settings
}
