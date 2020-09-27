package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceTagSet() *schema.Resource {
	return &schema.Resource{
		Create: resourceTagSetCreate,
		Read:   resourceTagSetRead,
		Update: resourceTagSetUpdate,
		Delete: resourceTagSetDelete,

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

func resourceTagSetRead(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.TagSets.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingTagSet, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

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

func resourceTagSetCreate(d *schema.ResourceData, m interface{}) error {
	newTagSet := buildTagSetResource(d)

	apiClient := m.(*client.Client)
	tagSet, err := apiClient.TagSets.Add(newTagSet)
	if err != nil {
		return createResourceOperationError(errorCreatingTagSet, newTagSet.Name, err)
	}

	d.SetId(tagSet.ID)
	return nil
}

func resourceTagSetUpdate(d *schema.ResourceData, m interface{}) error {
	tagSet := buildTagSetResource(d)
	tagSet.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)
	updatedTagSet, err := apiClient.TagSets.Update(tagSet)
	if err != nil {
		return createResourceOperationError(errorUpdatingTagSet, d.Id(), err)
	}

	d.SetId(updatedTagSet.ID)
	return nil
}

func resourceTagSetDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.TagSets.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingTagSet, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
