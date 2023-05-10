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

func GetAzureTestSetting(testcase string) string {

	var resource string

	switch testcase {
	case "ContainerRegistryExists":
		resource = "azurerm_container_registry"
	case "ContainerRegistryShow":
		resource = "azurerm_container_registry"
	case "VirtualNetworkExists":
		resource = "azurerm_virtual_network"
	case "SubnetExists":
		resource = "azurerm_subnet"
	case "PublicAddressExists":
		resource = "azurerm_public_ip"
	case "PublicAddressShow":
		resource = "azurerm_public_ip"
	case "NetworkInterfaceExists":
		resource = "azurerm_network_interface"
	default:
		LogMiss(testcase)
	}

	return resource
}

func AssertTrue(t *testing.T, exists bool) {
	assert.True(t, exists)
}

func AssertEqual(t *testing.T, options any, values string) {
	for key, value := range options.(map[string]interface{}) {
		jsonValue := gjson.Get(gjson.Get(values, "{#."+key+"}").String(), ""+key+".0").Value()
		assert.Equal(t, value, jsonValue)
	}
}

func GetAllValues(json string, resource string) string {
	getAllValues := gjson.Get(json, "values.root_module.child_modules.#.resources.#(type=\""+resource+"\")#.values")

	return getAllValues.String()
}

func GetValues(json string, idx string) string {
	getValues := gjson.Get(json, "#."+idx+"")

	return getValues.String()
}

func GetValue(json string, value string) string {
	getValue := gjson.Get(gjson.Get(json, "{#."+value+"}").String(), ""+value+".0")

	return getValue.String()
}

func GetIndex(json string) uint64 {
	getIndex := gjson.Get(gjson.Get(json, "#.#").String(), "0")

	return getIndex.Uint()
}
