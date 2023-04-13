package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandRunbook(ctx context.Context, d *schema.ResourceData) *runbooks.Runbook {
	name := d.Get("name").(string)
	projectId := d.Get("project_id").(string)

	runbook := runbooks.NewRunbook(name, projectId)
	runbook.ID = d.Id()

	if v, ok := d.GetOk("description"); ok {
		runbook.Description = v.(string)
	}

	if v, ok := d.GetOk("runbook_process_id"); ok {
		runbook.RunbookProcessID = v.(string)
	}

	if v, ok := d.GetOk("published_runbook_snapshot_id"); ok {
		runbook.PublishedRunbookSnapshotID = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		runbook.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("multi_tenancy_mode"); ok {
		runbook.MultiTenancyMode = core.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("connectivity_policy"); ok {
		runbook.ConnectivityPolicy = expandConnectivityPolicy(v.([]interface{}))
	}

	if v, ok := d.GetOk("environment_scope"); ok {
		runbook.EnvironmentScope = v.(string)
	}

	if v, ok := d.GetOk("environments"); ok {
		runbook.Environments = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("default_guided_failure_mode"); ok {
		runbook.DefaultGuidedFailureMode = v.(string)
	}

	return runbook
}

func flattenRunbook(ctx context.Context, d *schema.ResourceData, runbook *runbooks.Runbook) map[string]interface{} {
	if runbook == nil {
		return nil
	}

	runbookMap := map[string]interface{}{
		"name":                          runbook.Name,
		"project_id":                    runbook.ProjectID,
		"description":                   runbook.Description,
		"runbook_process_id":            runbook.RunbookProcessID,
		"published_runbook_snapshot_id": runbook.PublishedRunbookSnapshotID,
		"space_id":                      runbook.SpaceID,
		"multi_tenancy_mode":            runbook.MultiTenancyMode,
		"connectivity_policy":           runbook.ConnectivityPolicy,
		"environment_scope":             runbook.EnvironmentScope,
		"environments":                  runbook.Environments,
		"default_guided_failure_mode":   runbook.DefaultGuidedFailureMode,
	}

	return runbookMap
}

func getRunbookDataSchema() map[string]*schema.Schema {
	dataSchema := getProjectSchema()
	setDataSchema(&dataSchema)

	return map[string]*schema.Schema{
		"id":           getDataSchemaID(),
		"ids":          getQueryIDs(),
		"partial_name": getQueryPartialName(),
		"runbooks": {
			Computed:    true,
			Description: "A list of runbooks that match the filter(s).",
			Elem:        &schema.Resource{Schema: dataSchema},
			Optional:    true,
			Type:        schema.TypeList,
		},
		"skip": getQuerySkip(),
		"take": getQueryTake(),
	}
}

func getRunbookSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"name": {
			Description:      "The name of the project in Octopus Deploy. This name must be unique.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"description": {
			Computed:    true,
			Description: "The description of this runbook.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"project_id": {
			Description: "The project that this runbook belongs to",
			Required:    true,
			Type:        schema.TypeString,
		},
		"runbook_process_id": {
			Description: "The runbook process ID",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"published_runbook_snapshot_id": {
			Description: "The published snapshot ID",
			Computed:    true,
			Type:        schema.TypeString,
		},
		"space_id": {
			Computed:         true,
			Description:      "The space ID associated with this project.",
			Optional:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotWhiteSpace),
		},
		"multi_tenancy_mode": getTenantedDeploymentSchema(),
		"connectivity_policy": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getConnectivityPolicySchema()},
			MaxItems: 1,
			Optional: true,
			Type:     schema.TypeList,
		},
		"environment_scope": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
		},
		"environments": {
			Computed: true,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Type:     schema.TypeList,
		},
		"default_guided_failure_mode": {
			Computed: true,
			Optional: true,
			Type:     schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
				"EnvironmentDefault",
				"Off",
				"On",
			}, false)),
		},
	}
}

func setRunbook(ctx context.Context, d *schema.ResourceData, runbook *runbooks.Runbook) error {
	d.Set("id", runbook.ID)
	d.Set("name", runbook.Name)
	d.Set("project_id", runbook.ProjectID)
	d.Set("description", runbook.Description)
	d.Set("runbook_process_id", runbook.RunbookProcessID)
	d.Set("published_runbook_snapshot_id", runbook.PublishedRunbookSnapshotID)
	d.Set("space_id", runbook.SpaceID)
	d.Set("multi_tenancy_mode", runbook.MultiTenancyMode)
	if err := d.Set("connectivity_policy", flattenConnectivityPolicy(runbook.ConnectivityPolicy)); err != nil {
		return fmt.Errorf("error setting connectivity_policy: %s", err)
	}
	d.Set("environment_scope", runbook.EnvironmentScope)
	d.Set("environments", runbook.Environments)
	d.Set("default_guided_failure_mode", runbook.DefaultGuidedFailureMode)

	return nil
}
