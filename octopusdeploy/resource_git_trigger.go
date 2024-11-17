package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actions"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/filters"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/triggers"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceGitTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGitTriggerCreate,
		DeleteContext: resourceGitTriggerDelete,
		Description:   "This resource manages Git repository triggers in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceGitTriggerRead,
		Schema:        getGitTriggerSchema(),
		UpdateContext: resourceGitTriggerUpdate,
	}
}

func buildGitTriggerResource(d *schema.ResourceData, client *client.Client) (*triggers.ProjectTrigger, error) {
	name := d.Get("name").(string)
	spaceId := d.Get("space_id").(string)
	projectId := d.Get("project_id").(string)
	channelId := d.Get("channel_id").(string)

	isDisabled := false
	if v, ok := d.GetOk("is_disabled"); ok {
		isDisabled = v.(bool)
	}

	flattenedGitTriggerSources := d.Get("sources")
	gitTriggerSources := expandGitTriggerSources(flattenedGitTriggerSources)

	action := actions.NewCreateReleaseAction(channelId)
	filter := filters.NewGitTriggerFilter(gitTriggerSources)

	project, err := projects.GetByID(client, spaceId, projectId)
	if err != nil {
		return nil, err
	}

	createReleaseTrigger := triggers.NewProjectTrigger(name, "", isDisabled, project, action, filter)

	return createReleaseTrigger, nil
}

func resourceGitTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)

	projectTrigger, err := buildGitTriggerResource(d, client)
	if err != nil {
		return diag.FromErr(err)
	}

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

func resourceGitTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*client.Client)
	projectTrigger, err := client.ProjectTriggers.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if projectTrigger == nil {
		d.SetId("")
		return nil
	}

	action := projectTrigger.Action.(*actions.CreateReleaseAction)
	filter := projectTrigger.Filter.(*filters.GitTriggerFilter)

	d.Set("name", projectTrigger.Name)
	d.Set("space_id", projectTrigger.SpaceID)
	d.Set("project_id", projectTrigger.ProjectID)
	d.Set("is_disabled", projectTrigger.IsDisabled)
	d.Set("channel_id", action.ChannelID)
	d.Set("sources", flattenGitTriggerSources(filter.Sources))

	return nil
}

func resourceGitTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	projectTrigger, err := buildGitTriggerResource(d, client)
	if err != nil {
		return diag.FromErr(err)
	}
	projectTrigger.ID = d.Id() // set ID so Octopus API knows which project trigger to update

	resource, err := client.ProjectTriggers.Update(projectTrigger)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceGitTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	err := client.ProjectTriggers.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
