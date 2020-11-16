package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTag(tag map[string]interface{}) octopusdeploy.Tag {
	return octopusdeploy.Tag{
		CanonicalTagName: tag["canonical_tag_name"].(string),
		Color:            tag["color"].(string),
		Description:      tag["description"].(string),
		Name:             tag["name"].(string),
		SortOrder:        tag["sort_order"].(int),
	}
}

func getTagSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"canonical_tag_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"color": {
			Required: true,
			Type:     schema.TypeString,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"sort_order": {
			Optional: true,
			Type:     schema.TypeInt,
		},
	}
}
