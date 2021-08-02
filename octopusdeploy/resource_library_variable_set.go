package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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
	}
}

func resourceLibraryVariableSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := expandLibraryVariableSet(d)

	log.Printf("[INFO] creating library variable set: %#v", libraryVariableSet)

	client := m.(*octopusdeploy.Client)
	createdLibraryVariableSet, err := client.LibraryVariableSets.Add(libraryVariableSet)
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

	client := m.(*octopusdeploy.Client)
	err := client.LibraryVariableSets.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] library variable set deleted (%s)", d.Id())
	d.SetId("")
	return nil
}

func resourceLibraryVariableSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading library variable set (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	libraryVariableSet, err := client.LibraryVariableSets.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] library variable set (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
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

	client := m.(*octopusdeploy.Client)
	updatedLibraryVariableSet, err := client.LibraryVariableSets.Update(libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setLibraryVariableSet(ctx, d, updatedLibraryVariableSet); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] library variable set updated (%s)", d.Id())
	return nil
}
