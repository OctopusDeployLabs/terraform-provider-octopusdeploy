package octopusdeploy

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandArtifactoryGenericFeed(d *schema.ResourceData) (*feeds.ArtifactoryGenericFeed, error) {
	name := d.Get("name").(string)

	feed, err := feeds.NewArtifactoryGenericFeed(name)
	if err != nil {
		return nil, err
	}

	feed.ID = d.Id()

	if v, ok := d.GetOk("feed_uri"); ok {
		feed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("package_acquisition_location_options"); ok {
		feed.PackageAcquisitionLocationOptions = getSliceFromTerraformTypeList(v)
	}

	if v, ok := d.GetOk("password"); ok {
		feed.Password = core.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("space_id"); ok {
		feed.SpaceID = v.(string)
	}

	if v, ok := d.GetOk("username"); ok {
		feed.Username = v.(string)
	}

	if v, ok := d.GetOk("layout_regex"); ok {
		feed.LayoutRegex = v.(string)
	}

	if v, ok := d.GetOk("repository"); ok {
		feed.Repository = v.(string)
	}

	return feed, nil
}

func getArtifactoryGenericFeedSchema() map[string]*schema.Schema {
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
		"repository": {
			Computed: false,
			Required: true,
			Type:     schema.TypeString,
		},
		"layout_regex": {
			Computed: false,
			Required: false,
			Optional: true,
			Type:     schema.TypeString,
		},
	}
}

func setArtifactoryGenericFeed(ctx context.Context, d *schema.ResourceData, feed *feeds.ArtifactoryGenericFeed) error {
	d.Set("feed_uri", feed.FeedURI)
	d.Set("name", feed.Name)
	d.Set("space_id", feed.SpaceID)
	d.Set("username", feed.Username)
	d.Set("repository", feed.Repository)
	d.Set("layout_regex", feed.LayoutRegex)

	if err := d.Set("package_acquisition_location_options", feed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(feed.GetID())

	return nil
}
