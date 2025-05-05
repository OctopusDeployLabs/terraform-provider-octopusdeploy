package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbookprocess"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
)

// A wrapper of deployment- or runbook processes.
//
// Provides common "abstraction" so both processes can be managed by the process resource
type processWrapper interface { // Better name?
	GetID() string
	GetSpaceID() string
	GetProjectID() string
	PopulateState(state *schemas.ProcessResourceModel)
	AppendStep(step *deployments.DeploymentStep)
	RemoveStep(stepId string)
	ReplaceSteps(steps []*deployments.DeploymentStep)
	// Update sends underlying process to the server via corresponding API endpoint
	//
	// Returns new instance with updated process, original process remains unchanged
	Update(client *client.Client) (processWrapper, error)
	FindStepByID(stepID string) (*deployments.DeploymentStep, bool)
	FindStepByName(name string) (*deployments.DeploymentStep, bool)
	GetSteps() []*deployments.DeploymentStep
}

func findDeploymentStepByID(steps []*deployments.DeploymentStep, stepID string) (*deployments.DeploymentStep, bool) {
	for _, step := range steps {
		if step.ID == stepID {
			return step, true
		}
	}
	return nil, false
}

func findDeploymentStepByName(steps []*deployments.DeploymentStep, name string) (*deployments.DeploymentStep, bool) {
	for _, step := range steps {
		if step.Name == name {
			return step, true
		}
	}
	return nil, false
}

// loadProcessWrapperByProcessId determines projectId before loading deployment or runbook process.
//
// Returns wrapper of the process or error when process is not found or warning when corresponding project is version controlled.
func loadProcessWrapperByProcessId(client *client.Client, spaceId string, processId string) (processWrapper, diag.Diagnostics) {
	switch kind, ownerId := deconstructProcessIdentifier(processId); kind {
	case "deployment":
		return loadProcessWrapper(client, spaceId, ownerId, processId)
	case "runbook":
		runbook, err := runbooks.GetByID(client, spaceId, ownerId)
		if err != nil {
			runbookNotFound := diag.NewErrorDiagnostic("Unable to load runbook for process", err.Error())
			return nil, diag.Diagnostics{runbookNotFound}
		}

		return loadProcessWrapper(client, spaceId, runbook.ProjectID, processId)
	default:
		invalidIdentifier := diag.NewErrorDiagnostic("Unable to load process", fmt.Sprintf("Invalid process identifier '%s'", processId))
		return nil, diag.Diagnostics{invalidIdentifier}
	}
}

