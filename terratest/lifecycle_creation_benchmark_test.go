package terratest

//
//import (
//	"fmt"
//	"testing"
//
//	"github.com/gruntwork-io/terratest/modules/terraform"
//)
//
//func BenchmarkLifecycleCreation(b *testing.B) {
//	terraformTest := &terraform.Options{
//		TerraformDir: "../examples/Lifecycle-Creation",
//	}
//
//	defer terraform.Destroy(b, terraformTest)
//
//	if _, err := terraform.InitE(b, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//
//	if _, err := terraform.PlanE(b, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//
//	if _, err := terraform.ApplyE(b, terraformTest); err != nil {
//		fmt.Println(err)
//	}
//}
