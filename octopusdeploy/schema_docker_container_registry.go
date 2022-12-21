package octopusdeploy

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/feeds"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func expandDockerContainerRegistry(d *schema.ResourceData) (*feeds.DockerContainerRegistry, error) {
	name := d.Get("name").(string)

	feed, err := feeds.NewDockerContainerRegistry(name)
	if err != nil {
		return nil, err
	}

	feed.ID = d.Id()

	if v, ok := d.GetOk("api_version"); ok {
		feed.APIVersion = v.(string)
	}

	if v, ok := d.GetOk("feed_uri"); ok {
		feed.FeedURI = v.(string)
	}

	if v, ok := d.GetOk("registry_path"); ok {
		feed.RegistryPath = v.(string)
	}

	if v, ok := d.GetOk("space_id"); ok {
		feed.SpaceID = v.(string)
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

	return feed, nil
}

func getDockerContainerRegistrySchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"api_version": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"feed_uri": {
			Description:      "The URL to a Maven repository. This URL is the same value defined in the `repositories`/`repository`/`url` element of a Maven `settings.xml` file.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPorHTTPS),
		},
		"id": getIDSchema(),
		"name": {
			Description:      "A short, memorable, unique name for this feed. Example: ACME Builds.",
			Required:         true,
			Type:             schema.TypeString,
			ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
		},
		"password": getPasswordSchema(false),
		"package_acquisition_location_options": {
			Computed: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
			Optional: true,
			Type:     schema.TypeList,
		},
		"registry_path": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"space_id": getSpaceIDSchema(),
		"username": getUsernameSchema(false),
	}
}

func setDockerContainerRegistry(ctx context.Context, d *schema.ResourceData, feed *feeds.DockerContainerRegistry) error {
	d.Set("api_version", feed.APIVersion)
	d.Set("feed_uri", feed.FeedURI)
	d.Set("name", feed.Name)
	d.Set("registry_path", feed.RegistryPath)
	d.Set("space_id", feed.SpaceID)
	d.Set("username", feed.Username)

	if err := d.Set("package_acquisition_location_options", feed.PackageAcquisitionLocationOptions); err != nil {
		return fmt.Errorf("error setting package_acquisition_location_options: %s", err)
	}

	d.SetId(feed.GetID())

	return nil
}
