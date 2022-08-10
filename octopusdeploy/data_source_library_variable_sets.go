package octopusdeploy

import (
	"context"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
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
	query := variables.LibraryVariablesQuery{
		ContentType: d.Get("content_type").(string),
		IDs:         expandArray(d.Get("ids").([]interface{})),
		PartialName: d.Get("partial_name").(string),
		Skip:        d.Get("skip").(int),
		Take:        d.Get("take").(int),
	}

	client := m.(*client.Client)
	existingLibraryVariableSets, err := client.LibraryVariableSets.Get(query)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenedLibraryVariableSets := []interface{}{}
	for _, libraryVariableSet := range existingLibraryVariableSets.Items {
		flattenedLibraryVariableSets = append(flattenedLibraryVariableSets, flattenLibraryVariableSet(libraryVariableSet))
	}

	d.Set("library_variable_sets", flattenedLibraryVariableSets)
	d.SetId("Library Variables Sets " + time.Now().UTC().String())

	return nil
}
