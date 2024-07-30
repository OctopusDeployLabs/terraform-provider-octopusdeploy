package octopusdeploy_framework

import (
	"context"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"
)

type tenantsDataSource struct {
	*Config
}

func NewTenantsDataSource() datasource.DataSource {
	return &tenantsDataSource{}
}

func (*tenantsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("tenants")
}

func (e *tenantsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (*tenantsDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = datasourceSchema.Schema{
		Description: "Provides information about existing tenants.",
		Attributes:  schemas.GetTenantsDataSourceSchema(),
		Blocks: map[string]datasourceSchema.Block{
			"tenants": datasourceSchema.ListNestedBlock{
				Description: "A list of tenants that match the filter(s).",
				NestedObject: datasourceSchema.NestedBlockObject{
					Attributes: schemas.GetTenantDataSourceSchema(),
				},
			},
		},
	}
}

func (b *tenantsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.TenantsModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := tenants.TenantsQuery{
		ClonedFromTenantID: data.ClonedFromTenantId.ValueString(),
		IDs:                util.GetStringSlice(data.IDs),
		IsClone:            data.IsClone.ValueBool(),
		Name:               data.Name.ValueString(),
		PartialName:        data.PartialName.ValueString(),
		ProjectID:          data.ProjectId.ValueString(),
		Skip:               int(data.Skip.ValueInt64()),
		Tags:               util.GetStringSlice(data.Tags),
		Take:               int(data.Take.ValueInt64()),
	}

	existingTenants, err := tenants.Get(b.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load tenants", err.Error())
		return
	}

	flattenedTenants := []interface{}{}
	for _, tenant := range existingTenants.Items {
		flattenedTenants = append(flattenedTenants, schemas.FlattenTenant(tenant))
	}

	data.ID = types.StringValue("Tenants " + time.Now().UTC().String())
	data.Tenants, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.TenantObjectType()}, flattenedTenants)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
