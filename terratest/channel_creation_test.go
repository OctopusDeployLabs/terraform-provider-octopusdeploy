package terratest

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func TestChannelCreation(test *testing.T) {
	terraformTest := &terraform.Options{
		TerraformDir: "../examples/Channel-Creation",
		VarFiles:     []string{"/Users/michaellevan/gitrepos/terraform-provider-octopusdeploy/examples/terraform.tfvars"},
	}

	defer terraform.Destroy(test, terraformTest)

	if _, err := terraform.InitE(test, terraformTest); err != nil {
		fmt.Println(err)
	}

	if _, err := terraform.PlanE(test, terraformTest); err != nil {
		fmt.Println(err)
	}

	if _, err := terraform.ApplyE(test, terraformTest); err != nil {
		fmt.Println(err)
	}

}
