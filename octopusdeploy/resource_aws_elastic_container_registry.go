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
	feed, err := expandAwsElasticContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("creating AWS Elastic Container Registry, %s", feed.GetName()))

	client := m.(*client.Client)
	createdFeed, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, createdFeed.(*feeds.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(createdFeed.GetID())

	tflog.Info(ctx, fmt.Sprintf("AWS Elastic Container Registry created (%s)", d.Id()))
	return nil
}

func resourceAwsElasticContainerRegistryDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("deleting AWS Elastic Container Registry (%s)", d.Id()))

	client := m.(*client.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	tflog.Info(ctx, "AWS Elastic Container Registry deleted")
	return nil
}

func resourceAwsElasticContainerRegistryRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	tflog.Info(ctx, fmt.Sprintf("reading AWS Elastic Container Registry (%s)", d.Id()))

	client := m.(*client.Client)
	feed, err := client.Feeds.GetByID(d.Id())
	if err != nil {
		return errors.ProcessApiError(ctx, d, err, "AWS Elastic Container Registry")
	}

	awsElasticContainerRegistry := feed.(*feeds.AwsElasticContainerRegistry)
	if err := setAwsElasticContainerRegistry(ctx, d, awsElasticContainerRegistry); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("AWS Elastic Container Registry read: %s", awsElasticContainerRegistry.GetID()))
	return nil
}

func resourceAwsElasticContainerRegistryUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	awsElasticContainerRegistry, err := expandAwsElasticContainerRegistry(d)
	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("updating AWS Elastic Container Registry (%s)", awsElasticContainerRegistry.GetID()))

	client := m.(*client.Client)
	updatedFeed, err := client.Feeds.Update(awsElasticContainerRegistry)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := setAwsElasticContainerRegistry(ctx, d, updatedFeed.(*feeds.AwsElasticContainerRegistry)); err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, fmt.Sprintf("AWS Elastic Container Registry updated (%s)", d.Id()))
	return nil
}
