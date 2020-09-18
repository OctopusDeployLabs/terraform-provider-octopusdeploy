package octopusdeploy

import (
	"fmt"

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
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tag": getTagSchema(),
		},
	}
}

func getTagSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"name": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
				"color": {
					Type:        schema.TypeString,
					Description: "The name of the step",
					Required:    true,
				},
			},
		},
	}
}

func resourceTagSetRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	tagSetID := d.Id()
	tagSet, err := apiClient.TagSets.Get(tagSetID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading tagSet %s: %s", tagSetID, err.Error())
	}

	d.Set("name", tagSet.Name)

	return nil
}

func buildTagSetResource(d *schema.ResourceData) *model.TagSet {
	tagSetName := d.Get("name").(string)

	var tagSet = model.NewTagSet(tagSetName)

	if attr, ok := d.GetOk("tag"); ok {
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
		Name:  tfTag["name"].(string),
		Color: tfTag["color"].(string),
	}

	return tag
}

func resourceTagSetCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newTagSet := buildTagSetResource(d)
	tagSet, err := apiClient.TagSets.Add(newTagSet)

	if err != nil {
		return fmt.Errorf("error creating tagSet %s: %s", newTagSet.Name, err.Error())
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
		return fmt.Errorf("error updating tagSet id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedTagSet.ID)
	return nil
}

func resourceTagSetDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	tagSetID := d.Id()

	err := apiClient.TagSets.Delete(tagSetID)

	if err != nil {
		return fmt.Errorf("error deleting tagSet id %s: %s", tagSetID, err.Error())
	}

	d.SetId("")
	return nil
}
