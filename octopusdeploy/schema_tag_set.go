package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tagsets"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTagSet(d *schema.ResourceData) *tagsets.TagSet {
	name := d.Get("name").(string)

	tagSet := tagsets.NewTagSet(name)
	tagSet.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		tagSet.Description = v.(string)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		tagSet.SortOrder = int32(v.(int))
	}

	if v, ok := d.GetOk("space_id"); ok {
		tagSet.SpaceID = v.(string)
	}

	return tagSet
}

func flattenTagSet(tagSet *tagsets.TagSet) map[string]interface{} {
	if tagSet == nil {
		return nil
	}

	return map[string]interface{}{
		"description": tagSet.Description,
		"id":          tagSet.GetID(),
		"name":        tagSet.Name,
		"sort_order":  tagSet.SortOrder,
		"space_id":    tagSet.SpaceID,
	}
}

func getTagSetDataSchema() map[string]*schema.Schema {
	dataSchema := getTagSetSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"tag_sets": {
			Computed:    true,
			Description: "A list of tag sets that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"take": getQueryTake(),
	}
}

func getTagSetSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": getDescriptionSchema("tag set"),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"sort_order":  getSortOrderSchema(),
		"space_id":    getSpaceIDSchema(),
	}
}

func setTagSet(ctx context.Context, d *schema.ResourceData, tagSet *tagsets.TagSet) error {
	d.Set("description", tagSet.Description)
	d.Set("id", tagSet.GetID())
	d.Set("name", tagSet.Name)
	d.Set("sort_order", tagSet.SortOrder)
	d.Set("space_id", tagSet.SpaceID)

	return nil
}
