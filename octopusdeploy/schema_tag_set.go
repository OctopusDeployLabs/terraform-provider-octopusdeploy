package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTagSet(d *schema.ResourceData) *octopusdeploy.TagSet {
	name := d.Get("name").(string)

	var tagSet = octopusdeploy.NewTagSet(name)

	if v, ok := d.GetOk("tag"); ok {
		tags := v.([]interface{})
		for _, t := range tags {
			tag := expandTag(t.(map[string]interface{}))
			tagSet.Tags = append(tagSet.Tags, tag)
		}
	}

	return tagSet
}

func expandTag(tfTag map[string]interface{}) octopusdeploy.Tag {
	tag := octopusdeploy.Tag{
		CanonicalTagName: tfTag["canonical_tag_name"].(string),
		Color:            tfTag["color"].(string),
		Description:      tfTag["description"].(string),
		Name:             tfTag["name"].(string),
		SortOrder:        tfTag["sort_order"].(int),
	}

	return tag
}

func getTagSchema() *schema.Schema {
	return &schema.Schema{
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
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
			},
		},
		Optional: true,
		Type:     schema.TypeList,
	}
}

func getTagSetDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"tag_sets": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getTagSetSchema()},
			Type:     schema.TypeList,
		},
	}
}

func getTagSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required: true,
			Type:     schema.TypeString,
		},
		"tag": getTagSchema(),
	}
}
