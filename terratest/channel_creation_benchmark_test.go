package terratest

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func BenchmarkChannelCreation(b *testing.B) {
	terraformTest := &terraform.Options{
		TerraformDir: "../examples/Channel-Creation",
		VarFiles:     []string{"C:/Users/Mike/Desktop/terraform-provider-octopusdeploy/examples/terraform.tfvars"},
	}

	defer terraform.Destroy(b, terraformTest)

	if _, err := terraform.InitE(b, terraformTest); err != nil {
		fmt.Println(err)
	}

	if _, err := terraform.PlanE(b, terraformTest); err != nil {
		fmt.Println(err)
	}

	if _, err := terraform.ApplyE(b, terraformTest); err != nil {
		fmt.Println(err)
	}
}
