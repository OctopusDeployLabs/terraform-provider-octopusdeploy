package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTagSet(d *schema.ResourceData) *octopusdeploy.TagSet {
	name := d.Get("name").(string)

	var tagSet = octopusdeploy.NewTagSet(name)
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

	if v, ok := d.GetOk("tag"); ok {
		tags := v.([]interface{})
		for _, t := range tags {
			tag := expandTag(t.(map[string]interface{}))
			tagSet.Tags = append(tagSet.Tags, tag)
		}
	}

	return tagSet
}

func flattenTagSet(tagSet *octopusdeploy.TagSet) map[string]interface{} {
	if tagSet == nil {
		return nil
	}

	return map[string]interface{}{
		"description": tagSet.Description,
		"id":          tagSet.GetID(),
		"name":        tagSet.Name,
		"sort_order":  tagSet.SortOrder,
		"space_id":    tagSet.SpaceID,
		"tag":         flattenTags(tagSet.Tags),
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
		"tag": {
			Description: "A list of tags.",
			Elem:        &schema.Resource{Schema: getTagsSchema()},
			Optional:    true,
			Type:        schema.TypeList,
		},
	}
}

func setTagSet(ctx context.Context, d *schema.ResourceData, tagSet *octopusdeploy.TagSet) error {
	d.Set("description", tagSet.Description)
	d.Set("id", tagSet.GetID())
	d.Set("name", tagSet.Name)
	d.Set("sort_order", tagSet.SortOrder)
	d.Set("space_id", tagSet.SpaceID)

	if err := d.Set("tag", flattenTags(tagSet.Tags)); err != nil {
		return fmt.Errorf("error setting tag: %s", err)
	}

	return nil
}
