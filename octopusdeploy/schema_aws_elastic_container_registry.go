package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandAwsElasticContainerRegistry(d *schema.ResourceData) (*feeds.AwsElasticContainerRegistry, error) {
	accessKey := d.Get("access_key").(string)
	name := d.Get("name").(string)
	secretKey := core.NewSensitiveValue(d.Get("secret_key").(string))
	region := d.Get("region").(string)

	feed, err := feeds.NewAwsElasticContainerRegistry(name, accessKey, secretKey, region)
	if err != nil {
		return nil, err
	}

	feed.ID = d.Id()

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		feed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("space_id"); ok {
		feed.SpaceID = v.(string)
	}

	return feed, nil
}

func getAwsElasticContainerRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"access_key": {
			Description: "The AWS access key to use when authenticating against Amazon Web Services.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"id": {
			Computed:    true,
			Description: "The unique ID for this feed.",
			Optional:    true,
			Type:        schema.TypeString,
		},
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"region": {
			Description: "The AWS region where the registry resides.",
			Required:    true,
			Type:        schema.TypeString,
		},
		"secret_key": {
			Description: "The AWS secret key to use when authenticating against Amazon Web Services.",
			Required:    true,
			Sensitive:   true,
			Type:        schema.TypeString,
		},
		"space_id": {
			Computed:    true,
			Description: "The space ID associated with this feed.",
			Optional:    true,
			Type:        schema.TypeString,
			ForceNew:    true,
		},
	}
}

func setAwsElasticContainerRegistry(ctx context.Context, d *schema.ResourceData, feed *feeds.AwsElasticContainerRegistry) error {
	d.Set("access_key", feed.AccessKey)
	d.Set("name", feed.Name)
	d.Set("space_id", feed.SpaceID)
	d.Set("region", feed.Region)

	if err := d.Set("package_acquisition_location_options", feed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(feed.GetID())

	return nil
}
