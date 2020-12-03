package octopusdeploy

import "github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

func getAccountTypeQuery() *schema.Schema {
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

func getClonedFromProjectIDQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for cloned resources by a project ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getCommunicationStylesQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of communication styles. Valid communication styles are `AzureCloudService`, `AzureServiceFabricCluster`, `AzureWebApp`, `Ftp`, `Kubernetes`, `None`, `OfflineDrop`, `Ssh`, `TentacleActive`, or `TentaclePassive`.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
		ValidateDiagFunc: validateValueFunc([]string{
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
		}),
	}
}

func getDeploymentIDQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by deployment ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getEnvironmentsQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of environment IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getFilterQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter with which to search.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getIDDataSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A auto-generated identifier that includes the timestamp when this data source was last modified.",
		Type:        schema.TypeString,
	}
}

func getIDsQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getIsCloneQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search for cloned resources.",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getNameQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by name.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getPartialNameQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by the partial match of a name.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getRolesQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of role IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getSkipQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to specify the number of items to skip in the response.",
		Optional:    true,
		Type:        schema.TypeInt,
	}
}

func getTakeQuery() *schema.Schema {
	return &schema.Schema{
		Default:     1,
		Description: "A filter to specify the number of items to take (or return) in the response.",
		Type:        schema.TypeInt,
		Optional:    true,
	}
}

func getTenantsQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of tenant IDs.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getTenantTagsQuery() *schema.Schema {
	return &schema.Schema{
		Description: "A filter to search by a list of tenant tags.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}
