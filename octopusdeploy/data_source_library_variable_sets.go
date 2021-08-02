package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLibraryVariableSet() *schema.Resource {
	return &schema.Resource{
		Description: "Provides information about existing library variable sets.",
		ReadContext: dataSourceLibraryVariableSetReadByName,
		Schema:      getLibraryVariableSetDataSchema(),
	}
}

func dataSourceLibraryVariableSetReadByName(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	query := octopusdeploy.LibraryVariablesQuery{
		ContentType: d.Get("content_type").(string),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*octopusdeploy.Client)
	libraryVariableSets, err := client.LibraryVariableSets.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedLibraryVariableSets := []interface{}{}
	for _, libraryVariableSet := range libraryVariableSets.Items {
		flattenedLibraryVariableSets = append(flattenedLibraryVariableSets, flattenLibraryVariableSet(libraryVariableSet))
	}

	d.Set("library_variable_sets", flattenedLibraryVariableSets)
	d.SetId("Library Variables Sets " + time.Now().UTC().String())

	return nil
}
