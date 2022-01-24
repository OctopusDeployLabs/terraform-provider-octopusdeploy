package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandProjectGroup(d *schema.ResourceData) *octopusdeploy.ProjectGroup {
	name := d.Get("name").(string)

	projectGroup := octopusdeploy.NewProjectGroup(name)
	projectGroup.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		projectGroup.Description = v.(string)
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
		"id":                  projectGroup.GetID(),
		"name":                projectGroup.Name,
		"retention_policy_id": projectGroup.RetentionPolicyID,
	}
}

func getProjectGroupDataSchema() map[string]*schema.Schema {
	dataSchema := getProjectGroupSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"project_groups": {
			Computed:    true,
			Description: "A list of project groups that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"skip": getQuerySkip(),
		"take": getQueryTake(),
	}
}

func getProjectGroupSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"description": {
			Computed:    true,
			Description: "The description of this project group.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"id":   getIDSchema(),
		"name": getNameSchema(true),
		"retention_policy_id": {
			Computed:    true,
			Description: "The ID of the retention policy associated with this project group.",
			Optional:    true,
			Type:        schema.TypeString,
		},
	}
}

func setProjectGroup(ctx context.Context, d *schema.ResourceData, projectGroup *octopusdeploy.ProjectGroup) error {
	d.Set("description", projectGroup.Description)

	d.Set("name", projectGroup.Name)
	d.Set("retention_policy_id", projectGroup.RetentionPolicyID)

	d.SetId(projectGroup.GetID())

	return nil
}
