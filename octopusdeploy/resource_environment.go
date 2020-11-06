package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceEnvironment() *schema.Resource {
	resourceEnvironmentImporter := &schema.ResourceImporter{
		StateContext: schema.ImportStatePassthroughContext,
	}
	resourceEnvironmentSchema := map[string]*schema.Schema{
		constAllowDynamicInfrastructure: &schema.Schema{
			Optional: true,
			Type:     schema.TypeBool,
		},
		constDescription: &schema.Schema{
			Optional: true,
			Type:     schema.TypeString,
		},
		constName: &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
		constSortOrder: &schema.Schema{
			Computed: true,
			Type:     schema.TypeInt,
		},
		constUseGuidedFailure: &schema.Schema{
			Optional: true,
			Type:     schema.TypeBool,
		},
	}

	return &schema.Resource{
		CreateContext: resourceEnvironmentCreate,
		DeleteContext: resourceEnvironmentDelete,
		Importer:      resourceEnvironmentImporter,
		ReadContext:   resourceEnvironmentRead,
		Schema:        resourceEnvironmentSchema,
		UpdateContext: resourceEnvironmentUpdate,
	}
}

func buildEnvironmentResource(d *schema.ResourceData) *octopusdeploy.Environment {
	var name string
	if v, ok := d.GetOk(constName); ok {
		name = v.(string)
	}

	environment := octopusdeploy.NewEnvironment(name)

	if v, ok := d.GetOk(constAllowDynamicInfrastructure); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk(constDescription); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk(constSortOrder); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk(constUseGuidedFailure); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func resourceEnvironmentCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := buildEnvironmentResource(d)

	client := m.(*octopusdeploy.Client)
	createdEnvironment, err := client.Environments.Add(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	updateEnvironmentState(ctx, d, createdEnvironment)
	return nil
}

func resourceEnvironmentRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	environment, err := client.Environments.GetByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	updateEnvironmentState(ctx, d, environment)
	return nil
}

func resourceEnvironmentUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environment := buildEnvironmentResource(d)
	environment.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	updatedEnvironment, err := client.Environments.Update(environment)
	if err != nil {
		return diag.FromErr(err)
	}

	updateEnvironmentState(ctx, d, updatedEnvironment)
	return nil
}

func resourceEnvironmentDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Environments.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)
	return nil
}

func updateEnvironmentState(ctx context.Context, d *schema.ResourceData, environment *octopusdeploy.Environment) {
	d.Set(constAllowDynamicInfrastructure, environment.AllowDynamicInfrastructure)
	d.Set(constDescription, environment.Description)
	d.Set(constName, environment.Name)
	d.Set(constSortOrder, environment.SortOrder)
	d.Set(constUseGuidedFailure, environment.UseGuidedFailure)
	d.SetId(environment.ID)
}
