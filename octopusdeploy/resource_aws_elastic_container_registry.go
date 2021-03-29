package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAwsElasticContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsElasticContainerRegistryCreate,
		DeleteContext: resourceAwsElasticContainerRegistryDelete,
		Description:   "This resource manages a AWS Elastic Container Registry in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAwsElasticContainerRegistryRead,
		Schema:        getAwsElasticContainerRegistrySchema(),
		UpdateContext: resourceAwsElasticContainerRegistryUpdate,
	}
}

func resourceAwsElasticContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	dockerContainerRegistry := expandAwsElasticContainerRegistry(d)

	log.Printf("[INFO] creating AWS Elastic Container Registry: %#v", dockerContainerRegistry)

	client := m.(*octopusdeploy.Client)
	createdAwsElasticContainerRegistry, err := client.Feeds.Add(dockerContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, createdAwsElasticContainerRegistry.(*octopusdeploy.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdAwsElasticContainerRegistry.GetID())

	log.Printf("[INFO] AWS Elastic Container Registry created (%s)", d.Id())
	return nil
}

func resourceAwsElasticContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting AWS Elastic Container Registry (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	log.Printf("[INFO] AWS Elastic Container Registry deleted")
	return nil
}

func resourceAwsElasticContainerRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] reading AWS Elastic Container Registry (%s)", d.Id())

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		apiError := err.(*octopusdeploy.APIError)
		if apiError.StatusCode == 404 {
			log.Printf("[INFO] AWS Elastic Container Registry (%s) not found; deleting from state", d.Id())
			d.SetId("")
			return nil
		}
		return diag.FromErr(err)
	}

	feedResource, err = octopusdeploy.ToFeed(feedResource.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	dockerContainerRegistry := feedResource.(*octopusdeploy.AwsElasticContainerRegistry)
	if err := setAwsElasticContainerRegistry(ctx, d, dockerContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS Elastic Container Registry read: %#v", dockerContainerRegistry)
	return nil
}

func resourceAwsElasticContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := expandAwsElasticContainerRegistry(d)

	log.Printf("[INFO] updating AWS Elastic Container Registry: %#v", feed)

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := octopusdeploy.ToFeed(updatedFeed.(*octopusdeploy.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, feedResource.(*octopusdeploy.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS Elastic Container Registry updated (%s)", d.Id())
	return nil
}
