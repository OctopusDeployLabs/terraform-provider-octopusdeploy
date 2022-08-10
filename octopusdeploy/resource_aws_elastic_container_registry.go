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

func resourceAwsElasticContainerRegistry() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAwsElasticContainerRegistryCreate,
		DeleteContext: resourceAwsElasticContainerRegistryDelete,
		Description:   "This resource manages an AWS Elastic Container Registry in Octopus Deploy.",
		Importer:      getImporter(),
		ReadContext:   resourceAwsElasticContainerRegistryRead,
		Schema:        getAwsElasticContainerRegistrySchema(),
		UpdateContext: resourceAwsElasticContainerRegistryUpdate,
	}
}

func resourceAwsElasticContainerRegistryCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	awsElasticContainerRegistry, err := expandAwsElasticContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] creating AWS Elastic Container Registry, %s", awsElasticContainerRegistry.GetName())

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(awsElasticContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, createdFeed.(*feeds.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	log.Printf("[INFO] AWS Elastic Container Registry created (%s)", d.Id())
	return nil
}

func resourceAwsElasticContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	log.Printf("[INFO] deleting AWS Elastic Container Registry (%s)", d.Id())

	client := m.(*client.Client)
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

	client := m.(*client.Client)
	feedResource, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "AWS Elastic Container Registry")
	}

	feedResource, err = feeds.ToFeed(feedResource.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	awsElasticContainerRegistry := feedResource.(*feeds.AwsElasticContainerRegistry)
	if err := setAwsElasticContainerRegistry(ctx, d, awsElasticContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS Elastic Container Registry read: %s", awsElasticContainerRegistry.GetID())
	return nil
}

func resourceAwsElasticContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	awsElasticContainerRegistry, err := expandAwsElasticContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] updating AWS Elastic Container Registry (%s)", awsElasticContainerRegistry.GetID())

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(awsElasticContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource, err := feeds.ToFeed(updatedFeed.(*feeds.FeedResource))
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, feedResource.(*feeds.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] AWS Elastic Container Registry updated (%s)", d.Id())
	return nil
}
