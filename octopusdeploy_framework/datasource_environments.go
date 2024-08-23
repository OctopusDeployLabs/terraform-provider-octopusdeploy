package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type environmentDataSource struct {
	*Config
}

type environmentsDataSourceModel struct {
	ID           types.String `tfsdk:"id"`
	SpaceID      types.String `tfsdk:"space_id"`
	IDs          types.List   `tfsdk:"ids"`
	Name         types.String `tfsdk:"name"`
	PartialName  types.String `tfsdk:"partial_name"`
	Skip         types.Int64  `tfsdk:"skip"`
	Take         types.Int64  `tfsdk:"take"`
	Environments types.List   `tfsdk:"environments"`
}

func NewEnvironmentsDataSource() datasource.DataSource {
	return &environmentDataSource{}
}

func (*environmentDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("environments")
}

func (*environmentDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Provides information about existing environments.",
		Attributes: map[string]schema.Attribute{
			//request
			"ids":          schemas.GetQueryIDsDatasourceSchema(),
			"space_id":     schemas.GetSpaceIdDatasourceSchema(schemas.EnvironmentResourceDescription, false),
			"name":         schemas.GetQueryNameDatasourceSchema(),
			"partial_name": schemas.GetQueryPartialNameDatasourceSchema(),
			"skip":         schemas.GetQuerySkipDatasourceSchema(),
			"take":         schemas.GetQueryTakeDatasourceSchema(),

			//response
			"id": schemas.GetIdDatasourceSchema(true),
			"environments": schema.ListNestedAttribute{
				Computed: true,
				Optional: false,
				NestedObject: schema.NestedAttributeObject{
					Attributes: schemas.GetEnvironmentDatasourceSchema(),
				},
			},
		},
	}
}

func (e *environmentDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (e *environmentDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data environmentsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	query := environments.EnvironmentsQuery{
		IDs:         util.GetIds(data.IDs),
		PartialName: data.PartialName.ValueString(),
		Name:        data.Name.ValueString(),
		Skip:        util.GetNumber(data.Skip),
		Take:        util.GetNumber(data.Take),
	}

	util.DatasourceReading(ctx, "environments", query)

	existingEnvironments, err := environments.Get(e.Client, data.SpaceID.ValueString(), query)
	if err != nil {
		resp.Diagnostics.AddError("unable to load environments", err.Error())
		return
	}

	var mappedEnvironments []schemas.EnvironmentTypeResourceModel
	if data.Name.IsNull() {
		tflog.Debug(ctx, fmt.Sprintf("environments returned from API: %#v", existingEnvironments))
		for _, environment := range existingEnvironments.Items {
			mappedEnvironments = append(mappedEnvironments, schemas.MapFromEnvironment(ctx, environment))
		}
	} else { // if name has been specified, match by exact name rather than partial name as the API does
		var matchedEnvironment *environments.Environment
		tflog.Debug(ctx, fmt.Sprintf("matching environment by name: %s", data.Name))
		for _, env := range existingEnvironments.Items {
			if strings.EqualFold(env.Name, data.Name.ValueString()) {
				matchedEnvironment = env
			}
		}
		if matchedEnvironment != nil {
			mappedEnvironments = append(mappedEnvironments, schemas.MapFromEnvironment(ctx, matchedEnvironment))
		}
	}

	util.DatasourceResultCount(ctx, "environments", len(mappedEnvironments))

	data.Environments, _ = types.ListValueFrom(ctx, types.ObjectType{AttrTypes: schemas.EnvironmentObjectType()}, mappedEnvironments)
	data.ID = types.StringValue("Environments " + time.Now().UTC().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