// loadProcessWrapper loads deployment or runbook process and returns a wrapper of the loaded process.
//
// Returns error when process is not found and warning when corresponding project is version controlled.
func loadProcessWrapper(client *client.Client, spaceId string, projectId string, processId string) (processWrapper, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	// Load corresponding project to check if it's version controlled
	project, projectError := projects.GetByID(client, spaceId, projectId)
	if projectError != nil {
		diags.AddError("Unable to load project for the process", projectError.Error())
		return nil, diags
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled {
		diags.AddWarning("Process persisted under version control system", "Version controlled resources will not be modified via terraform")
		return nil, diags
	}

	switch kind, _ := deconstructProcessIdentifier(processId); kind {
	case "deployment":
		process, processError := deployments.GetDeploymentProcessByID(client, spaceId, processId)
		if processError != nil {
			diags.AddError("Unable to load deployment process", processError.Error())
			return nil, diags
		}

		return deploymentProcessWrapper{process}, diags
	case "runbook":
		process, runbookError := runbookprocess.GetByID(client, spaceId, processId)
		if runbookError != nil {
			diags.AddError("Unable to load runbook process", runbookError.Error())
			return nil, diags
		}

		return runbookProcessWrapper{process}, diags
	default:
		diags.AddError("Unable to load process", fmt.Sprintf("Invalid process identifier '%s'", processId))
		return nil, diags
	}
}

// deconstructProcessIdentifier determines what kind of the process given identifier represents.
//
// Returns determined kind and extracted owner identifier.
//
// Relies on the fact that owner id is embedded in the processId
func deconstructProcessIdentifier(processId string) (kind string, owner string) {
	const deploymentPrefix = "deploymentprocess-"
	if strings.HasPrefix(processId, deploymentPrefix) {
		projectId := processId[len(deploymentPrefix):]
		return "deployment", projectId
	}

	const runbookPrefix = "RunbookProcess-"
	if strings.HasPrefix(processId, runbookPrefix) {
		runbookId := processId[len(runbookPrefix):]
		return "runbook", runbookId
	}

	return "unknown", ""
}

type deploymentProcessWrapper struct {
	process *deployments.DeploymentProcess
}

func (w deploymentProcessWrapper) GetID() string {
	return w.process.GetID()
}

func (w deploymentProcessWrapper) GetSpaceID() string {
	return w.process.SpaceID
}

func (w deploymentProcessWrapper) GetProjectID() string {
	return w.process.ProjectID
}

func (w deploymentProcessWrapper) PopulateState(state *schemas.ProcessResourceModel) {
	state.ID = types.StringValue(w.process.ID)
	state.SpaceID = types.StringValue(w.process.SpaceID)
	state.ProjectID = types.StringValue(w.process.ProjectID)
	state.RunbookID = types.StringNull()
}

func (w deploymentProcessWrapper) AppendStep(step *deployments.DeploymentStep) {
	w.process.Steps = append(w.process.Steps, step)
}

func (w deploymentProcessWrapper) RemoveStep(stepId string) {
	var filteredSteps []*deployments.DeploymentStep
	for _, step := range w.process.Steps {
		if stepId != step.GetID() {
			filteredSteps = append(filteredSteps, step)
		}
	}
	w.process.Steps = filteredSteps
}

func (w deploymentProcessWrapper) ReplaceSteps(steps []*deployments.DeploymentStep) {
	w.process.Steps = steps
}

func (w deploymentProcessWrapper) Update(client *client.Client) (processWrapper, error) {
	updated, err := deployments.UpdateDeploymentProcess(client, w.process)
	if err != nil {
		return nil, err
	}

	return deploymentProcessWrapper{updated}, nil
}

func (w deploymentProcessWrapper) FindStepByID(stepID string) (*deployments.DeploymentStep, bool) {
	return findDeploymentStepByID(w.process.Steps, stepID)
}

func (w deploymentProcessWrapper) FindStepByName(name string) (*deployments.DeploymentStep, bool) {
	return findDeploymentStepByName(w.process.Steps, name)
}

func (w deploymentProcessWrapper) GetSteps() []*deployments.DeploymentStep {
	return w.process.Steps
}

type runbookProcessWrapper struct {
	process *runbookprocess.RunbookProcess
}

func (w runbookProcessWrapper) GetID() string {
	return w.process.GetID()
}

func (w runbookProcessWrapper) GetSpaceID() string {
	return w.process.SpaceID
}

func (w runbookProcessWrapper) GetProjectID() string {
	return w.process.ProjectID
}

func (w runbookProcessWrapper) PopulateState(state *schemas.ProcessResourceModel) {
	state.ID = types.StringValue(w.process.ID)
	state.SpaceID = types.StringValue(w.process.SpaceID)
	state.ProjectID = types.StringValue(w.process.ProjectID)
	state.RunbookID = types.StringValue(w.process.RunbookID)
}

func (w runbookProcessWrapper) AppendStep(step *deployments.DeploymentStep) {
	w.process.Steps = append(w.process.Steps, step)
}

func (w runbookProcessWrapper) RemoveStep(stepId string) {
	var filteredSteps []*deployments.DeploymentStep
	for _, step := range w.process.Steps {
		if stepId != step.GetID() {
			filteredSteps = append(filteredSteps, step)
		}
	}
	w.process.Steps = filteredSteps
}

func (w runbookProcessWrapper) ReplaceSteps(steps []*deployments.DeploymentStep) {
	w.process.Steps = steps
}

func (w runbookProcessWrapper) Update(client *client.Client) (processWrapper, error) {
	updated, err := runbookprocess.Update(client, w.process)
	if err != nil {
		return nil, err
	}

	return runbookProcessWrapper{updated}, nil
}

func (w runbookProcessWrapper) FindStepByID(stepID string) (*deployments.DeploymentStep, bool) {
	return findDeploymentStepByID(w.process.Steps, stepID)
}

func (w runbookProcessWrapper) FindStepByName(name string) (*deployments.DeploymentStep, bool) {
	return findDeploymentStepByName(w.process.Steps, name)
}

func (w runbookProcessWrapper) GetSteps() []*deployments.DeploymentStep {
	return w.process.Steps
}
