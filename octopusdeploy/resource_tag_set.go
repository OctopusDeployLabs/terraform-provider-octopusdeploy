package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTagSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagSetCreate,
		DeleteContext: resourceTagSetDelete,
		ReadContext:   resourceTagSetRead,
		UpdateContext: resourceTagSetUpdate,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			constTag: getTagSchema(),
		},
	}
}

func getTagSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				constCanonicalTagName: {
					Optional: true,
					Type:     schema.TypeString,
				},
				constColor: {
					Required: true,
					Type:     schema.TypeString,
				},
				"description": {
					Optional: true,
					Type:     schema.TypeString,
				},
				"name": {
					Required: true,
					Type:     schema.TypeString,
				},
				constSortOrder: {
					Optional: true,
					Type:     schema.TypeInt,
				},
			},
		},
	}
}

func resourceTagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set(constName, tagSet.Name)
	d.SetId(tagSet.GetID())

	return nil
}

func buildTagSetResource(d *schema.ResourceData) *octopusdeploy.TagSet {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	var tagSet = octopusdeploy.NewTagSet(name)

	if attr, ok := d.GetOk(constTag); ok {
		tfTags := attr.([]interface{})

		for _, tfTag := range tfTags {
			tag := buildTagResource(tfTag.(map[string]interface{}))
			tagSet.Tags = append(tagSet.Tags, tag)
		}
	}

	return tagSet
}

func buildTagResource(tfTag map[string]interface{}) octopusdeploy.Tag {
	tag := octopusdeploy.Tag{
		CanonicalTagName: tfTag[constCanonicalTagName].(string),
		Color:            tfTag[constColor].(string),
		Description:      tfTag["description"].(string),
		Name:             tfTag[constName].(string),
		SortOrder:        tfTag[constSortOrder].(int),
	}

	return tag
}

func resourceTagSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := buildTagSetResource(d)

	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.Add(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagSet.GetID())

	return nil
}

func resourceTagSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := buildTagSetResource(d)
	tagSet.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.TagSets.Update(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceTagSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.TagSets.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
