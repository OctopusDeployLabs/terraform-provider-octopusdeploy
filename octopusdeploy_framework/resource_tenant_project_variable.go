package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/variables"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
)

var _ resource.Resource = &tenantProjectVariableResource{}
var _ resource.ResourceWithImportState = &tenantProjectVariableResource{}

type tenantProjectVariableResource struct {
	*Config
}

type tenantProjectVariableResourceModel struct {
	ID            types.String `tfsdk:"id"`
	SpaceID       types.String `tfsdk:"space_id"`
	TenantID      types.String `tfsdk:"tenant_id"`
	ProjectID     types.String `tfsdk:"project_id"`
	EnvironmentID types.String `tfsdk:"environment_id"`
	TemplateID    types.String `tfsdk:"template_id"`
	Value         types.String `tfsdk:"value"`
}

func NewTenantProjectVariableResource() resource.Resource {
	return &tenantProjectVariableResource{}
}

func (t *tenantProjectVariableResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.TenantProjectVariableResourceName)
}

func (t *tenantProjectVariableResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetTenantProjectVariableResourceSchema()
}

func (t *tenantProjectVariableResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tenantProjectVariableResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan tenantProjectVariableResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating tenant project variable")

	id := fmt.Sprintf("%s:%s:%s:%s", plan.TenantID.ValueString(), plan.ProjectID.ValueString(), plan.EnvironmentID.ValueString(), plan.TemplateID.ValueString())

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

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantProjectVariable(tenantVariables, plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	plan.ID = types.StringValue(id)

	tflog.Debug(ctx, "Tenant project variable created", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (t *tenantProjectVariableResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state tenantProjectVariableResourceModel
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

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[state.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[state.EnvironmentID.ValueString()]; ok {
			if value, ok := environment[state.TemplateID.ValueString()]; ok {
				if !isSensitive {
					state.Value = types.StringValue(value.Value)
				}
			} else {
				resp.State.RemoveResource(ctx)
				return
			}
		}
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (t *tenantProjectVariableResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan tenantProjectVariableResourceModel
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

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, plan)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if err := updateTenantProjectVariable(tenantVariables, plan, isSensitive); err != nil {
		resp.Diagnostics.AddError("Error updating tenant project variable", err.Error())
		return
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (t *tenantProjectVariableResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state tenantProjectVariableResourceModel
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

	isSensitive, err := checkIfVariableIsSensitive(tenantVariables, state)
	if err != nil {
		resp.Diagnostics.AddError("Error checking if variable is sensitive", err.Error())
		return
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[state.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[state.EnvironmentID.ValueString()]; ok {
			if isSensitive {
				environment[state.TemplateID.ValueString()] = core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}}
			} else {
				delete(environment, state.TemplateID.ValueString())
			}
		}
	}

	_, err = t.Client.Tenants.UpdateVariables(tenant, tenantVariables)
	if err != nil {
		resp.Diagnostics.AddError("Error updating tenant variables", err.Error())
		return
	}
}

func (t *tenantProjectVariableResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ":")

	if len(idParts) != 4 {
		resp.Diagnostics.AddError(
			"Incorrect Import Format",
			"ID must be in the format: TenantID:ProjectID:EnvironmentID:TemplateID (e.g. Tenants-123:Projects-456:Environments-789:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c)",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("tenant_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("project_id"), idParts[1])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("environment_id"), idParts[2])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("template_id"), idParts[3])...)
}

func checkIfVariableIsSensitive(tenantVariables *variables.TenantVariables, plan tenantProjectVariableResourceModel) (bool, error) {
	if projectVariable, ok := tenantVariables.ProjectVariables[plan.ProjectID.ValueString()]; ok {
		for _, template := range projectVariable.Templates {
			if template.GetID() == plan.TemplateID.ValueString() {
				return template.DisplaySettings["Octopus.ControlType"] == "Sensitive", nil
			}
		}
	}
	return false, fmt.Errorf("unable to find template for tenant variable")
}

func updateTenantProjectVariable(tenantVariables *variables.TenantVariables, plan tenantProjectVariableResourceModel, isSensitive bool) error {
	if projectVariable, ok := tenantVariables.ProjectVariables[plan.ProjectID.ValueString()]; ok {
		if environment, ok := projectVariable.Variables[plan.EnvironmentID.ValueString()]; ok {
			environment[plan.TemplateID.ValueString()] = core.NewPropertyValue(plan.Value.ValueString(), isSensitive)
			return nil
		}
	}
	return fmt.Errorf("unable to locate tenant variable for tenant ID %s", plan.TenantID.ValueString())
}
