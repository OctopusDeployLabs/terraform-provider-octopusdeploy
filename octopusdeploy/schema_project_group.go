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

func flattenProjectGroup(ctx context.Context, d *schema.ResourceData, projectGroup *octopusdeploy.ProjectGroup) {
	d.Set("description", projectGroup.Description)
	d.Set("environments", projectGroup.EnvironmentIDs)
	d.Set("name", projectGroup.Name)
	d.Set("retention_policy_id", projectGroup.RetentionPolicyID)

	d.SetId(projectGroup.GetID())
}

func getProjectGroupDataSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": &schema.Schema{
			Required:     true,
			Type:         schema.TypeString,
			ValidateFunc: validation.StringIsNotEmpty,
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
		"name": &schema.Schema{
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
