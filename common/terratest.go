package common

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
	"gopkg.in/yaml.v3"
)

var (
	testEnvVars = make(map[string]string)
	authEnvVars = make(map[string]string)
)

type TerratestSettings struct {
	Package   string      `yaml:"package"`
	Functions []string    `yaml:"functions"`
	Options   map[any]any `yaml:"options"`
}

func GetTerratestSettings(file string) *TerratestSettings {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var settings TerratestSettings
	err = yaml.NewDecoder(f).Decode(&settings)
	if err != nil {
		panic(err)
	}

	return &settings
}

func TestSetup() map[string]string {
	// Getting enVars from environment variables
	testEnvVars["TEST_TYPE"] = os.Getenv("TEST_TYPE")

	return testEnvVars
}

func AzureAuthentication() map[string]string {
	// Getting enVars from environment variables
	authEnvVars["ARM_CLIENT_ID"] = os.Getenv("AZURE_CLIENT_ID")
	authEnvVars["ARM_CLIENT_SECRET"] = os.Getenv("AZURE_CLIENT_SECRET")
	authEnvVars["ARM_SUBSCRIPTION_ID"] = os.Getenv("AZURE_SUBSCRIPTION_ID")
	authEnvVars["ARM_TENANT_ID"] = os.Getenv("AZURE_TENANT_ID")

	return authEnvVars
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
