package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandEnvironment(d *schema.ResourceData) *octopusdeploy.Environment {
	name := d.Get("name").(string)

	environment := octopusdeploy.NewEnvironment(name)
	environment.ID = d.Id()

	if v, ok := d.GetOk("allow_dynamic_infrastructure"); ok {
		environment.AllowDynamicInfrastructure = v.(bool)
	}

	if v, ok := d.GetOk("description"); ok {
		environment.Description = v.(string)
	}

	if v, ok := d.GetOk("sort_order"); ok {
		environment.SortOrder = v.(int)
	}

	if v, ok := d.GetOk("use_guided_failure"); ok {
		environment.UseGuidedFailure = v.(bool)
	}

	return environment
}

func flattenEnvironment(environment *octopusdeploy.Environment) map[string]interface{} {
	if environment == nil {
		return nil
	}

	return map[string]interface{}{
		"allow_dynamic_infrastructure": environment.AllowDynamicInfrastructure,
		"description":                  environment.Description,
		"id":                           environment.GetID(),
		"name":                         environment.Name,
		"sort_order":                   environment.SortOrder,
		"use_guided_failure":           environment.UseGuidedFailure,
	}
}
func getEnvironmentDataSchema() map[string]*schema.Schema {
	dataSchema := getEnvironmentSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"environment": {
			Computed:    true,
			Description: "A list of environments that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"name":         getQueryName(),
		"partial_name": getQueryPartialName(),
		"skip":         getQuerySkip(),
		"take":         getQueryTake(),
	}
}

func getEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_dynamic_infrastructure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": getDescriptionSchema(),
		"id":          getIDSchema(),
		"name":        getNameSchema(true),
		"sort_order": {
			Computed: true,
			Type:     schema.TypeInt,
		},
		"use_guided_failure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
	}
}

func setEnvironment(ctx context.Context, d *schema.ResourceData, environment *octopusdeploy.Environment) error {
	d.Set("allow_dynamic_infrastructure", environment.AllowDynamicInfrastructure)
	d.Set("description", environment.Description)
	d.Set("name", environment.Name)
	d.Set("sort_order", environment.SortOrder)
	d.Set("use_guided_failure", environment.UseGuidedFailure)

	d.SetId(environment.GetID())

	return nil
}
