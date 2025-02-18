package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
	"strings"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ resource.Resource = &processStepResource{}

type processStepResource struct {
	*Config
}

func NewProcessStepResource() resource.Resource {
	return &processStepResource{}
}

func (r *processStepResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName(schemas.ProcessStepResourceName)
}

func (r *processStepResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.ProcessStepSchema{}.GetResourceSchema()
}

func (r *processStepResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *processStepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("creating process step: %s", data.Name.ValueString()))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("Error creating process step, unable to find a process", err.Error())
		return
	}

	step := deployments.NewDeploymentStep(data.Name.ValueString())

	diagnostics := mapFromStateToProcessStep(ctx, data, step)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	process.Steps = append(process.Steps, step)

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to create process step", err.Error())
		return
	}

	createdStep, exists := findStepFromProcessByName(updatedProcess, step.Name)
	if !exists {
		resp.Diagnostics.AddError("unable to create process step '%s'", step.Name)
		return
	}

	mapFromProcessStepToState(updatedProcess, createdStep, data)

	tflog.Info(ctx, fmt.Sprintf("process step created (%s)", data.ID))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("reading process step (%s)", data.ID))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to find process", err.Error())
		return
	}

	step, exists := findStepFromProcessByID(process, stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find process step '%s'", stepId)
		return
	}

	mapFromProcessStepToState(process, step, data)

	tflog.Info(ctx, fmt.Sprintf("process step read (%s)", step.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("updating process step (%s)", stepId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	step, exists := findStepFromProcessByID(process, stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find process step '%s'", stepId)
		return
	}

	diagnostics := mapFromStateToProcessStep(ctx, data, step)
	if diagnostics.HasError() {
		resp.Diagnostics.Append(diagnostics...)
		return
	}

	updatedProcess, err := deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to update process step", err.Error())
		return
	}

	updatedStep, exists := findStepFromProcessByID(updatedProcess, stepId)
	if !exists {
		resp.Diagnostics.AddError("unable to find updated process step '%s'", stepId)
		return
	}

	mapFromProcessStepToState(updatedProcess, updatedStep, data)

	tflog.Info(ctx, fmt.Sprintf("process step updated (%s)", updatedStep.GetID()))
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *processStepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *schemas.ProcessStepResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	spaceId := data.SpaceID.ValueString()
	processId := data.ProcessID.ValueString()
	stepId := data.ID.ValueString()

	internal.KeyedMutex.Lock(processId)
	defer internal.KeyedMutex.Unlock(processId)

	tflog.Info(ctx, fmt.Sprintf("deleting process step (%s)", stepId))

	client := r.Config.Client
	process, err := deployments.GetDeploymentProcessByID(client, spaceId, processId)
	if err != nil {
		resp.Diagnostics.AddError("unable to load process", err.Error())
		return
	}

	var filteredSteps []*deployments.DeploymentStep
	for _, step := range process.Steps {
		if stepId != step.GetID() {
			filteredSteps = append(filteredSteps, step)
		}
	}
	process.Steps = filteredSteps

	_, err = deployments.UpdateDeploymentProcess(client, process)
	if err != nil {
		resp.Diagnostics.AddError("unable to delete process step", err.Error())
		return
	}

	resp.State.RemoveResource(ctx)
}

func mapFromStateToProcessStep(ctx context.Context, state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) diag.Diagnostics {
	step.Condition = deployments.DeploymentStepConditionType(state.Condition.ValueString())
	step.StartTrigger = deployments.DeploymentStepStartTrigger(state.StartTrigger.ValueString())

	targetRoles, diagnostics := util.SetToStringArray(ctx, state.TargetRoles)
	if diagnostics.HasError() {
		return diagnostics
	}
	step.Properties["Octopus.Action.TargetRoles"] = core.NewPropertyValue(strings.Join(targetRoles, ","), false)

	step.Properties["Octopus.Action.MaxParallelism"] = core.NewPropertyValue(state.WindowSize.ValueString(), false)

	mapFromStateToProcessStepFirstAction(state, step)

	return nil
}

func mapFromStateToProcessStepFirstAction(state *schemas.ProcessStepResourceModel, step *deployments.DeploymentStep) {
	actionType := state.ActionType.ValueString()
	name := state.Name.ValueString()

	if step.Actions == nil || len(step.Actions) == 0 {
		newAction := deployments.NewDeploymentAction(name, actionType)
		step.Actions = []*deployments.DeploymentAction{newAction}
	}

	if step.Actions[0] == nil {
		step.Actions[0] = deployments.NewDeploymentAction(name, actionType)
	}

	mapFromStateToProcessStepAction(state, step.Actions[0])
}

func mapFromStateToProcessStepAction(state *schemas.ProcessStepResourceModel, action *deployments.DeploymentAction) {
	action.Name = state.Name.ValueString()
	action.ActionType = state.ActionType.ValueString()

	runOnServer := "False"
	if state.RunOnServer.ValueBool() {
		runOnServer = "True"
	}
	action.Properties["Octopus.Action.RunOnServer"] = core.NewPropertyValue(runOnServer, false)

	action.Properties["Octopus.Action.Script.Syntax"] = core.NewPropertyValue(state.ScriptSyntax.ValueString(), false)
	action.Properties["Octopus.Action.Script.ScriptBody"] = core.NewPropertyValue(state.ScriptBody.ValueString(), false)
}

func mapFromProcessStepToState(process *deployments.DeploymentProcess, step *deployments.DeploymentStep, state *schemas.ProcessStepResourceModel) {
	state.ID = types.StringValue(step.GetID())
	state.SpaceID = types.StringValue(process.SpaceID)

	state.Condition = types.StringValue(string(step.Condition))
	state.StartTrigger = types.StringValue(string(step.StartTrigger))

	parsedTargetRoles := strings.Split(step.Properties["Octopus.Action.TargetRoles"].Value, ",")
	targetRoles := make([]attr.Value, len(parsedTargetRoles))
	for i, value := range parsedTargetRoles {
		targetRoles[i] = types.StringValue(value)
	}
	state.TargetRoles, _ = types.SetValue(types.StringType, targetRoles)
	state.WindowSize = types.StringValue(step.Properties["Octopus.Action.MaxParallelism"].Value)

	if len(step.Actions) > 0 && step.Actions[0] != nil {
		mapFromProcessStepActionToState(step.Actions[0], state)
	}
}

func mapFromProcessStepActionToState(action *deployments.DeploymentAction, state *schemas.ProcessStepResourceModel) {
	state.ActionType = types.StringValue(action.ActionType)

	value, _ := strconv.ParseBool(action.Properties["Octopus.Action.RunOnServer"].Value)
	state.RunOnServer = types.BoolValue(value)

	state.ScriptSyntax = types.StringValue(action.Properties["Octopus.Action.Script.Syntax"].Value)
	state.ScriptBody = types.StringValue(action.Properties["Octopus.Action.Script.ScriptBody"].Value)
}

func findStepFromProcessByID(process *deployments.DeploymentProcess, stepID string) (*deployments.DeploymentStep, bool) {
	for _, step := range process.Steps {
		if step.ID == stepID {
			return step, true
		}
	}
	return nil, false
}

func findStepFromProcessByName(process *deployments.DeploymentProcess, name string) (*deployments.DeploymentStep, bool) {
	for _, step := range process.Steps {
		if step.Name == name {
			return step, true
		}
	}
	return nil, false
}
