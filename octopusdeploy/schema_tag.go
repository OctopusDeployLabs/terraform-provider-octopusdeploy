package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getTagSchema() map[string]*schema.Schema {
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
		"name":        getNameSchema(true),
		"sort_order": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeInt,
		},
		"tag_set_id": {
			Description: "The ID of the associated tag set.",
			Required:    true,
			Type:        schema.TypeString,
		},
	}
}

func expandTag(d *schema.ResourceData) *tagsets.Tag {
	color := d.Get("color").(string)
	name := d.Get("name").(string)

	tag := tagsets.NewTag(name, color)
	tag.ID = d.Id()

	if v, ok := d.GetOk("canonical_tag_name"); ok {
		tag.CanonicalTagName = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		tag.Description = v.(string)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		tag.SortOrder = v.(int)
	}

	return tag
}

func setTag(ctx context.Context, d *schema.ResourceData, tag *tagsets.Tag, tagSet *tagsets.TagSet) error {
	d.Set("canonical_tag_name", tag.CanonicalTagName)
	d.Set("color", tag.Color)
	d.Set("description", tag.Description)
	d.Set("name", tag.Name)
	d.Set("sort_order", tag.SortOrder)
	d.Set("tag_set_id", tagSet.GetID())
	d.SetId(tag.ID)

	return nil
}
