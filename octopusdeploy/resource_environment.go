package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceEnvironment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		DeleteContext: resourceEnvironmentDelete,
		Importer:      getImporter(),
		ReadContext:   resourceEnvironmentRead,
		Schema:        getEnvironmentSchema(),
		UpdateContext: resourceEnvironmentUpdate,
	}
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := expandEnvironment(d)

	client := m.(*octopusdeploy.Client)
	createdEnvironment, err := client.Environments.Add(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenEnvironment(ctx, d, createdEnvironment)
	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	environment, err := client.Environments.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	flattenEnvironment(ctx, d, environment)
	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := expandEnvironment(d)

	client := m.(*octopusdeploy.Client)
	updatedEnvironment, err := client.Environments.Update(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	flattenEnvironment(ctx, d, updatedEnvironment)
	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Environments.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
