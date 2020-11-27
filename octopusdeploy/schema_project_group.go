package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandProjectGroup(d *schema.ResourceData) *octopusdeploy.ProjectGroup {
	name := d.Get("name").(string)

	projectGroup := octopusdeploy.NewProjectGroup(name)
	projectGroup.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		projectGroup.Description = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		projectGroup.EnvironmentIDs = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("retention_policy_id"); ok {
		projectGroup.RetentionPolicyID = v.(string)
	}

	return projectGroup
}

func flattenProjectGroup(projectGroup *octopusdeploy.ProjectGroup) map[string]interface{} {
	if projectGroup == nil {
		return nil
	}

	return map[string]interface{}{
		"description":         projectGroup.Description,
		"environments":        projectGroup.EnvironmentIDs,
		"id":                  projectGroup.GetID(),
		"name":                projectGroup.Name,
		"retention_policy_id": projectGroup.RetentionPolicyID,
	}
}

func getProjectGroupDataSchema() map[string]*schema.Schema {
	projectGroupSchema := getProjectGroupSchema()
	for _, field := range projectGroupSchema {
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
		"ids": {
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"partial_name": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"project_groups": {
			Computed: true,
			Elem:     &schema.Resource{Schema: projectGroupSchema},
			Type:     schema.TypeList,
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

func getProjectGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Computed: true,
			Elem: &schema.Schema{
				Type: schema.TypeString,
			},
			Type: schema.TypeList,
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
		"retention_policy_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}

func setProjectGroup(ctx context.Context, d *schema.ResourceData, projectGroup *octopusdeploy.ProjectGroup) {
	d.Set("description", projectGroup.Description)
	d.Set("environments", projectGroup.EnvironmentIDs)
	d.Set("name", projectGroup.Name)
	d.Set("retention_policy_id", projectGroup.RetentionPolicyID)

	d.SetId(projectGroup.GetID())
}
