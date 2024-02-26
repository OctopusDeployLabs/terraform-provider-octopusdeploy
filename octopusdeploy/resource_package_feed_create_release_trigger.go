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

func buildPackageFeedCreateReleaseTriggerResource(d *schema.ResourceData, client *client.Client, feedCategory string) (*triggers.ProjectTrigger, error) {
	name := d.Get("name").(string)
	spaceId := d.Get("space_id").(string)
	projectId := d.Get("project_id").(string)
	channelId := d.Get("channel_id").(string)

	isDisabled := false
	if v, ok := d.GetOk("is_disabled"); ok {
		isDisabled = v.(bool)
	}

	flattenedPackages := d.Get("package")
	packages := expandDeploymentActionPackages(flattenedPackages)

	action := actions.NewCreateReleaseAction(channelId)
	filter := filters.NewFeedTriggerFilter(feedCategory, packages)

	project, err := projects.GetByID(client, spaceId, projectId)
	if err != nil {
		return nil, err
	}

	createReleaseTrigger := triggers.NewProjectTrigger(name, "", isDisabled, project, action, filter)

	return createReleaseTrigger, nil
}

func resourcePackageFeedCreateReleaseTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}, feedCategory string) diag.Diagnostics {
	client := m.(*client.Client)

	projectTrigger, err := buildPackageFeedCreateReleaseTriggerResource(d, client, feedCategory)
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

func resourcePackageFeedCreateReleaseTriggerRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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
	filter := projectTrigger.Filter.(*filters.FeedTriggerFilter)

	d.Set("name", projectTrigger.Name)
	d.Set("space_id", projectTrigger.SpaceID)
	d.Set("project_id", projectTrigger.ProjectID)
	d.Set("is_disabled", projectTrigger.IsDisabled)
	d.Set("channel_id", action.ChannelID)
	d.Set("package", flattenDeploymentActionPackages(filter.Packages))

	return nil
}

func resourcePackageFeedCreateReleaseTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}, feedCategory string) diag.Diagnostics {
	client := m.(*client.Client)
	projectTrigger, err := buildPackageFeedCreateReleaseTriggerResource(d, client, feedCategory)
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

func resourcePackageFeedCreateReleaseTriggerDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*client.Client)
	err := client.ProjectTriggers.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
