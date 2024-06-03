package octopusdeploy

import (
	"hash/crc32"
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

	if v, ok := list.([]string); ok {
		return v
	}

	terraformList, ok := list.([]interface{})
	if !ok {
		terraformSet, ok := list.(*schema.Set)
		if ok {
			terraformList = terraformSet.List()
		} else {
			// It's not a list or set type
			return nil
		}
	}
	var newSlice []string
	for _, v := range terraformList {
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

func stringHashCode(s string) int {
	v := int(crc32.ChecksumIEEE([]byte(s)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	return 0
}

func formatBoolForActionProperty(b bool) string {
	if b {
		return "True"
	}
	return "False"
}
