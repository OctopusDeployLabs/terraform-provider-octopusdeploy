package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTagSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTagSetCreate,
		DeleteContext: resourceTagSetDelete,
		Description:   "This resource manages tag sets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceTagSetRead,
		Schema:        getTagSetSchema(),
		UpdateContext: resourceTagSetUpdate,
	}
}

func resourceTagSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := expandTagSet(d)

	log.Printf("[INFO] creating tag set: %#v", tagSet)

	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.Add(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(tagSet.GetID())

	log.Printf("[INFO] tag set created (%s)", d.Id())
	return nil
}

func resourceTagSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting tag set (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	if err := client.TagSets.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] tag set deleted")
	return nil
}

func resourceTagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading tag set (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	tagSet, err := client.TagSets.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] tag set (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	if err := setTagSet(ctx, d, tagSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tag set read: %#v", tagSet)
	return nil
}

func resourceTagSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := expandTagSet(d)

	log.Printf("[INFO] updating tag set: %#v", tagSet)

	client := m.(*octopusdeploy.Client)
	updatedTagSet, err := client.TagSets.Update(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTagSet(ctx, d, updatedTagSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tag set updated (%s)", d.Id())
	return nil
}
