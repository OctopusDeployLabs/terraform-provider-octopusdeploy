package octopusdeploy

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func getAccountTypeSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description: "Specifies the type of the account. Valid account types are `AmazonWebServicesAccount`, `AmazonWebServicesRoleAccount`, `AzureServicePrincipal`, `AzureOIDC`, `AzureSubscription`, `None`, `SshKeyPair`, `Token`, or `UsernamePassword`.",
		ForceNew:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"AmazonWebServicesAccount",
			"AmazonWebServicesRoleAccount",
			"AzureServicePrincipal",
			"AzureOIDC",
			"AzureSubscription",
			"None",
			"SshKeyPair",
			"Token",
			"UsernamePassword",
		}, false)),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getAccessKeySchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description: "The access key associated with this resource.",
		Type:        schema.TypeString,
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getApplicationIDSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The application ID of this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getAuthenticationEndpointSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The authentication endpoint URI for this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getAzureEnvironmentSchema() *schema.Schema {
	return &schema.Schema{
		Computed: true,
		//Default:     "AzureCloud",
		Description: "The Azure environment associated with this resource. Valid Azure environments are `AzureCloud`, `AzureChinaCloud`, `AzureGermanCloud`, or `AzureUSGovernment`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"AzureCloud",
			"AzureChinaCloud",
			"AzureGermanCloud",
			"AzureUSGovernment",
		}, false)),
	}
}

func getCertificateDataFormatSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "Specifies the archive file format used for storing cryptography objects in the certificate. Valid formats are `Der`, `Pem`, `Pkcs12`, or `Unknown`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"Der",
			"Pem",
			"Pkcs12",
			"Unknown",
		}, false)),
	}
}

func getDescriptionSchema(resourceName string) *schema.Schema {
	return &schema.Schema{
		Description: fmt.Sprintf("The description of this %s.", resourceName),
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getDisplayNameSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description: "The display name of this resource.",
		Type:        schema.TypeString,
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getEmailAddressSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description: "The email address of this resource.",
		Type:        schema.TypeString,
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getEnvironmentsSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A list of environment IDs associated with this resource.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getHealthStatusSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "Represents the health status of this deployment target. Valid health statuses are `HasWarnings`, `Healthy`, `Unavailable`, `Unhealthy`, or `Unknown`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"HasWarnings",
			"Healthy",
			"Unavailable",
			"Unhealthy",
			"Unknown",
		}, false)),
	}
}

func getIDSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "The unique ID for this resource.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getIsSensitiveSchema() *schema.Schema {
	return &schema.Schema{
		Default:     false,
		Description: "Indicates whether or not this resource is considered sensitive and should be kept secret.",
		Optional:    true,
		Type:        schema.TypeBool,
	}
}

func getPasswordSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The password associated with this resource.",
		Sensitive:        true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getNameSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The name of this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getNameSchemaWithMaxLength(isRequired bool, maxLength int) *schema.Schema {
	schema := &schema.Schema{
		Description:      fmt.Sprintf("The name of this resource, no more than %d characters long", maxLength),
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringLenBetween(1, maxLength)),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getResourceManagerEndpointSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The resource manager endpoint URI for this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsURLWithHTTPS),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getSecretKeySchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description: "The secret key associated with this resource.",
		Sensitive:   true,
		Type:        schema.TypeString,
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getSortOrderSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "The sort order associated with this resource.",
		Optional:    true,
		Type:        schema.TypeInt,
	}
}

func getSpaceIDSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "The space ID associated with this resource.",
		Optional:    true,
		Type:        schema.TypeString,
		ForceNew:    true,
	}
}

func getStatusSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "The status of this resource. Valid statuses are `CalamariNeedsUpgrade`, `Disabled`, `NeedsUpgrade`, `Offline`, `Online`, or `Unknown`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"CalamariNeedsUpgrade",
			"Disabled",
			"NeedsUpgrade",
			"Offline",
			"Online",
			"Unknown",
		}, false)),
	}
}

func getStatusSummarySchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A summary elaborating on the status of this resource.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}

func getSubscriptionIDSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The subscription ID of this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getTenantedDeploymentSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "The tenanted deployment mode of the resource. Valid account types are `Untenanted`, `TenantedOrUntenanted`, or `Tenanted`.",
		Optional:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"Untenanted",
			"TenantedOrUntenanted",
			"Tenanted",
		}, false)),
	}
}

func getTenantIDSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The tenant ID of this resource.",
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.IsUUID),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getTenantsSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A list of tenant IDs associated with this resource.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getTenantTagsSchema() *schema.Schema {
	return &schema.Schema{
		Computed:    true,
		Description: "A list of tenant tags associated with this resource.",
		Elem:        &schema.Schema{Type: schema.TypeString},
		Optional:    true,
		Type:        schema.TypeList,
	}
}

func getTokenSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The token of this resource.",
		Sensitive:        true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getUsernameSchema(isRequired bool) *schema.Schema {
	schema := &schema.Schema{
		Description:      "The username associated with this resource.",
		Sensitive:        true,
		Type:             schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringIsNotEmpty),
	}

	if isRequired {
		schema.Required = true
	} else {
		schema.Optional = true
	}

	return schema
}

func getVariableTypeSchema() *schema.Schema {
	return &schema.Schema{
		Description: "The type of variable represented by this resource. Valid types are `AmazonWebServicesAccount`, `AzureAccount`, `GoogleCloudAccount`, `Certificate`, `Sensitive`, `String`, or `WorkerPool`.",
		Required:    true,
		Type:        schema.TypeString,
		ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{
			"AmazonWebServicesAccount",
			"AzureAccount",
			"GoogleCloudAccount",
			"Certificate",
			"Sensitive",
			"String",
			"WorkerPool",
		}, false)),
	}
}

func setDataSchema(schema *map[string]*schema.Schema) {
	for _, field := range *schema {
		field.Computed = true
		field.Default = nil
		field.DefaultFunc = nil
		field.AtLeastOneOf = nil
		field.ConflictsWith = nil
		field.ExactlyOneOf = nil
		field.MaxItems = 0
		field.MinItems = 0
		field.Optional = false
		field.Required = false
		field.ValidateDiagFunc = nil
	}
}

func getSubjectKeysSchema(description string) *schema.Schema {

	return &schema.Schema{
		Optional:    true,
		Description: description,
		Type:        schema.TypeList,
		Elem:        &schema.Schema{Type: schema.TypeString},
	}
}

func getOidcAudienceSchema() *schema.Schema {
	return &schema.Schema{
		Description: "Federated credentials audience, this value is used to establish a connection between external workload identities and Microsoft Entra ID.",
		Optional:    true,
		Type:        schema.TypeString,
	}
}
