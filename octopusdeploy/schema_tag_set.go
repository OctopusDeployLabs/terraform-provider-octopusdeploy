package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandTagSet(d *schema.ResourceData) *octopusdeploy.TagSet {
	name := d.Get("name").(string)

	var tagSet = octopusdeploy.NewTagSet(name)

	if v, ok := d.GetOk("tags"); ok {
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
		"id":   tagSet.GetID(),
		"name": tagSet.Name,
		"tags": flattenTags(tagSet.Tags),
	}
}

func getTagSetDataSchema() map[string]*schema.Schema {
	dataSchema := getTagSetSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"tag_set": {
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
		"id":   getIDSchema(),
		"name": getNameSchema(true),
		"tags": {
			Elem:     &schema.Resource{Schema: getTagsSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
	}
}
