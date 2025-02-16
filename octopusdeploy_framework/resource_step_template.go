package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type stepTemplateTypeResource struct {
	*Config
}

var _ resource.ResourceWithImportState = &stepTemplateTypeResource{}

func NewStepTemplateResource() resource.Resource {
	return &stepTemplateTypeResource{}
}

func (r *stepTemplateTypeResource) Metadata(_ context.Context, _ resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("step_template")
}

func (r *stepTemplateTypeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.StepTemplateSchema{}.GetResourceSchema()
}

func (r *stepTemplateTypeResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (*stepTemplateTypeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func (r *stepTemplateTypeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data schemas.StepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newActionTemplate, dg := mapStepTemplateResourceModelToActionTemplate(ctx, data)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionTemplate, err := actiontemplates.Add(r.Config.Client, newActionTemplate)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to create step template", err.Error())
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(&data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data schemas.StepTemplateTypeResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	actionTemplate, err := actiontemplates.GetByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, data, err, "action template"); err != nil {
			util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load environment", err.Error())
		}
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(&data, actionTemplate)...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.StepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	at, err := actiontemplates.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to load step template", err.Error())
		return
	}

	actionTemplateUpdate, dg := mapStepTemplateResourceModelToActionTemplate(ctx, data)
	resp.Diagnostics.Append(dg...)
	if resp.Diagnostics.HasError() {
		return
	}
	actionTemplateUpdate.ID = at.ID
	actionTemplateUpdate.SpaceID = at.SpaceID
	actionTemplateUpdate.Version = at.Version

	updatedActionTemplate, err := actiontemplates.Update(r.Config.Client, actionTemplateUpdate)
	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to update step template", err.Error())
		return
	}

	resp.Diagnostics.Append(mapStepTemplateToResourceModel(&data, updatedActionTemplate)...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.StepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := actiontemplates.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		util.AddDiagnosticError(resp.Diagnostics, r.Config.SystemInfo, "unable to delete step template", err.Error())
		return
	}
}

func mapStepTemplateToResourceModel(data *schemas.StepTemplateTypeResourceModel, at *actiontemplates.ActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	data.ID = types.StringValue(at.ID)
	data.SpaceID = types.StringValue(at.SpaceID)
	data.Name = types.StringValue(at.Name)
	data.Version = types.Int32Value(at.Version)
	data.Description = types.StringValue(at.Description)
	data.CommunityActionTemplateId = types.StringValue(at.CommunityActionTemplateID)
	data.ActionType = types.StringValue(at.ActionType)

	// Parameters
	sParams, dg := convertStepTemplateToParameterAttributes(at.Parameters)
	resp.Append(dg...)
	data.Parameters = sParams

	// Properties
	stringProps := make(map[string]attr.Value, len(at.Properties))
	for keys, value := range at.Properties {
		stringProps[keys] = types.StringValue(value.Value)
	}
	props, dg := types.MapValue(types.StringType, stringProps)
	resp.Append(dg...)
	data.Properties = props

	// Packages
	pkgs, dg := convertStepTemplateToPackageAttributes(at.Packages)
	resp.Append(dg...)
	data.Packages = pkgs

	return resp
}

func mapStepTemplateResourceModelToActionTemplate(ctx context.Context, data schemas.StepTemplateTypeResourceModel) (*actiontemplates.ActionTemplate, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	at := actiontemplates.NewActionTemplate(data.Name.ValueString(), data.ActionType.ValueString())

	at.SpaceID = data.SpaceID.ValueString()
	at.Description = data.Description.ValueString()
	if !data.CommunityActionTemplateId.IsNull() {
		at.CommunityActionTemplateID = data.CommunityActionTemplateId.ValueString()
	}

	pkgs := make([]schemas.StepTemplatePackageType, 0, len(data.Packages.Elements()))
	resp.Append(data.Packages.ElementsAs(ctx, &pkgs, false)...)
	if resp.HasError() {
		return at, resp
	}

	props := make(map[string]types.String, len(data.Properties.Elements()))
	resp.Append(data.Properties.ElementsAs(ctx, &props, false)...)
	if resp.HasError() {
		return at, resp
	}

	params := make([]schemas.StepTemplateParameterType, 0, len(data.Parameters.Elements()))
	resp.Append(data.Parameters.ElementsAs(ctx, &params, false)...)
	if resp.HasError() {
		return at, resp
	}

	if len(props) > 0 {
		templateProps := make(map[string]core.PropertyValue, len(props))
		for key, val := range props {
			templateProps[key] = core.NewPropertyValue(val.ValueString(), false)
		}
		at.Properties = templateProps
	} else {
		at.Properties = make(map[string]core.PropertyValue)
	}

	at.Packages = make([]packages.PackageReference, len(pkgs))
	if len(pkgs) > 0 {
		for i, val := range pkgs {
			pkgProps := convertAttributeStepTemplatePackageProperty(val.Properties.Attributes())
			pkgRef := packages.PackageReference{
				AcquisitionLocation: val.AcquisitionLocation.ValueString(),
				FeedID:              val.FeedID.ValueString(),
				Properties:          pkgProps,
				Name:                val.Name.ValueString(),
				PackageID:           val.PackageID.ValueString(),
			}
			pkgRef.ID = val.ID.ValueString()
			at.Packages[i] = pkgRef
		}
	}

	at.Parameters = make([]actiontemplates.ActionTemplateParameter, len(params))
	if len(params) > 0 {
		paramIDMap := make(map[string]bool, len(params))
		for i, val := range params {
			defaultValue := core.NewPropertyValue(val.DefaultValue.ValueString(), false)
			at.Parameters[i] = actiontemplates.ActionTemplateParameter{
				DefaultValue:    &defaultValue,
				Name:            val.Name.ValueString(),
				Label:           val.Label.ValueString(),
				HelpText:        val.HelpText.ValueString(),
				DisplaySettings: util.ConvertAttrStringMapToStringMap(val.DisplaySettings.Elements()),
			}
			id := val.ID.ValueString()
			if _, ok := paramIDMap[id]; ok {
				resp.AddError("ID conflict", fmt.Sprintf("conflicting UUID's within parameters list: %s", id))
			}
			paramIDMap[val.ID.ValueString()] = true
			at.Parameters[i].ID = id
			at.Parameters[i].ID = val.ID.ValueString()
		}
	}
	if resp.HasError() {
		return at, resp
	}
	return at, resp
}

