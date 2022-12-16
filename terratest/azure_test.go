package terratest

import (
	"common"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestAzure(t *testing.T) {
	t.Parallel()

	common.LogColor("yellow", "Azure")

	// get common test settings
	testSettings := common.GetTestSettings()
	path := testSettings["path"]

	// prepare Terratest Seetings
	// website::tag::1:: Configure Terraform setting up a path to Terraform code.
	terraformOptions := &terraform.Options{
		TerraformDir: path,
	}

	// unit tests
	common.LogColor("yellow", "Unit Tests")

	// website::tag::2:: Run `terraform init` and `terraform plan`. Fail the test if there are any errors.
	terraform.InitAndPlan(t, terraformOptions)
}
