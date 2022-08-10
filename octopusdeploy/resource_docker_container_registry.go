package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
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

	log.Printf("[INFO] creating Docker container registry, %s", dockerContainerRegistry.GetName())

	client := m.(*client.Client)
	createdDockerContainerRegistry, err := client.Feeds.Add(dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, createdDockerContainerRegistry.(*feeds.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDockerContainerRegistry.GetID())

	log.Printf("[INFO] Docker container registry created (%s)", d.Id())
	return nil
}

func resourceDockerContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Docker container registry (%s)", d.Id())

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] Docker container registry deleted")
	return nil
}

func resourceDockerContainerRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading Docker container registry (%s)", d.Id())

	client := m.(*client.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "Docker container registry")
	}

	feedResource, err = feeds.ToFeed(feedResource.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	dockerContainerRegistry := feedResource.(*feeds.DockerContainerRegistry)
	if err := setDockerContainerRegistry(ctx, d, dockerContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Docker container registry read (%s)", dockerContainerRegistry.GetID())
	return nil
}

func resourceDockerContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed, err := expandDockerContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating Docker container registry (%s)", feed.GetID())

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := feeds.ToFeed(updatedFeed.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, feedResource.(*feeds.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Docker container registry updated (%s)", d.Id())
	return nil
}
