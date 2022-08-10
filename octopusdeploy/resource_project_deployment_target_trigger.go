package octopusdeploy

import (
	"context"
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectDeploymentTargetTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectDeploymentTargetTriggerCreate,
		DeleteContext: resourceProjectDeploymentTargetTriggerDelete,
		Importer:      getImporter(),
		ReadContext:   resourceProjectDeploymentTargetTriggerRead,
		Schema:        getProjectDeploymentTargetTriggerSchema(),
		UpdateContext: resourceProjectDeploymentTargetTriggerUpdate,
	}
}

func buildProjectDeploymentTargetTriggerResource(d *schema.ResourceData) (*triggers.ProjectTrigger, error) {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	shouldRedeploy := d.Get("should_redeploy").(bool)

	action := actions.NewAutoDeployAction(shouldRedeploy)
	filter := filters.NewDeploymentTargetFilter([]string{}, []string{}, []string{}, []string{})

	if attr, ok := d.GetOk("event_groups"); ok {
		eventGroups := getSliceFromTerraformTypeList(attr)

		// need to validate here "ValidateFunc is not yet supported on lists or sets."
		validValues := []string{
			"Machine",
			"MachineCritical",
			"MachineAvailableForDeployment",
			"MachineUnavailableForDeployment",
			"MachineHealthChanged",
		}

		if invalidValue, ok := validateAllSliceItemsInSlice(eventGroups, validValues); !ok {
			return nil, fmt.Errorf("invalid value for event_groups. %s not in %v", invalidValue, validValues)
		}

		filter.EventGroups = append(filter.EventGroups, eventGroups...)
	}

	if attr, ok := d.GetOk("event_categories"); ok {
		eventCategories := getSliceFromTerraformTypeList(attr)

		// need to validate here "ValidateFunc is not yet supported on lists or sets."
		validValues := []string{
			"MachineCleanupFailed",
			"MachineAdded",
			"MachineDeploymentRelatedPropertyWasUpdated",
			"MachineDisabled",
			"MachineEnabled",
			"MachineHealthy",
			"MachineUnavailable",
			"MachineUnhealthy",
			"MachineHasWarnings",
		}

		if invalidValue, ok := validateAllSliceItemsInSlice(eventCategories, validValues); !ok {
			return nil, fmt.Errorf("invalid value for event_categories. %s not in %v", invalidValue, validValues)
		}

		filter.EventCategories = append(filter.EventCategories, eventCategories...)
	}

	if attr, ok := d.GetOk("roles"); ok {
		filter.Roles = getSliceFromTerraformTypeList(attr)
	}

	if attr, ok := d.GetOk("environment_ids"); ok {
		filter.Environments = getSliceFromTerraformTypeList(attr)
	}

	deploymentTargetTrigger := triggers.NewProjectTrigger(name, "", false, projectID, action, filter)

	return deploymentTargetTrigger, nil
}

func resourceProjectDeploymentTargetTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectTrigger, err := buildProjectDeploymentTargetTriggerResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*client.Client)
	resource, err := client.ProjectTriggers.Add(projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceProjectDeploymentTargetTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*client.Client)
	resource, err := client.ProjectTriggers.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource("project_trigger", m)

	action := resource.Action.(*actions.AutoDeployAction)
	filter := resource.Filter.(*filters.DeploymentTargetFilter)

	d.Set("environment_ids", filter.Environments)
	d.Set("event_groups", filter.EventGroups)
	d.Set("event_categories", filter.EventCategories)
	d.Set("name", resource.Name)
	d.Set("roles", filter.Roles)
	d.Set("should_redeploy", action.ShouldRedeploy)

	return nil
}

func resourceProjectDeploymentTargetTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectTrigger, err := buildProjectDeploymentTargetTriggerResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	projectTrigger.ID = d.Id() // set ID so Octopus API knows which project trigger to update

	client := m.(*client.Client)
	resource, err := client.ProjectTriggers.Update(*projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceProjectDeploymentTargetTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	err := client.ProjectTriggers.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
