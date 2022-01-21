package octopusdeploy

import (
	"context"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDeploymentProcess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeploymentProcessCreate,
		DeleteContext: resourceDeploymentProcessDelete,
		Description:   "This resource manages deployment processes in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceDeploymentProcessRead,
		Schema:        getDeploymentProcessSchema(),
		UpdateContext: resourceDeploymentProcessUpdate,
	}
}

func getDeploymentProcessSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"id": getIDSchema(),
		"branch": {
			Computed:    true,
			Description: "The branch name associated with this deployment process (i.e. `main`). This value is optional and only applies to associated projects that are stored in version control.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"last_snapshot_id": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"project_id": {
			Description: "The project ID associated with this deployment process.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"step":     getDeploymentStepSchema(),
		"version": {
			Computed:    true,
			Description: "The version number of this deployment process.",
			Optional:    true,
			Type:        schema.TypeInt,
		},
	}
}

func resourceDeploymentProcessCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	deploymentProcess := expandDeploymentProcess(d, client)

	log.Printf("[INFO] creating deployment process: %#v", deploymentProcess)

	project, err := client.Projects.GetByID(deploymentProcess.ProjectID)
	if err != nil {
		return diag.FromErr(err)
	}

	var current *octopusdeploy.DeploymentProcess
	if project.PersistenceSettings != nil && project.PersistenceSettings.GetType() == "VersionControlled" {
		current, err = client.DeploymentProcesses.Get(project, deploymentProcess.Branch)
		if err != nil {
			return diag.FromErr(err)
		}
	} else {
		current, err = client.DeploymentProcesses.GetByID(project.DeploymentProcessID)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	deploymentProcess.ID = current.ID
	deploymentProcess.Links = current.Links
	deploymentProcess.Version = current.Version

	createdDeploymentProcess, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDeploymentProcess(ctx, d, createdDeploymentProcess); err != nil {
		return diag.FromErr(err)
	}

	id := createdDeploymentProcess.GetID()
	if project.PersistenceSettings != nil && project.PersistenceSettings.GetType() == "VersionControlled" {
		id = "deploymentprocess-" + createdDeploymentProcess.ProjectID + "-" + deploymentProcess.Branch
	}

	d.SetId(id)

	log.Printf("[INFO] deployment process created (%s)", d.Id())
	return nil
}

func resourceDeploymentProcessDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting deployment process (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	current, err := client.DeploymentProcesses.GetByID(d.Id())
	if err == nil {
		deploymentProcess := &octopusdeploy.DeploymentProcess{
			Version: current.Version,
		}
		deploymentProcess.Links = current.Links
		deploymentProcess.ID = d.Id()

		_, err = client.DeploymentProcesses.Update(deploymentProcess)
		if err != nil {
			return diag.FromErr(err)
		}

		d.SetId("")
		log.Printf("[INFO] deployment process deleted")
		return nil
	}

	r, _ := regexp.Compile(`Projects-\d+`)
	projectID := r.FindString(d.Id())

	project, err := client.Projects.GetByID(projectID)
	if err != nil {
		return diag.FromErr(err)
	}

	r, _ = regexp.Compile(`\d+-\w+`)
	gitRef := strings.SplitAfter(r.FindString(d.Id()), "-")[1]

	current, err = client.DeploymentProcesses.Get(project, gitRef)
	if err != nil {
		return diag.FromErr(err)
	}

	deploymentProcess := &octopusdeploy.DeploymentProcess{
		Version: current.Version,
	}
	deploymentProcess.Links = current.Links
	deploymentProcess.ID = d.Id()

	_, err = client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	log.Printf("[INFO] deployment process deleted")
	return nil
}

func resourceDeploymentProcessRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading deployment process (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentProcess, err := client.DeploymentProcesses.GetByID(d.Id())
	if err == nil {
		if err := setDeploymentProcess(ctx, d, deploymentProcess); err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] deployment process read (%s)", d.Id())
		return nil
	}

	r, _ := regexp.Compile(`Projects-\d+`)
	projectID := r.FindString(d.Id())

	project, err := client.Projects.GetByID(projectID)
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == http.StatusNotFound {
				log.Printf("[INFO] project (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	r, _ = regexp.Compile(`\d+-\w+`)
	gitRef := strings.SplitAfter(r.FindString(d.Id()), "-")[1]

	deploymentProcess, err = client.DeploymentProcesses.Get(project, gitRef)
	if err == nil {
		if err := setDeploymentProcess(ctx, d, deploymentProcess); err != nil {
			return diag.FromErr(err)
		}

		log.Printf("[INFO] deployment process read (%s)", d.Id())
		return nil
	}

	log.Printf("[INFO] deployment process (%s) not found; deleting from state", d.Id())
	d.SetId("")
	return nil
}

func resourceDeploymentProcessUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] updating deployment process (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	deploymentProcess := expandDeploymentProcess(d, client)
	current, err := client.DeploymentProcesses.GetByID(d.Id())
	if err != nil {
		r, _ := regexp.Compile(`Projects-\d+`)
		projectID := r.FindString(d.Id())

		project, err := client.Projects.GetByID(projectID)
		if err != nil {
			return diag.FromErr(err)
		}

		r, _ = regexp.Compile(`\d+-\w+`)
		gitRef := strings.SplitAfter(r.FindString(d.Id()), "-")[1]

		current, err = client.DeploymentProcesses.Get(project, gitRef)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	deploymentProcess.Links = current.Links
	deploymentProcess.Version = current.Version

	updatedDeploymentProcess, err := client.DeploymentProcesses.Update(deploymentProcess)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDeploymentProcess(ctx, d, updatedDeploymentProcess); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] deployment process updated (%s)", d.Id())
	return nil
}
