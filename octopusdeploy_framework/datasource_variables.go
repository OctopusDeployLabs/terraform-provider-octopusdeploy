package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type variablesDataSource struct {
	*Config
}

func NewVariablesDataSource() datasource.DataSource {
	return &variablesDataSource{}
}

func (*variablesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.VariablesDataSourceDescription)
}

func (*variablesDataSource) Schema(_ context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schemas.GetVariableDatasourceSchema()
}

func (e *variablesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	e.Config = DataSourceConfiguration(req, resp)
}

func (v *variablesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var err error
	var data schemas.VariablesDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	scope := schemas.MapToVariableScope(data.Scope)
	variables, err := variables.GetByName(v.Client, data.SpaceID.ValueString(), data.OwnerID.ValueString(), data.Name.ValueString(), &scope)
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("error reading variable with owner ID %s with name %s", data.OwnerID, data.Name), err.Error())
		return
	}

	if variables == nil {
		return
	}

	if len(variables) != 1 {
		resp.Diagnostics.AddError("could not find variable by name", fmt.Sprintf("expected to find 1 variable but got %d", len(variables)))
		return
	}

	variable := variables[0]

	data.Description = types.StringValue(variable.Description)
	data.IsEditable = types.BoolValue(variable.IsEditable)
	data.IsSensitive = types.BoolValue(variable.IsSensitive)
	data.Name = types.StringValue(variable.Name)
	data.Type = types.StringValue(variable.Type)
	if variable.IsSensitive {
		data.Value = types.StringNull()
	} else {
		data.Value = types.StringValue(variable.Value)
	}

	if variable.Prompt != nil {
		data.Prompt = types.ListValueMust(
			types.ObjectType{AttrTypes: schemas.VariablePromptOptionsObjectType()},
			[]attr.Value{schemas.MapFromVariablePromptOptions(variable.Prompt)},
		)
	}
	if !variable.Scope.IsEmpty() {
		data.Scope = types.ListValueMust(
			types.ObjectType{AttrTypes: schemas.VariableScopeObjectType()},
			[]attr.Value{schemas.MapFromVariableScope(variable.Scope)},
		)
	}

	data.ID = types.StringValue(variable.GetID())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
