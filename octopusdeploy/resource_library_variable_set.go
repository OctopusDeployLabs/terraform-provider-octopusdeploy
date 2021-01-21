package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceLibraryVariableSetCreate,
		DeleteContext: resourceLibraryVariableSetDelete,
		Importer:      getImporter(),
		ReadContext:   resourceLibraryVariableSetRead,
		Schema:        getLibraryVariableSetSchema(),
		UpdateContext: resourceLibraryVariableSetUpdate,
	}
}

func resourceLibraryVariableSetCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := expandLibraryVariableSet(d)

	client := m.(*octopusdeploy.Client)
	createdLibraryVariableSet, err := client.LibraryVariableSets.Add(libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdLibraryVariableSet.GetID())
	return resourceLibraryVariableSetRead(ctx, d, m)
}

func resourceLibraryVariableSetDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.LibraryVariableSets.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}

func resourceLibraryVariableSetRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	libraryVariableSet, err := client.LibraryVariableSets.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	setLibraryVariableSet(ctx, d, libraryVariableSet)
	return nil
}

func resourceLibraryVariableSetUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSet := expandLibraryVariableSet(d)

	client := m.(*octopusdeploy.Client)
	_, err := client.LibraryVariableSets.Update(libraryVariableSet)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceLibraryVariableSetRead(ctx, d, m)
}
