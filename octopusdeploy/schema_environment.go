package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
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
	environmentSchema := getEnvironmentSchema()
	for _, field := range environmentSchema {
		field.Computed = true
		field.Default = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
		field.ValidateFunc = nil
	}

	return map[string]*schema.Schema{
		"environments": {
			Computed: true,
			Elem:     &schema.Resource{Schema: environmentSchema},
			Type:     schema.TypeList,
		},
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"skip": {
			Default:  0,
			Type:     schema.TypeInt,
			Optional: true,
		},
		"take": {
			Default:  1,
			Type:     schema.TypeInt,
			Optional: true,
		},
	}
}

func getEnvironmentSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"allow_dynamic_infrastructure": {
			Optional: true,
			Type:     schema.TypeBool,
		},
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"name": {
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
		},
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

func setEnvironment(ctx context.Context, d *schema.ResourceData, environment *octopusdeploy.Environment) {
	d.Set("allow_dynamic_infrastructure", environment.AllowDynamicInfrastructure)
	d.Set("description", environment.Description)
	d.Set("name", environment.Name)
	d.Set("sort_order", environment.SortOrder)
	d.Set("use_guided_failure", environment.UseGuidedFailure)

	d.SetId(environment.GetID())
}
