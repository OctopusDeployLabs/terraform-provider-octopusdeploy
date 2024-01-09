package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceDockerContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDockerContainerRegistryCreate,
		DeleteContext: resourceDockerContainerRegistryDelete,
		Description:   "This resource manages a Docker Container Registry in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceDockerContainerRegistryRead,
		Schema:        getDockerContainerRegistrySchema(),
		UpdateContext: resourceDockerContainerRegistryUpdate,
	}
}

func resourceDockerContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dockerContainerRegistry, err := expandDockerContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("creating Docker container registry, %s", dockerContainerRegistry.GetName()))

	client := m.(*client.Client)
	createdDockerContainerRegistry, err := feeds.Add(client, dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, createdDockerContainerRegistry.(*feeds.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDockerContainerRegistry.GetID())

	tflog.Info(ctx, fmt.Sprintf("Docker container registry created (%s)", d.Id()))
	return nil
}

func resourceDockerContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting Docker container registry (%s)", d.Id()))

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "Docker container registry deleted")
	return nil
}

func resourceDockerContainerRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading Docker container registry (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := feeds.GetByID(client, d.Get("space_id").(string), d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Docker container registry")
	}

	dockerContainerRegistry := feed.(*feeds.DockerContainerRegistry)
	if err := setDockerContainerRegistry(ctx, d, dockerContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Docker container registry read (%s)", dockerContainerRegistry.GetID()))
	return nil
}

func resourceDockerContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandDockerContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating Docker container registry (%s)", feed.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := feeds.Update(client, feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, updatedFeed.(*feeds.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("Docker container registry updated (%s)", d.Id()))
	return nil
}
