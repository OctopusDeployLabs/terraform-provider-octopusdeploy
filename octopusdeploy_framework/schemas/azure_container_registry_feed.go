package schemas

import (
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AzureContainerRegistryFeedSchema struct{}

var _ EntitySchema = AzureContainerRegistryFeedSchema{}

func (d AzureContainerRegistryFeedSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "This resource manages Azure Container Registry feed in Octopus Deploy (alias of Docker Container Registry feed)",
		Attributes: map[string]resourceSchema.Attribute{
			"api_version": resourceSchema.StringAttribute{
				Optional: true,
			},
			"feed_uri": GetFeedUriResourceSchema(),
			"id":       GetIdResourceSchema(),
			"name":     GetNameResourceSchema(true),
			"password": GetPasswordResourceSchema(false),
			"space_id": GetSpaceIdResourceSchema("Azure container registry feed"),
			"username": GetUsernameResourceSchema(false),
			"registry_path": resourceSchema.StringAttribute{
				Optional: true,
			},
			"oidc_authentication": resourceSchema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]resourceSchema.Attribute{
					"client_id": resourceSchema.StringAttribute{
						Description: "Unique identifier representing the application requesting authentication",
						Optional:    true,
						Computed:    true,
					},
					"tenant_id": resourceSchema.StringAttribute{
						Description: "Unique identifier representing the Azure AD instance hosting the authenticating application",
						Optional:    true,
						Computed:    true,
					},
					"audience": resourceSchema.StringAttribute{
						Description: "Audience representing the intended recipient of the OIDC token",
						Optional:    true,
						Computed:    true,
						Default:     stringdefault.StaticString(""),
					},
					"subject_keys": GetOidcSubjectKeysSchema("Keys to include in a deployment or runbook. Valid options are `space`, `feed`.", false),
				},
			},
		},
	}
}

func (d AzureContainerRegistryFeedSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{}
}

type AzureContainerRegistryFeedTypeResourceModel struct {
	APIVersion         types.String                                           `tfsdk:"api_version"`
	FeedUri            types.String                                           `tfsdk:"feed_uri"`
	Name               types.String                                           `tfsdk:"name"`
	Password           types.String                                           `tfsdk:"password"`
	SpaceID            types.String                                           `tfsdk:"space_id"`
	Username           types.String                                           `tfsdk:"username"`
	RegistryPath       types.String                                           `tfsdk:"registry_path"`
	OidcAuthentication *AzureContainerRegistryOidcAuthenticationResourceModel `tfsdk:"oidc_authentication"`

	ResourceModel
}

type AzureContainerRegistryOidcAuthenticationResourceModel struct {
	ClientId   types.String `tfsdk:"client_id"`
	TenantId   types.String `tfsdk:"tenant_id"`
	Audience   types.String `tfsdk:"audience"`
	SubjectKey types.List   `tfsdk:"subject_keys"`
}
