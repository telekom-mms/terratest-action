package terratest

import (
	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"os"
	"strings"
	"terratest-action/common"
	"testing"
)

func TestAzure(t *testing.T) {
	t.Parallel()

	common.LogColor("yellow", "Azure")

	// get common test settings
	testSettings := common.GetTestSettings()
	path := testSettings["path"]

	// get terratest settings
	terratestSettings := path + "/terratest.yaml"
	settings := common.GetTerratestSettings(terratestSettings)

	// prepare Terratest Seetings
	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: path,
	}

	// unit tests
	common.LogColor("yellow", "Unit Test")

	// website::tag::2:: Run `terraform init` and `terraform plan`. Fail the test if there are any errors.
	terraform.InitAndPlan(t, terraformOptions)

	// integration tests
	common.LogColor("yellow", "Integration Test")

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

			jsonQuery := "values.root_module.child_modules.#.resources.#.values"
			tfValues := common.GetValues(tfShow, jsonQuery)

			resourceGroupName := common.GetValue(tfValues, "resource_group_name")
			resourceName := common.GetValue(tfValues, "name")
			subscriptionID := os.Getenv("ARM_SUBSCRIPTION_ID")

			// website::tag::4:: Assert
			switch function {
			case "ContainerRegistryExists":
				exists := azure.ContainerRegistryExists(t, resourceName, resourceGroupName, subscriptionID)
				common.AssertTrue(t, exists)

			case "GetContainerRegistry":
				common.AssertEqual(t, options, tfValues)

			default:
				common.LogMiss(function)
			}
		}
	}
}
