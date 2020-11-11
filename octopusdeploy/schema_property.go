package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getPropertySchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeSet,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				constKey: {
					Type:        schema.TypeString,
					Description: "The name of the action",
					Required:    true,
				},
				constValue: {
					Type:        schema.TypeString,
					Description: "The type of action",
					Required:    true,
				},
			},
		},
	}
}

func buildPropertiesMap(tfProperties interface{}) map[string]string {
	properties := map[string]string{}
	if tfProperties != nil {
		for _, tfProp := range tfProperties.(*schema.Set).List() {
			m := tfProp.(map[string]interface{})
			properties[m[constKey].(string)] = m[constValue].(string)
		}
	}
	return properties
}
