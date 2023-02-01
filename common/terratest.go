package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
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

func AssertTrue(t *testing.T, exists bool) {
	assert.True(t, exists)
}

func AssertEqual(t *testing.T, options any, values string) {
	for key, value := range options.(map[string]interface{}) {
		jsonValue := gjson.Get(values, key).Value()
		assert.Equal(t, value, jsonValue)
	}
}

func GetValues(json string, jsonQuery string) string {
	jsonValues := gjson.Get(json, jsonQuery)

	var getValues string

	for _, array1 := range jsonValues.Array() {
		for _, array2 := range array1.Array() {
			getValues = array2.String()
		}
	}

	return getValues
}

func GetValue(json string, jsonQuery string) string {
	getValue := gjson.Get(json, jsonQuery).String()

	return getValue
}
