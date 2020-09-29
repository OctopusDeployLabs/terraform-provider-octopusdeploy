package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGroup() *schema.Resource {
	return &schema.Resource{
		Description: "Creates a project group",
		Create:      resourceProjectGroupCreate,
		Read:        resourceProjectGroupRead,
		Update:      resourceProjectGroupUpdate,
		Delete:      resourceProjectGroupDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constDescription: {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func buildProjectGroupResource(d *schema.ResourceData) *model.ProjectGroup {
	name := d.Get(constName).(string)

	projectGroup := model.NewProjectGroup(name)

	if attr, ok := d.GetOk(constDescription); ok {
		projectGroup.Description = attr.(string)
	}

	return projectGroup
}

func resourceProjectGroupCreate(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := buildProjectGroupResource(d)

	apiClient := m.(*client.Client)
	resource, err := apiClient.ProjectGroups.Add(projectGroup)
	if err != nil {
		// return createResourceOperationError(errorCreatingProjectGroup, projectGroup.ID, err)
		return diag.FromErr(err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return diags
}

func resourceProjectGroupRead(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.ProjectGroups.GetByID(id)
	if err != nil {
		// return createResourceOperationError(errorReadingProjectGroup, id, err)
		diag.Errorf(err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constProjectGroup, m)

	d.Set(constName, resource.Name)
	d.Set(constDescription, resource.Description)

	return diag
}

func resourceProjectGroupUpdate(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	projectGroup := buildProjectGroupResource(d)
	projectGroup.ID = d.Id() // set ID so Octopus API knows which project group to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.ProjectGroups.Update(projectGroup)
	if err != nil {
		// return createResourceOperationError(errorUpdatingProjectGroup, d.Id(), err)
		diag.Errorf(err)
	}

	d.SetId(resource.ID)

	return diag
}

func resourceProjectGroupDelete(d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.ProjectGroups.DeleteByID(id)
	if err != nil {
		// return createResourceOperationError(errorDeletingProjectGroup, id, err)
		diag.Errorf(err)
	}

	d.SetId(constEmptyString)

	return diag
}
