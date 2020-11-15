package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func expandOfflineDrop(d *schema.ResourceData) *octopusdeploy.OfflineDropEndpoint {
	endpoint := octopusdeploy.NewOfflineDropEndpoint()
	endpoint.ID = d.Id()

	if v, ok := d.GetOk("applications_directory"); ok {
		endpoint.ApplicationsDirectory = v.(string)
	}

	if v, ok := d.GetOk("destination"); ok {
		endpoint.Destination = expandOfflineDropDestination(v)
	}

	if v, ok := d.GetOk("sensitive_variables_encryption_password"); ok {
		endpoint.SensitiveVariablesEncryptionPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	if v, ok := d.GetOk("working_directory"); ok {
		endpoint.SensitiveVariablesEncryptionPassword = octopusdeploy.NewSensitiveValue(v.(string))
	}

	return endpoint
}

func flattenOfflineDrop(endpoint *octopusdeploy.OfflineDropEndpoint) []interface{} {
	if endpoint == nil {
		return nil
	}

	rawEndpoint := map[string]interface{}{
		"applications_directory": endpoint.ApplicationsDirectory,
		"destination":            flattenOfflineDropDestination(endpoint.Destination),
		"id":                     endpoint.GetID(),
		"working_directory":      endpoint.WorkingDirectory,
	}

	return []interface{}{rawEndpoint}
}

func getOfflineDropSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"applications_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"destination": {
			Computed: true,
			Elem:     &schema.Resource{Schema: getOfflineDropDestinationSchema()},
			Optional: true,
			Type:     schema.TypeList,
		},
		"id": {
			Computed: true,
			Type:     schema.TypeString,
		},
		"working_directory": {
			Optional: true,
			Type:     schema.TypeString,
		},
		"sensitive_variables_encryption_password": {
			Optional:  true,
			Sensitive: true,
			Type:      schema.TypeString,
		},
	}
}
