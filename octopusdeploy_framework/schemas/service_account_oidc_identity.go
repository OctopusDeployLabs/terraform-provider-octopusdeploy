package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const ServiceAccountOIDCIdentityResourceName = "service_account_oidc_identity"
const ServiceAccountOIDCIdentityDatasourceName = "service_account_oidc_identity"

type ServiceAccountOIDCIdentitySchema struct{}

var _ EntitySchema = ServiceAccountOIDCIdentitySchema{}

func (d ServiceAccountOIDCIdentitySchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			"id":   GetIdResourceSchema(),
			"name": GetNameResourceSchema(true),
			"service_account_id": util.ResourceString().
				Description("ID of the user to associate this identity to").
				Required().
				PlanModifiers(stringplanmodifier.RequiresReplace()).
				Build(),
			"issuer": util.ResourceString().
				Description("OIDC issuer url").
				Required().
				Build(),
			"subject": util.ResourceString().
				Description("OIDC subject claims").
				Required().
				Build(),
		},
		Description: "This resource manages manages OIDC service account for the associated user",
	}
}

func (d ServiceAccountOIDCIdentitySchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Attributes: map[string]datasourceSchema.Attribute{
			"id": GetIdDatasourceSchema(false),
			"service_account_id": util.DataSourceString().
				Description("ID of the user associated to this identity").
				Required().
				Build(),
			// Response
			"name": util.DataSourceString().
				Description("Name of the user associated to this identity").
				Computed().
				Build(),
			"issuer": util.DataSourceString().
				Description("OIDC issuer url").
				Computed().
				Build(),
			"subject": util.DataSourceString().
				Description("OIDC subject claims").
				Computed().
				Build(),
		},
	}
}

type OIDCServiceAccountSchemaModel struct {
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Issuer           types.String `tfsdk:"issuer"`
	Subject          types.String `tfsdk:"subject"`

	ResourceModel
}

type OIDCServiceAccountDatasourceSchemaModel struct {
	ID               types.String `tfsdk:"id"`
	ServiceAccountID types.String `tfsdk:"service_account_id"`
	Name             types.String `tfsdk:"name"`
	Issuer           types.String `tfsdk:"issuer"`
	Subject          types.String `tfsdk:"subject"`
}
