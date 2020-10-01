package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTagSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagSetCreate,
		ReadContext:   resourceTagSetRead,
		UpdateContext: resourceTagSetUpdate,
		DeleteContext: resourceTagSetDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
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
				constName: {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				constColor: {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
			},
		},
	}
}

func resourceTagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.TagSets.GetByID(id)
	if err != nil {
		// return createResourceOperationError(errorReadingTagSet, id, err)
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constTagSet, m)

	d.Set(constName, resource.Name)

	return nil
}

func buildTagSetResource(d *schema.ResourceData) *model.TagSet {
	tagSetName := d.Get(constName).(string)

	var tagSet = model.NewTagSet(tagSetName)

	if attr, ok := d.GetOk(constTag); ok {
		tfTags := attr.([]interface{})

		for _, tfTag := range tfTags {
			tag := buildTagResource(tfTag.(map[string]interface{}))
			tagSet.Tags = append(tagSet.Tags, tag)
		}
	}

	return tagSet
}

func buildTagResource(tfTag map[string]interface{}) model.Tag {
	tag := model.Tag{
		Name:  tfTag[constName].(string),
		Color: tfTag[constColor].(string),
	}

	return tag
}

func resourceTagSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := buildTagSetResource(d)
	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.TagSets.Add(tagSet)
	if err != nil {
		// return createResourceOperationError(errorCreatingTagSet, tagSet.Name, err)
		return diag.FromErr(err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceTagSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := buildTagSetResource(d)
	tagSet.ID = d.Id() // set ID so Octopus API knows which tag set to update

	diagValidate()

	apiClient := m.(*client.Client)
	resource, err := apiClient.TagSets.Update(*tagSet)
	if err != nil {
		// return createResourceOperationError(errorUpdatingTagSet, d.Id(), err)
		return diag.FromErr(err)
	}

	d.SetId(resource.ID)

	return nil
}

func resourceTagSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()
	diagValidate()

	apiClient := m.(*client.Client)
	err := apiClient.TagSets.DeleteByID(id)
	if err != nil {
		// return createResourceOperationError(errorDeletingTagSet, id, err)
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
