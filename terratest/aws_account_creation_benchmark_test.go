package terratest

import (
	"fmt"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

func BenchmarkAWSCreation(b *testing.B) {
	b.Skip("examples tests are outdated and are likely to be removed.")
	terraformTest := &terraform.Options{
		TerraformDir: "../examples/AWS-Account",
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
