package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/libraryvariablesets"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLibraryVariableSetCreate,
		DeleteContext: resourceLibraryVariableSetDelete,
		Description:   "This resource manages library variable sets in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceLibraryVariableSetRead,
		Schema:        getLibraryVariableSetSchema(),
		UpdateContext: resourceLibraryVariableSetUpdate,
		CustomizeDiff: fixTemplateIds,
	}
}

// fixTemplateIds uses the suggestion from https://github.com/hashicorp/terraform/issues/18863
// to ensure that the template_ids field has keys to match the list of template names.
func fixTemplateIds(ctx context.Context, d *schema.ResourceDiff, meta interface{}) error {
	templates := d.Get("template")
	templateIds := map[string]string{}
	if templates != nil {
		for _, t := range templates.([]interface{}) {
			template := t.(map[string]interface{})
			templateIds[template["name"].(string)] = template["id"].(string)
		}
	}
	if err := d.SetNew("template_ids", templateIds); err != nil {
		return err
	}

	return nil
}

func resourceLibraryVariableSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := expandLibraryVariableSet(d)

	log.Printf("[INFO] creating library variable set: %#v", libraryVariableSet)

	client := m.(*client.Client)

	createdLibraryVariableSet, err := libraryvariablesets.Add(client, libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLibraryVariableSet(ctx, d, createdLibraryVariableSet); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdLibraryVariableSet.GetID())

	log.Printf("[INFO] library variable set created (%s)", d.Id())
	return nil
}

func resourceLibraryVariableSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting library variable set (%s)", d.Id())

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	client := m.(*client.Client)
	err := libraryvariablesets.DeleteByID(client, spaceID, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] library variable set deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceLibraryVariableSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading library variable set (%s)", d.Id())

	var spaceID string
	if v, ok := d.GetOk("space_id"); ok {
		spaceID = v.(string)
	}

	client := m.(*client.Client)
	libraryVariableSet, err := libraryvariablesets.GetByID(client, spaceID, d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "library variable set")
	}

	if err := setLibraryVariableSet(ctx, d, libraryVariableSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] library variable set read (%s)", d.Id())
	return nil
}

func resourceLibraryVariableSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating library variable set (%s)", d.Id())

	libraryVariableSet := expandLibraryVariableSet(d)

	client := m.(*client.Client)
	updatedLibraryVariableSet, err := libraryvariablesets.Update(client, libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLibraryVariableSet(ctx, d, updatedLibraryVariableSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] library variable set updated (%s)", d.Id())
	return nil
}
