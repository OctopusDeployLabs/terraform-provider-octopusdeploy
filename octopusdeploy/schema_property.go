package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func expandProperties(properties interface{}) map[string]string {
	if properties == nil {
		return nil
	}

	expandedProperties := make(map[string]string)
	for k, v := range properties.(map[string]interface{}) {
		expandedProperties[k] = v.(string)
	}
	return expandedProperties
}

func getPropertySchema() *schema.Schema {
	return &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"key": {
					Description: "The name of the action",
					Required:    true,
					Type:        schema.TypeString,
				},
				"value": {
					Description: "The type of action",
					Required:    true,
					Type:        schema.TypeString,
				},
			},
		},
		Optional: true,
		Type:     schema.TypeSet,
	}
}

func buildPropertiesMap(tfProperties interface{}) map[string]string {
	properties := map[string]string{}
	if tfProperties != nil {
		for _, tfProp := range tfProperties.(*schema.Set).List() {
			m := tfProp.(map[string]interface{})
			properties[m["key"].(string)] = m["value"].(string)
		}
	}
	return properties
}
