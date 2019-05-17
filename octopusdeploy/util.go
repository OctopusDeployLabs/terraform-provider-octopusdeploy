package octopusdeploy

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
)

// Octopus can get itself into some race conditions, so this mutex can be used to ensure
// that we wait for other commands to finish first.
var octoMutex = mutexkv.NewMutexKV()

// Validate a value against a set of possible values
func validateValueFunc(values []string) schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		value := v.(string)
		valid := false
		for _, val := range values {
			if value == val {
				valid = true
				break
			}
		}

		if !valid {
			errors = append(errors, fmt.Errorf("%#v is an invalid value for argument %s. Must be one of %#v", value, k, values))
		}
		return
	}
}

// validateStringInSlice checks if a string is in the given slice
func validateStringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}

	return false
}

func validateAllSliceItemsInSlice(givenSlice, validationSlice []string) (string, bool) {
	for _, v := range givenSlice {
		if !validateStringInSlice(v, validationSlice) {
			return v, false
		}
	}

	return "", true
}

func getSliceFromTerraformTypeList(inputTypeList interface{}) []string {
	var newSlice []string

	terraformList := inputTypeList.([]interface{})

	for _, item := range terraformList {
		newSlice = append(newSlice, item.(string))
	}

	return newSlice
}

func getStringOrEmpty(tfAttr interface{}) string {
	if tfAttr == nil {
		return ""
	}
	return tfAttr.(string)

}

func getTenantedDeploymentSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Default:  "Untenanted",
		ValidateFunc: validateValueFunc([]string{
			"Untenanted",
			"TenantedOrUntenanted",
			"Tenanted",
		}),
	}
}
