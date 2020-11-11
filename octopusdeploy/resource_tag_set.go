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
		Importer:      getImporter(),
		ReadContext:   resourceTagSetRead,
		Schema:        getTagSetSchema(),
		UpdateContext: resourceTagSetUpdate,
	}
}

func resourceTagSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := expandTagSet(d)

	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.Add(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagSet.GetID())

	return nil
}

func resourceTagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("name", tagSet.Name)
	d.SetId(tagSet.GetID())

	return nil
}

func resourceTagSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := expandTagSet(d)
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
