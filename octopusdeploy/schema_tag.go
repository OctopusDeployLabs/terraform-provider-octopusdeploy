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
		ID:               tag["id"].(string),
		Name:             tag["name"].(string),
		SortOrder:        tag["sort_order"].(int),
	}
}

func flattenTags(tags []octopusdeploy.Tag) []map[string]interface{} {
	var flattenedTags = make([]map[string]interface{}, len(tags))
	for key, tag := range tags {
		flattenedTags[key] = map[string]interface{}{
			"canonical_tag_name": tag.CanonicalTagName,
			"color":              tag.Color,
			"description":        tag.Description,
			"id":                 tag.ID,
			"name":               tag.Name,
			"sort_order":         tag.SortOrder,
		}
	}

	return flattenedTags
}

func getTagsSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"canonical_tag_name": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"color": {
			Required: true,
			Type:     schema.TypeString,
		},
		"description": getDescriptionSchema("tag"),
		"id": {
			Computed:    true,
			Description: "The unique ID for the tag.",
			Type:        schema.TypeString,
		},
		"name": getNameSchema(true),
		"sort_order": {
			Computed: true,
			Type:     schema.TypeInt,
		},
	}
}
