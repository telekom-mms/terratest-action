package terratest

import (
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"strconv"
	"strings"
	"terratest-action/common"
	"testing"
)

func TestAzure(t *testing.T) {
	t.Parallel()

	common.LogColor("yellow", "Azure")

	// get common test settings
	testSetup := common.TestSetup()
	azureCredentials := common.AzureAuthentication()
	testSettings := common.GetTestSettings()
	path := testSettings["path"]

	// prepare Terratest Seetings
	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		// path to code that will be tested
		TerraformDir: path,
		// globalvariables for user account
		EnvVars: azureCredentials,
	}

	// unit tests
	if testSetup["TEST_TYPE"] == "unit" {
		common.LogColor("yellow", "Unit Test")

		// website::tag::2:: Run `terraform init` and `terraform plan`. Fail the test if there are any errors.
		terraform.InitAndPlan(t, terraformOptions)
	}

	// integration tests
	if testSetup["TEST_TYPE"] == "integration" {
		common.LogColor("yellow", "Integration Test")

		// get terratest settings for integration test
		terratestSettings := path + "/terratest.yaml"
		settings := common.GetTerratestSettings(terratestSettings)

		// website::tag::5:: At the end of the test, run `terraform destroy` to clean up any resources that were created
		defer terraform.Destroy(t, terraformOptions)

		// website::tag::2:: Run `terraform init` and `terraform apply`. Fail the test if there are any errors.
		terraform.InitAndApply(t, terraformOptions)

		for _, function := range settings.Functions {
			functionCase := strings.Contains(function, "Exists")
			options := settings.Options[function]

			if functionCase == false && options == nil {
				common.LogMiss("options for function " + function)
			} else {
				common.LogColor("yellow", "> "+function)

				// website::tag::3:: Run `terraform show` to get the values after build
				tfShow := terraform.Show(t, terraformOptions)
				subscriptionID := azureCredentials["ARM_SUBSCRIPTION_ID"]

				// website::tag::4:: Assert
				getAzureTestSetting := common.GetAzureTestSetting(function)
				tfValues := common.GetAllValues(tfShow, getAzureTestSetting)
				tfIndex := common.GetIndex(tfValues)

				for i := uint64(0); i < tfIndex; i++ {
					// convert int to string, needed as string for GetValue regex
					idx := strconv.FormatUint(uint64(i), 10)

					resourceValues := common.GetValues(tfValues, idx)
					resourceGroupName := common.GetValue(resourceValues, "resource_group_name")
					resourceName := common.GetValue(resourceValues, "name")

					common.LogColor("green", resourceName)

					switch function {
					case "ContainerRegistryExists":
						exists := azure.ContainerRegistryExists(t, resourceName, resourceGroupName, subscriptionID)
						common.AssertTrue(t, exists)
					case "ContainerRegistryShow":
						common.AssertEqual(t, options, resourceValues)
					case "VirtualNetworkExists":
						exists := azure.VirtualNetworkExists(t, resourceName, resourceGroupName, subscriptionID)
						common.AssertTrue(t, exists)
					case "SubnetExists":
						virtualNetworkName := common.GetValue(resourceValues, "virtual_network_name")
						exists := azure.SubnetExists(t, resourceName, virtualNetworkName, resourceGroupName, subscriptionID)
						common.AssertTrue(t, exists)
					case "PublicAddressExists":
						exists := azure.PublicAddressExists(t, resourceName, resourceGroupName, subscriptionID)
						common.AssertTrue(t, exists)
					case "PublicAddressShow":
						common.AssertEqual(t, options, resourceValues)
					case "NetworkInterfaceExists":
						exists := azure.NetworkInterfaceExists(t, resourceName, resourceGroupName, subscriptionID)
						common.AssertTrue(t, exists)
					default:
						common.LogMiss(function)
					}
				}
			}
		}
	}
}
