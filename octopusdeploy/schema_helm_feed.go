package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandHelmFeed(d *schema.ResourceData) (*octopusdeploy.HelmFeed, error) {
	name := d.Get("name").(string)

	helmFeed, err := octopusdeploy.NewHelmFeed(name)
	if err != nil {
		return nil, err
	}

	helmFeed.ID = d.Id()

	if v, ok := d.GetOk("feed_uri"); ok {
		helmFeed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		helmFeed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		helmFeed.Password = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		helmFeed.Username = v.(string)
	}

	return helmFeed, nil
}

func getHelmFeedSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"feed_uri": {
			Required: true,
			Type:     schema.TypeString,
		},
		"id": getIDSchema(),
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
		"password": getPasswordSchema(false),
		"space_id": getSpaceIDSchema(),
		"username": getUsernameSchema(false),
	}
}

func setHelmFeed(ctx context.Context, d *schema.ResourceData, mavenFeed *octopusdeploy.HelmFeed) error {
	d.Set("feed_uri", mavenFeed.FeedURI)
	d.Set("name", mavenFeed.Name)
	d.Set("space_id", mavenFeed.SpaceID)
	d.Set("username", mavenFeed.Username)

	if err := d.Set("package_acquisition_location_options", mavenFeed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(mavenFeed.GetID())

	return nil
}
