package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
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

	client := m.(*octopusdeploy.Client)
	createdDockerContainerRegistry, err := client.Feeds.Add(dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, createdDockerContainerRegistry.(*octopusdeploy.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdDockerContainerRegistry.GetID())

	log.Printf("[INFO] Docker container registry created (%s)", d.Id())
	return nil
}

func resourceDockerContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting Docker container registry (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
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

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] Docker container registry (%s) not found; deleting from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	dockerContainerRegistry := feedResource.(*octopusdeploy.DockerContainerRegistry)
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

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setDockerContainerRegistry(ctx, d, feedResource.(*octopusdeploy.DockerContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Docker container registry updated (%s)", d.Id())
	return nil
}
