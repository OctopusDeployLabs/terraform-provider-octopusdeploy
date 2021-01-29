package octopusdeploy

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getQueryAccountType() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of account types.  Valid account types are `AmazonWebServicesAccount`, `AmazonWebServicesRoleAccount`, `AzureServicePrincipal`, `AzureSubscription`, `None`, `SshKeyPair`, `Token`, or `UsernamePassword`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validateValueFunc([]string{
			"AmazonWebServicesAccount",
			"AmazonWebServicesRoleAccount",
			"AzureServicePrincipal",
			"AzureSubscription",
			"None",
			"SshKeyPair",
			"Token",
			"UsernamePassword",
		}),
	}
}

func getQueryArchived() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for resources that have been archived.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryClonedFromProjectID() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for cloned resources by a project ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryClonedFromTenantID() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for a cloned tenant by its ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryCommunicationStyles() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of communication styles. Valid communication styles are `AzureCloudService`, `AzureServiceFabricCluster`, `AzureWebApp`, `Ftp`, `Kubernetes`, `None`, `OfflineDrop`, `Ssh`, `TentacleActive`, or `TentaclePassive`.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"AzureCloudService",
				"AzureServiceFabricCluster",
				"AzureWebApp",
				"Ftp",
				"Kubernetes",
				"None",
				"OfflineDrop",
				"Ssh",
				"TentacleActive",
				"TentaclePassive",
			}, false)),
		},
		Optional: true,
		Type:     schema.TypeList,
	}
}

func getQueryDeploymentID() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by deployment ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryEnvironments() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of environment IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryFeedType() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by feed type. Valid feed types are `AwsElasticContainerRegistry`, `BuiltIn`, `Docker`, `GitHub`, `Helm`, `Maven`, `NuGet`, or `OctopusProject`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
			"AwsElasticContainerRegistry",
			"BuiltIn",
			"Docker",
			"GitHub",
			"Helm",
			"Maven",
			"NuGet",
			"OctopusProject",
		}, false)),
	}
}

func getQueryFilter() *schema.Schema {
	return &schema.Schema{
		Description: "A filter with which to search.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryFirstResult() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to define the first result.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryHealthStatuses() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of health statuses of resources. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.",
		Elem: &schema.Schema{
			Type: schema.TypeString,
			ValidateDiagFunc: validateDiagFunc(validation.StringInSlice([]string{
				"HasWarnings",
				"Healthy",
				"Unavailable",
				"Unhealthy",
				"Unknown",
			}, false)),
		},
		Optional: true,
		Type:     schema.TypeList,
	}
}

func getDataSchemaID() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A auto-generated identifier that includes the timestamp when this data source was last modified.",
		Type:        schema.TypeString,
	}
}

func getQueryShellNames() *schema.Schema {
	return &schema.Schema{
		Description: "A list of shell names to match in the query and/or search",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryIDs() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryIncludeSystem() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to include system teams.",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getQueryIsClone() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for cloned resources.",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getQueryIsDisabled() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by the disabled status of a resource.",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getQueryName() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by name.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryOrderBy() *schema.Schema {
	return &schema.Schema{
		Description: "A filter used to order the search results.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryPartialName() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by the partial match of a name.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryProjectID() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a project ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryRoles() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of role IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQuerySearch() *schema.Schema {
	return &schema.Schema{
		Description: "A filter of terms used the search operation.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQuerySkip() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to specify the number of items to skip in the response.",
		Optional:    true,
		Type:        schema.TypeInt,
	}
}

func getQueryTags() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of tags.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryTake() *schema.Schema {
	return &schema.Schema{
		Default:     1,
		Description: "A filter to specify the number of items to take (or return) in the response.",
		Type:        schema.TypeInt,
		Optional:    true,
	}
}

func getQueryTenant() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a tenant ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getQueryTenants() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of tenant IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryTenantTags() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of tenant tags.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getQueryThumbprint() *schema.Schema {
	return &schema.Schema{
		Description: "The thumbprint of the deployment target to match in the query and/or search",
		Optional:    true,
		Type:        schema.TypeString,
	}
}
