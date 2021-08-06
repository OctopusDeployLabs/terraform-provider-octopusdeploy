package octopusdeploy

import (
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getImporter() *schema.ResourceImporter {
	return &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
}

func expandArray(values []interface{}) []string {
	if values == nil {
		return nil
	}

	s := make([]string, len(values))
	for i, v := range values {
		s[i] = v.(string)
	}
	return s
}

func flattenArray(values []string) []interface{} {
	if values == nil {
		return nil
	}

	s := make([]interface{}, len(values))
	for i, v := range values {
		s[i] = v
	}
	return s
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

func getSliceFromTerraformTypeList(list interface{}) []string {
	if list == nil {
		return nil
	}

	if _, ok := list.([]interface{}); !ok {
		return nil
	}

	var newSlice []string
	for _, v := range list.([]interface{}) {
		if v != nil {
			newSlice = append(newSlice, v.(string))
		}
	}
	return newSlice
}

func isEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func logResource(name string, resource interface{}) {
	log.Printf("[DEBUG] %s: %#v", name, resource)
}

func getStringOrEmpty(tfAttr interface{}) string {
	if tfAttr == nil {
		return ""
	}
	return tfAttr.(string)
}