func convertStepTemplateToPackageAttributes(atPackage []packages.PackageReference) (types.List, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	pkgs := make([]attr.Value, len(atPackage))
	for key, val := range atPackage {
		mapVal, dg := convertStepTemplatePackageAttribute(val)
		resp.Append(dg...)
		if resp.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}), resp
		}
		pkgs[key] = mapVal
	}
	pkgSet, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}, pkgs)
	resp.Append(dg...)
	if resp.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackageTypeAttributes()}), resp
	}
	return pkgSet, dg
}

func convertStepTemplateToParameterAttributes(atParams []actiontemplates.ActionTemplateParameter) (types.List, diag.Diagnostics) {
	resp := diag.Diagnostics{}
	params := make([]attr.Value, len(atParams))
	for i, val := range atParams {
		objVal, dg := convertStepTemplateParameterAttribute(val)
		resp.Append(dg...)
		if resp.HasError() {
			return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), resp
		}
		params[i] = objVal
	}
	sParams, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}, params)
	resp.Append(dg...)
	if resp.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), resp
	}
	return sParams, resp
}

func convertStepTemplateParameterAttribute(atp actiontemplates.ActionTemplateParameter) (types.Object, diag.Diagnostics) {
	displaySettings, dg := types.MapValue(types.StringType, util.ConvertStringMapToAttrStringMap(atp.DisplaySettings))
	if dg.HasError() {
		return types.ObjectNull(schemas.GetStepTemplateParameterTypeAttributes()), dg
	}
	return types.ObjectValue(schemas.GetStepTemplateParameterTypeAttributes(), map[string]attr.Value{
		"id":               types.StringValue(atp.ID),
		"name":             types.StringValue(atp.Name),
		"label":            types.StringValue(atp.Label),
		"help_text":        types.StringValue(atp.HelpText),
		"default_value":    types.StringValue(atp.DefaultValue.Value),
		"display_settings": displaySettings,
	})
}

func convertStepTemplatePackageAttribute(atp packages.PackageReference) (types.Object, diag.Diagnostics) {
	props, dg := convertStepTemplatePackagePropertyAttribute(atp.Properties)
	if dg.HasError() {
		return types.ObjectNull(schemas.GetStepTemplatePackageTypeAttributes()), dg
	}
	return types.ObjectValue(schemas.GetStepTemplatePackageTypeAttributes(), map[string]attr.Value{
		"id":                   types.StringValue(atp.ID),
		"acquisition_location": types.StringValue(atp.AcquisitionLocation),
		"name":                 types.StringValue(atp.Name),
		"feed_id":              types.StringValue(atp.FeedID),
		"package_id":           types.StringValue(atp.PackageID),
		"properties":           props,
	})
}

func convertStepTemplatePackagePropertyAttribute(atpp map[string]string) (types.Object, diag.Diagnostics) {
	prop := make(map[string]attr.Value)
	diags := diag.Diagnostics{}

	// We need to manually convert the string map to ensure all fields are set.
	if extract, ok := atpp["Extract"]; ok {
		prop["extract"] = types.StringValue(extract)
	} else {
		diags.AddWarning("Package property missing value.", "extract value missing from package property")
		prop["extract"] = types.StringNull()
	}

	if purpose, ok := atpp["Purpose"]; ok {
		prop["purpose"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "purpose value missing from package property")
		prop["purpose"] = types.StringNull()
	}

	if purpose, ok := atpp["PackageParameterName"]; ok {
		prop["package_parameter_name"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "package_parameter_name value missing from package property")
		prop["package_parameter_name"] = types.StringNull()
	}

	if selectionMode, ok := atpp["SelectionMode"]; ok {
		prop["selection_mode"] = types.StringValue(selectionMode)
	} else {
		diags.AddWarning("Package property missing value.", "selection_mode value missing from package property")
		prop["selection_mode"] = types.StringNull()
	}

	propMap, dg := types.ObjectValue(schemas.GetStepTemplatePackagePropertiesTypeAttributes(), prop)
	if dg.HasError() {
		diags.Append(dg...)
		return types.ObjectNull(schemas.GetStepTemplatePackagePropertiesTypeAttributes()), diags
	}
	return propMap, diags
}

func convertAttributeStepTemplatePackageProperty(prop map[string]attr.Value) map[string]string {
	atpp := make(map[string]string)

	if extract, ok := prop["extract"]; ok {
		atpp["Extract"] = extract.(types.String).ValueString()
	}

	if purpose, ok := prop["purpose"]; ok {
		atpp["Purpose"] = purpose.(types.String).ValueString()
	}

	if purpose, ok := prop["package_parameter_name"]; ok {
		atpp["PackageParameterName"] = purpose.(types.String).ValueString()
	}

	if selectionMode, ok := prop["selection_mode"]; ok {
		atpp["SelectionMode"] = selectionMode.(types.String).ValueString()
	}
	return atpp
}
