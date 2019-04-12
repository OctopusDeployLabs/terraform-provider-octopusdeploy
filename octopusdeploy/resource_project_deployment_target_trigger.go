package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceProjectDeploymentTargetTrigger() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectDeploymentTargetTriggerCreate,
		Read:   resourceProjectDeploymentTargetTriggerRead,
		Update: resourceProjectDeploymentTargetTriggerUpdate,
		Delete: resourceProjectDeploymentTargetTriggerDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the trigger.",
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The project_id of the Project to attach the trigger to.",
			},
			"should_redeploy": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Enable to re-deploy to the deployment targets even if they are already up-to-date with the current deployment.",
			},
			"event_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Apply event group filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			},
			"event_categories": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Apply event category filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			},
			"roles": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Apply event role filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			},
			"environment_ids": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional:    true,
				Description: "Apply environment id filters to restrict which deployment targets will actually cause the trigger to fire, and consequently, which deployment targets will be automatically deployed to.",
			},
		},
	}
}

func buildProjectDeploymentTargetTriggerResource(d *schema.ResourceData) (*octopusdeploy.ProjectTrigger, error) {
	name := d.Get("name").(string)
	projectID := d.Get("project_id").(string)
	shouldRedeploy := d.Get("should_redeploy").(bool)

	deploymentTargetTrigger := octopusdeploy.NewProjectDeploymentTargetTrigger(name, projectID, shouldRedeploy, nil, nil, nil)

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
			return nil, fmt.Errorf("Invalid value for event_groups. %s not in %v", invalidValue, validValues)
		}

		deploymentTargetTrigger.AddEventGroups(eventGroups)
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
			return nil, fmt.Errorf("Invalid value for event_categories. %s not in %v", invalidValue, validValues)
		}

		deploymentTargetTrigger.AddEventCategories(eventCategories)
	}

	if attr, ok := d.GetOk("roles"); ok {
		deploymentTargetTrigger.Filter.Roles = getSliceFromTerraformTypeList(attr)
	}

	if attr, ok := d.GetOk("environment_ids"); ok {
		deploymentTargetTrigger.Filter.EnvironmentIds = getSliceFromTerraformTypeList(attr)
	}

	return deploymentTargetTrigger, nil
}

func resourceProjectDeploymentTargetTriggerCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	deploymentTargetTrigger, err := buildProjectDeploymentTargetTriggerResource(d)

	if err != nil {
		return err
	}

	createdProjectDeploymentTargetTrigger, err := client.ProjectTrigger.Add(deploymentTargetTrigger)

	if err != nil {
		return fmt.Errorf("error creating project deployment target trigger: %s", err.Error())
	}

	d.SetId(createdProjectDeploymentTargetTrigger.ID)
	return nil
}

func resourceProjectDeploymentTargetTriggerRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectTriggerID := d.Id()

	projectTrigger, err := client.ProjectTrigger.Get(projectTriggerID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading project trigger id %s: %s", projectTrigger.ID, err.Error())
	}

	log.Printf("[DEBUG] project trigger: %v", m)
	d.Set("name", projectTrigger.Name)
	d.Set("should_redeploy", projectTrigger.Action.ShouldRedeployWhenMachineHasBeenDeployedTo)
	d.Set("event_groups", projectTrigger.Filter.EventGroups)
	d.Set("event_categories", projectTrigger.Filter.EventCategories)
	d.Set("roles", projectTrigger.Filter.Roles)
	d.Set("environment_ids", projectTrigger.Filter.EnvironmentIds)
	return nil
}

func resourceProjectDeploymentTargetTriggerUpdate(d *schema.ResourceData, m interface{}) error {
	deploymentTargetTrigger, err := buildProjectDeploymentTargetTriggerResource(d)

	if err != nil {
		return err
	}

	deploymentTargetTrigger.ID = d.Id() // set deploymenttrigger struct ID so octopus knows which to update

	client := m.(*octopusdeploy.Client)

	updatedProjectTrigger, err := client.ProjectTrigger.Update(deploymentTargetTrigger)

	if err != nil {
		return fmt.Errorf("error updating project trigger id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedProjectTrigger.ID)
	return nil
}

func resourceProjectDeploymentTargetTriggerDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	projectTriggerID := d.Id()

	err := client.ProjectTrigger.Delete(projectTriggerID)

	if err != nil {
		return fmt.Errorf("error deleting project trigger id %s: %s", projectTriggerID, err.Error())
	}

	d.SetId("")
	return nil
}
