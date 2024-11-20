package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/serviceaccounts"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type serviceAccountOIDCIdentityDataSource struct {
	*Config
}

func NewServiceAccountOIDCIdentityDataSource() datasource.DataSource {
	return &serviceAccountOIDCIdentityDataSource{}
}

func (*serviceAccountOIDCIdentityDataSource) Metadata(_ context.Context, _ datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ServiceAccountOIDCIdentityDatasourceName)
}

func (s *serviceAccountOIDCIdentityDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	s.Config = DataSourceConfiguration(req, resp)
}

func (*serviceAccountOIDCIdentityDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.ServiceAccountOIDCIdentitySchema{}.GetDatasourceSchema()
}

func (s *serviceAccountOIDCIdentityDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.OIDCServiceAccountDatasourceSchemaModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	oidcIdentity, err := serviceaccounts.GetOIDCIdentityByID(s.Client, data.ServiceAccountID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load service account OIDC Identity", err.Error())
		return
	}

	updateServiceAccountOIDCDataModel(oidcIdentity, &data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func updateServiceAccountOIDCDataModel(request *serviceaccounts.OIDCIdentity, model *schemas.OIDCServiceAccountDatasourceSchemaModel) {
	model.Name = types.StringValue(request.Name)
	model.Issuer = types.StringValue(request.Issuer)
	model.Subject = types.StringValue(request.Subject)
	model.ID = types.StringValue(request.ID)
	model.ServiceAccountID = types.StringValue(request.ServiceAccountID)
}
