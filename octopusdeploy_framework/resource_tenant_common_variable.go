package octopusdeploy_framework

import (
	"context"
	"fmt"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &tenantCommonVariableResource{}
var _ resource.ResourceWithImportState = &tenantCommonVariableResource{}

type tenantCommonVariableResource struct {
	*Config
}

type tenantCommonVariableResourceModel struct {
	SpaceID              types.String `tfsdk:"space_id"`
	TenantID             types.String `tfsdk:"tenant_id"`
	LibraryVariableSetID types.String `tfsdk:"library_variable_set_id"`
	TemplateID           types.String `tfsdk:"template_id"`
	Value                types.String `tfsdk:"value"`

	schemas.ResourceModel
}

func NewTenantCommonVariableResource() resource.Resource {
	return &tenantCommonVariableResource{}
}

func (t *tenantCommonVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TenantCommonVariableResourceName)
}

func (t *tenantCommonVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetTenantCommonVariableResourceSchema()
}

func (t *tenantCommonVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tenantCommonVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating tenant common variable")

	id := fmt.Sprintf("%s:%s:%s", plan.TenantID.ValueString(), plan.LibraryVariableSetID.ValueString(), plan.TemplateID.ValueString())

	tenant, err := tenants.GetByID(t.Client, plan.SpaceID.ValueString(), plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	err = checkIfCandidateVariableRequiredForTenant(tenant, tenantVariables, plan)
	if err != nil {
		resp.Diagnostics.AddError("Tenant doesn't need a value for this Common Variable", "Tenants must be connected to a Project with an included Library Variable Set that defines Common Variable templates, before common variable values can be provided ("+err.Error()+")")
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantCommonVariable(tenantVariables, plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	plan.ID = types.StringValue(id)
	plan.SpaceID = types.StringValue(tenant.SpaceID)

	tflog.Debug(ctx, "Tenant common variable created", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (t *tenantCommonVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, err := tenants.GetByID(t.Client, state.SpaceID.ValueString(), state.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[state.LibraryVariableSetID.ValueString()]; ok {
		if value, ok := libraryVariable.Variables[state.TemplateID.ValueString()]; ok {
			if !isSensitive {
				state.Value = types.StringValue(value.Value)
			}
		} else {
			resp.State.RemoveResource(ctx)
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (t *tenantCommonVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, err := tenants.GetByID(t.Client, plan.SpaceID.ValueString(), plan.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantCommonVariable(tenantVariables, plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant common variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (t *tenantCommonVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tenantCommonVariableResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tenant, err := tenants.GetByID(t.Client, state.SpaceID.ValueString(), state.TenantID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant", err.Error())
		return
	}

	tenantVariables, err := t.Client.Tenants.GetVariables(tenant)
	if err != nil {
		resp.Diagnostics.AddError("Error retrieving tenant variables", err.Error())
		return
	}

	isSensitive, err := checkIfCommonVariableIsSensitive(tenantVariables, state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[state.LibraryVariableSetID.ValueString()]; ok {
		if isSensitive {
			libraryVariable.Variables[state.TemplateID.ValueString()] = core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}}
		} else {
			delete(libraryVariable.Variables, state.TemplateID.ValueString())
		}
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantCommonVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 3 {
		resp.Diagnostics.AddError(
			"Incorrect Import Format",
			"ID must be in the format: TenantID:LibraryVariableSetID:TemplateID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c)",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("library_variable_set_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), idParts[2])...)
}

func checkIfCommonVariableIsSensitive(tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel) (bool, error) {
	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		for _, template := range libraryVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return template.DisplaySettings["Octopus.ControlType"] == "Sensitive", nil
			}
		}
	}
	return false, fmt.Errorf("unable to find template for tenant variable")
}

func checkIfCandidateVariableRequiredForTenant(tenant *tenants.Tenant, tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel) error {
	if tenant.ProjectEnvironments == nil || len(tenant.ProjectEnvironments) == 0 {
		return fmt.Errorf("tenant not connected to any projects")
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		for _, template := range libraryVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return nil
			}
		}
	} else {
		return fmt.Errorf("tenant not connected to a project that includes variable set " + plan.LibraryVariableSetID.ValueString())
	}

	return fmt.Errorf("common template " + plan.TemplateID.ValueString() + " not found in variable set " + plan.LibraryVariableSetID.ValueString())
}

func updateTenantCommonVariable(tenantVariables *variables.TenantVariables, plan tenantCommonVariableResourceModel, isSensitive bool) error {
	if libraryVariable, ok := tenantVariables.LibraryVariables[plan.LibraryVariableSetID.ValueString()]; ok {
		libraryVariable.Variables[plan.TemplateID.ValueString()] = core.NewPropertyValue(plan.Value.ValueString(), isSensitive)
		return nil
	}
	return fmt.Errorf("unable to locate tenant variable for tenant ID %s", plan.TenantID.ValueString())
}
