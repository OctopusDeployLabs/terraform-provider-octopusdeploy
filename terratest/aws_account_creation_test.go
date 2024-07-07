//package terratest
//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/gruntwork-io/terratest/modules/terraform"
//)
//
//func TestAWSCreation(t *testing.T) {
//	terraformTest := &terraform.Options{
//		TerraformDir: "../examples/AWS-Account",
//	}
//
//	defer terraform.Destroy(t, terraformTest)
//
//	if _, err := terraform.InitE(t, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//
//	if _, err := terraform.PlanE(t, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//
//	if _, err := terraform.ApplyE(t, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//}
