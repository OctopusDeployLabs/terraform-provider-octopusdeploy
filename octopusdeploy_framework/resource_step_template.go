package octopusdeploy_framework

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/actiontemplates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/environments"
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

func (r *stepTemplateTypeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("step_template")
}

func (r *stepTemplateTypeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
	pkgs := make([]schemas.StepTemplatePackageType, 0, len(data.Packages.Elements()))
	resp.Diagnostics.Append(data.Packages.ElementsAs(ctx, &pkgs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := make(map[string]types.String, len(data.Properties.Elements()))
	resp.Diagnostics.Append(data.Properties.ElementsAs(ctx, &props, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := make([]schemas.StepTemplateParameterType, 0, len(data.Parameters.Elements()))
	resp.Diagnostics.Append(data.Parameters.ElementsAs(ctx, &params, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	newActionTemplate := actiontemplates.NewActionTemplate(data.Name.ValueString(), data.ActionType.ValueString())
	newActionTemplate.SpaceID = data.SpaceID.ValueString()
	newActionTemplate.Description = data.Description.ValueString()
	if !data.CommunityActionTemplateId.IsNull() {
		newActionTemplate.CommunityActionTemplateID = data.CommunityActionTemplateId.ValueString()
	}

	if len(props) > 0 {
		templateProps := make(map[string]core.PropertyValue, len(props))
		for key, val := range props {
			templateProps[key] = core.NewPropertyValue(val.ValueString(), false)
		}
		newActionTemplate.Properties = templateProps
	} else {
		newActionTemplate.Properties = make(map[string]core.PropertyValue, 0)
	}

	newActionTemplate.Packages = make([]packages.PackageReference, len(pkgs))
	if len(pkgs) > 0 {
		for i, val := range pkgs {
			//		pkgProps := make(map[string]types.String, len(val.Properties.Attributes()))
			// TODO: fix
			// resp.Diagnostics.Append(val.Properties.(ctx, &pkgProps, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			pkgRef := packages.PackageReference{
				AcquisitionLocation: val.AcquisitionLocation.ValueString(),
				FeedID:              val.FeedID.ValueString(),
				// 				Properties:          util.ConvertAttrStringMapToStringMap(pkgProps),
				Name:      val.Name.ValueString(),
				PackageID: val.PackageID.ValueString(),
			}
			newActionTemplate.Packages[i] = pkgRef
		}
	}

	newActionTemplate.Parameters = make([]actiontemplates.ActionTemplateParameter, len(params))
	if len(params) > 0 {
		for i, val := range params {
			defaultValue := core.NewPropertyValue(val.DefaultValue.ValueString(), false)
			newActionTemplate.Parameters[i] = actiontemplates.ActionTemplateParameter{
				DefaultValue:    &defaultValue,
				Name:            val.Name.ValueString(),
				Label:           val.Label.ValueString(),
				HelpText:        val.HelpText.ValueString(),
				DisplaySettings: util.ConvertAttrStringMapToStringMap(val.DisplaySettings.Elements()),
			}
		}
	}

	actionTemplate, err := actiontemplates.Add(r.Config.Client, newActionTemplate)
	if err != nil {
		resp.Diagnostics.AddError("unable to create environment", err.Error())
		return
	}

	updateStepTemplate(ctx, &data, actionTemplate)
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
			resp.Diagnostics.AddError("unable to load environment", err.Error())
		}
		return
	}

	updateStepTemplate(ctx, &data, actionTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state schemas.StepTemplateTypeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	pkgs := make([]schemas.StepTemplatePackageType, 0, len(data.Packages.Elements()))
	resp.Diagnostics.Append(data.Packages.ElementsAs(ctx, &pkgs, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	props := make(map[string]types.String, len(data.Properties.Elements()))
	resp.Diagnostics.Append(data.Properties.ElementsAs(ctx, &props, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	params := make([]schemas.StepTemplateParameterType, 0, len(data.Parameters.Elements()))
	resp.Diagnostics.Append(data.Parameters.ElementsAs(ctx, &props, false)...)
	if resp.Diagnostics.HasError() {
		return
	}

	at, err := actiontemplates.GetByID(r.Config.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("unable to load step template", err.Error())
		return
	}

	actionTemplateUpdate := actiontemplates.NewActionTemplate(data.Name.ValueString(), data.ActionType.ValueString())
	actionTemplateUpdate.ID = at.ID
	actionTemplateUpdate.SpaceID = at.SpaceID
	actionTemplateUpdate.Version = at.Version

	actionTemplateUpdate.Description = data.Description.ValueString()
	if !data.CommunityActionTemplateId.IsNull() {
		actionTemplateUpdate.CommunityActionTemplateID = data.CommunityActionTemplateId.ValueString()
	}

	if len(props) > 0 {
		templateProps := make(map[string]core.PropertyValue, len(props))
		for key, val := range props {
			templateProps[key] = core.NewPropertyValue(val.ValueString(), false)
		}
		actionTemplateUpdate.Properties = templateProps
	} else {
		actionTemplateUpdate.Properties = make(map[string]core.PropertyValue, 0)
	}

	actionTemplateUpdate.Packages = make([]packages.PackageReference, 0, len(pkgs))
	if len(pkgs) > 0 {
		for i, val := range pkgs {
			// 			pkgProps := make(map[string]types.String, len(val.Properties.Attributes()))
			// resp.Diagnostics.Append(val.Properties.ElementsAs(ctx, &pkgProps, false)...)
			// TODO: fix
			if resp.Diagnostics.HasError() {
				return
			}
			pkgRef := packages.PackageReference{
				AcquisitionLocation: val.AcquisitionLocation.ValueString(),
				FeedID:              val.FeedID.ValueString(),
				// 				Properties:          util.ConvertAttrStringMapToStringMap(pkgProps),
				Name:      val.Name.ValueString(),
				PackageID: val.PackageID.ValueString(),
			}
			actionTemplateUpdate.Packages[i] = pkgRef
		}
	}

	actionTemplateUpdate.Parameters = make([]actiontemplates.ActionTemplateParameter, len(params))
	if len(params) > 0 {
		for i, val := range params {
			defaultValue := core.NewPropertyValue(val.DefaultValue.ValueString(), false)
			//			displaySetting := make(map[string]types.String, len(val.DisplaySettings.Elements()))
			resp.Diagnostics.Append(val.DisplaySettings.ElementsAs(ctx, &val.DisplaySettings, false)...)
			if resp.Diagnostics.HasError() {
				return
			}
			actionTemplateUpdate.Parameters[i] = actiontemplates.ActionTemplateParameter{
				DefaultValue: &defaultValue,
				Name:         val.Name.ValueString(),
				Label:        val.Label.ValueString(),
				HelpText:     val.HelpText.ValueString(),
				//				DisplaySettings: util.ConvertAttrStringMapToStringMap(displaySetting),
			}
		}
	}

	updatedActionTemplate, err := actiontemplates.Update(r.Config.Client, actionTemplateUpdate)
	if err != nil {
		resp.Diagnostics.AddError("unable to update environment", err.Error())
		return
	}

	updateStepTemplate(ctx, &data, updatedActionTemplate)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *stepTemplateTypeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.EnvironmentTypeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := environments.DeleteByID(r.Config.Client, data.SpaceID.ValueString(), data.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError("unable to delete environment", err.Error())
		return
	}
}

func updateStepTemplate(ctx context.Context, data *schemas.StepTemplateTypeResourceModel, at *actiontemplates.ActionTemplate) diag.Diagnostics {
	resp := diag.Diagnostics{}

	data.ID = types.StringValue(at.ID)
	data.SpaceID = types.StringValue(at.SpaceID)
	data.Name = types.StringValue(at.Name)
	data.Version = types.Int32Value(at.Version)
	data.Description = types.StringValue(at.Description)
	data.CommunityActionTemplateId = types.StringValue(at.CommunityActionTemplateID)
	data.ActionType = types.StringValue(at.ActionType)

	sParams, dg := convertStepTemplateToParameterAttributes(at.Parameters)
	resp.Append(dg...)
	if resp.HasError() {
		return resp
	}
	data.Parameters = sParams

	props, dg := types.MapValueFrom(ctx, types.StringType, at.Properties)
	resp.Append(dg...)
	if resp.HasError() {
		return resp
	}
	data.Properties = props

	pkgs, dg := convertStepTemplateToPackageAttributes(at.Packages)
	resp.Append(dg...)
	if resp.HasError() {
		return resp
	}
	data.Packages = pkgs

	return resp
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
	sParams, dg := types.ListValue(types.ObjectType{AttrTypes: schemas.GetStepTemplatePackagePropertiesTypeAttributes()}, params)
	resp.Append(dg...)
	if resp.HasError() {
		return types.ListNull(types.ObjectType{AttrTypes: schemas.GetStepTemplateParameterTypeAttributes()}), resp
	}
	return sParams, resp
}

func convertStepTemplateParameterAttribute(atp actiontemplates.ActionTemplateParameter) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(schemas.GetStepTemplateParameterTypeAttributes(), map[string]attr.Value{
		"id":               types.StringValue(atp.ID),
		"name":             types.StringValue(atp.Name),
		"label":            types.StringValue(atp.Label),
		"help_text":        types.StringValue(atp.HelpText),
		"display_settings": types.StringValue(atp.HelpText),
		"default_value":    types.StringValue(atp.DefaultValue.Value),
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

	if extract, ok := atpp["Extract"]; ok {
		prop["extract"] = types.StringValue(extract)
	} else {
		diags.AddWarning("Package property missing value.", "Extract value missing from package property")
		prop["extract"] = types.StringNull()
	}

	if purpose, ok := atpp["Purpose"]; ok {
		prop["purpose"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "Purpose value missing from package property")
		prop["purpose"] = types.StringNull()
	}

	if purpose, ok := atpp["PackageParameterName"]; ok {
		prop["package_parameter_name"] = types.StringValue(purpose)
	} else {
		diags.AddWarning("Package property missing value.", "PackageParameterName value missing from package property")
		prop["package_parameter_name"] = types.StringNull()
	}

	if selectionMode, ok := atpp["SelectionMode"]; ok {
		prop["selection_mode"] = types.StringValue(selectionMode)
	} else {
		diags.AddWarning("Package property missing value.", "SelectionMode value missing from package property")
		prop["package_parameter_name"] = types.StringNull()
	}

	propMap, dg := types.ObjectValue(schemas.GetStepTemplatePackagePropertiesTypeAttributes(), prop)
	if dg.HasError() {
		diags.Append(dg...)
		return types.ObjectNull(schemas.GetStepTemplatePackagePropertiesTypeAttributes()), diags
	}
	return propMap, diags
}
