package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type variableTypeResource struct {
	*Config
}

var _ resource.ResourceWithImportState = &variableTypeResource{}

func NewVariableResource() resource.Resource {
	return &variableTypeResource{}
}

func (r *variableTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Info(ctx, fmt.Sprintf("importing variable (%s)", req.ID))

	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 {
		resp.Diagnostics.AddError(
			"unexpected import identifier",
			fmt.Sprintf("%s_variable import must be in the form of OwnerID:VariableID (e.g. Projects-62:0906031f-68ba-4a15-afaa-657c1564e07b)", util.GetProviderName()))
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("owner_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[1])...)
}

func (r *variableTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetProviderName() + "_variable"
}

func (r *variableTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetVariableResourceSchema()
}

func (r *variableTypeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *variableTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	variableOwnerId, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	name := data.Name.ValueString()
	newVariable := variables.NewVariable(name)
	newVariable.Description = data.Description.ValueString()
	newVariable.IsEditable = data.IsEditable.ValueBool()
	newVariable.IsSensitive = data.IsSensitive.ValueBool()
	newVariable.Type = data.Type.ValueString()
	newVariable.Scope = schemas.MapToVariableScope(data.Scope)
	newVariable.Prompt = schemas.MapToVariablePrompOptions(data.Prompt)
	newVariable.SpaceID = data.SpaceID.ValueString()

	if newVariable.IsSensitive {
		newVariable.Type = schemas.VariableTypeNames.Sensitive
		newVariable.Value = data.SensitiveValue.ValueString()
	} else {
		newVariable.Value = data.Value.ValueString()
	}

	tflog.Info(ctx, fmt.Sprintf("creating variable: %#v", newVariable))

	variableSet, err := variables.AddSingle(r.Config.Client, r.Config.SpaceID, variableOwnerId.ValueString(), newVariable)
	if err != nil {
		resp.Diagnostics.AddError("create variable failed", err.Error())
		return
	}

	err = validateVariable(&variableSet, newVariable, variableOwnerId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("create variable failed", err.Error())
		return
	}

	mapVariableToState(&data, newVariable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("reading variable (%s)", data.ID))

	variableOwnerID, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	variable, err := variables.GetByID(r.Config.Client, data.SpaceID.ValueString(), variableOwnerID.ValueString(), data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load variable", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable read (%s)", data.ID))
	mapVariableToState(&data, variable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data, state schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("updating variable (%s)", data.ID))

	variableOwnerId, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	name := data.Name.ValueString()
	updatedVariable := variables.NewVariable(name)
	updatedVariable.Description = data.Description.ValueString()
	updatedVariable.IsEditable = data.IsEditable.ValueBool()
	updatedVariable.IsSensitive = data.IsSensitive.ValueBool()
	updatedVariable.Type = data.Type.ValueString()
	updatedVariable.Scope = schemas.MapToVariableScope(data.Scope)
	updatedVariable.Prompt = schemas.MapToVariablePrompOptions(data.Prompt)
	updatedVariable.SpaceID = state.SpaceID.ValueString()

	if updatedVariable.IsSensitive {
		updatedVariable.Type = schemas.VariableTypeNames.Sensitive
		updatedVariable.Value = data.SensitiveValue.ValueString()
	} else {
		updatedVariable.Value = data.Value.ValueString()
	}

	updatedVariable.ID = state.ID.ValueString()

	variableSet, err := variables.UpdateSingle(r.Config.Client, state.SpaceID.ValueString(), variableOwnerId.ValueString(), updatedVariable)
	if err != nil {
		resp.Diagnostics.AddError("update variable failed", err.Error())
		return
	}

	err = validateVariable(&variableSet, updatedVariable, variableOwnerId.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("update variable failed", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable updated (%s)", data.ID))

	mapVariableToState(&data, updatedVariable)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *variableTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	internal.Mutex.Lock()
	defer internal.Mutex.Unlock()

	var data schemas.VariableTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("deleting variable (%s)", data.ID.ValueString()))
	variableOwnerID, err := getVariableOwnerID(&data)
	if err != nil {
		resp.Diagnostics.AddError("invalid resource configuration", err.Error())
		return
	}

	if _, err := variables.DeleteSingle(r.Config.Client, r.Config.SpaceID, variableOwnerID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete variable", err.Error())
		return
	}

	tflog.Info(ctx, fmt.Sprintf("variable deleted (%s)", data.ID))
}

func (r *variableTypeResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data schemas.VariableTypeResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	isSensitive := data.IsSensitive
	variableType := data.Type

	if isSensitive.ValueBool() && variableType.ValueString() != schemas.VariableTypeNames.Sensitive {
		resp.Diagnostics.AddError(
			"invalid resource configuration",
			fmt.Sprintf("when %s is set to true, type needs to be '%s'", schemas.VariableSchemaAttributeNames.IsSensitive, schemas.VariableTypeNames.Sensitive),
		)
		return
	}

	if !isSensitive.ValueBool() && variableType.ValueString() == schemas.VariableTypeNames.Sensitive {
		resp.Diagnostics.AddError(
			"invalid resource configuration",
			fmt.Sprintf("when type is set to '%s', %s needs to be true", schemas.VariableSchemaAttributeNames.IsSensitive, schemas.VariableTypeNames.Sensitive),
		)
	}
}

func getVariableOwnerID(data *schemas.VariableTypeResourceModel) (*basetypes.StringValue, error) {
	if data.ProjectID.IsNull() && data.OwnerID.IsNull() {
		return nil, fmt.Errorf("one of %s or %s must be configured", schemas.VariableSchemaAttributeNames.ProjectID, schemas.VariableSchemaAttributeNames.OwnerID)
	} else if !data.ProjectID.IsNull() {
		return &data.ProjectID, nil
	} else {
		return &data.OwnerID, nil
	}
}

func validateVariable(variableSet *variables.VariableSet, newVariable *variables.Variable, variableOwnerId string) error {
	for _, v := range variableSet.Variables {
		if v.Name == newVariable.Name && v.Type == newVariable.Type && (v.IsSensitive || v.Value == newVariable.Value) && v.Description == newVariable.Description && v.IsSensitive == newVariable.IsSensitive {
			scopeMatches, _, err := variables.MatchesScope(v.Scope, &newVariable.Scope)
			if err != nil || !scopeMatches {
				return err
			}
			if scopeMatches {
				newVariable.ID = v.GetID()
				return nil
			}
		}
	}

	return fmt.Errorf("unable to locate variable for owner ID %s", variableOwnerId)
}

func mapVariableToState(data *schemas.VariableTypeResourceModel, variable *variables.Variable) {
	data.SpaceID = types.StringValue(variable.SpaceID)
	data.Name = types.StringValue(variable.Name)
	data.Description = types.StringValue(variable.Description)
	if !data.IsEditable.IsNull() {
		data.IsEditable = types.BoolValue(variable.IsEditable)
	}
	if !data.IsSensitive.IsNull() {
		data.IsSensitive = types.BoolValue(variable.IsSensitive)
	}
	data.Type = types.StringValue(variable.Type)

	if variable.IsSensitive {
		data.Value = types.StringNull()
	} else {
		if !data.Value.IsNull() {
			data.Value = types.StringValue(variable.Value)
		}
	}

	if !data.Prompt.IsNull() {
		data.Prompt = types.ListValueMust(
			types.ObjectType{AttrTypes: schemas.VariablePromptOptionsObjectType()},
			[]attr.Value{schemas.MapFromVariablePromptOptions(variable.Prompt)},
		)
	}

	if !data.Scope.IsNull() {
		data.Scope = types.ListValueMust(
			types.ObjectType{AttrTypes: schemas.VariableScopeObjectType()},
			[]attr.Value{schemas.MapFromVariableScope(variable.Scope)},
		)
	}

	data.EncryptedValue = types.StringNull()
	data.KeyFingerprint = types.StringNull()

	data.ID = types.StringValue(variable.GetID())
}
