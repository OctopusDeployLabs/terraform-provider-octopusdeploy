package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandHelmFeed(d *schema.ResourceData) (*feeds.HelmFeed, error) {
	name := d.Get("name").(string)

	helmFeed, err := feeds.NewHelmFeed(name)
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
		helmFeed.Password = core.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("username"); ok {
		helmFeed.Username = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		helmFeed.SpaceID = v.(string)
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

func setHelmFeed(ctx context.Context, d *schema.ResourceData, helmFeed *feeds.HelmFeed) error {
	d.Set("feed_uri", helmFeed.FeedURI)
	d.Set("name", helmFeed.Name)
	d.Set("space_id", helmFeed.SpaceID)
	d.Set("username", helmFeed.Username)

	if err := d.Set("package_acquisition_location_options", helmFeed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(helmFeed.GetID())

	return nil
}
