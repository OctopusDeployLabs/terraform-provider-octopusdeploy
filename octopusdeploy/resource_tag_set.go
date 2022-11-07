package octopusdeploy

import (
	"context"
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	octopus := m.(*client.Client)
	createdTagSet, err := octopus.TagSets.Add(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTagSet(ctx, d, createdTagSet); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdTagSet.GetID())

	log.Printf("[INFO] tag set created (%s)", d.Id())
	return nil
}

func resourceTagSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting tag set (%s)", d.Id())

	octopus := m.(*client.Client)
	if err := octopus.TagSets.DeleteByID(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] tag set deleted")
	return nil
}

func resourceTagSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading tag set (%s)", d.Id()))

	octopus := m.(*client.Client)
	tagSet, err := octopus.TagSets.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "tag set")
	}

	if err := setTagSet(ctx, d, tagSet); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("tag set read (%s)", d.Id()))
	return nil
}

func resourceTagSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tagSet := expandTagSet(d)

	log.Printf("[INFO] updating tag set: %#v", tagSet)

	octopus := m.(*client.Client)
	existingTagSet, err := octopus.TagSets.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	tagSet.Tags = existingTagSet.Tags

	updatedTagSet, err := octopus.TagSets.Update(tagSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setTagSet(ctx, d, updatedTagSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] tag set updated (%s)", d.Id())
	return nil
}
