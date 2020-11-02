package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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

	client := m.(*octopusdeploy.Client)
	resource, err := client.TagSets.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingTagSet, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constTagSet, m)

	d.Set(constName, resource.Name)

	return nil
}

func buildTagSetResource(d *schema.ResourceData) *octopusdeploy.TagSet {
	tagSetName := d.Get(constName).(string)

	var tagSet = octopusdeploy.NewTagSet(tagSetName)

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
		Name:  tfTag[constName].(string),
		Color: tfTag[constColor].(string),
	}

	return tag
}

func resourceTagSetCreate(d *schema.ResourceData, m interface{}) error {
	tagSet := buildTagSetResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.TagSets.Add(tagSet)
	if err != nil {
		return createResourceOperationError(errorCreatingTagSet, tagSet.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceTagSetUpdate(d *schema.ResourceData, m interface{}) error {
	tagSet := buildTagSetResource(d)
	tagSet.ID = d.Id() // set ID so Octopus API knows which tag set to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.TagSets.Update(*tagSet)
	if err != nil {
		return createResourceOperationError(errorUpdatingTagSet, d.Id(), err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceTagSetDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.TagSets.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingTagSet, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}
