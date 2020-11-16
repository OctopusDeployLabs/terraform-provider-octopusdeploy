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

func getTagSetDataSchema() map[string]*schema.Schema {
	tagSetSchema := getTagSetSchema()
	for _, field := range tagSetSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

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
		"tag_sets": {
			Computed: true,
			Elem:     &schema.Resource{Schema: tagSetSchema},
			Type:     schema.TypeList,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
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
		"tag": {
			Elem:     &schema.Resource{Schema: getTagSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
