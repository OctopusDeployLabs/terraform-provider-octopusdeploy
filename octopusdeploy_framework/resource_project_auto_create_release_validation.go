package octopusdeploy_framework

import (
	"context"
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/channels"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/deployments"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/packages"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r *projectAutoCreateReleaseResource) validateAutoCreateReleaseConfiguration(ctx context.Context, project *projects.Project, data *schemas.ProjectAutoCreateReleaseResourceModel) error {
	// Validate channel exists
	if err := r.validateChannelExists(ctx, project.SpaceID, data.ChannelID.ValueString()); err != nil {
		return err
	}

	// Validate release creation package configuration
	if err := r.validateReleaseCreationPackageConfiguration(ctx, project, data); err != nil {
		return err
	}

	return nil
}

func (r *projectAutoCreateReleaseResource) validateChannelExists(ctx context.Context, spaceID, channelID string) error {
	if channelID == "" {
		return fmt.Errorf("channel_id is required")
	}

	_, err := channels.GetByID(r.Client, spaceID, channelID)
	if err != nil {
		return fmt.Errorf("channel with ID %s does not exist: %w", channelID, err)
	}

	return nil
}

func (r *projectAutoCreateReleaseResource) validateReleaseCreationPackageConfiguration(ctx context.Context, project *projects.Project, data *schemas.ProjectAutoCreateReleaseResourceModel) error {
	if len(data.ReleaseCreationPackage) == 0 {
		return fmt.Errorf("release_creation_package block is required")
	}

	if len(data.ReleaseCreationPackage) > 1 {
		return fmt.Errorf("only one release_creation_package block is supported")
	}

	pkg := data.ReleaseCreationPackage[0]
	actionName := pkg.DeploymentAction.ValueString()
	packageName := pkg.PackageReference.ValueString()

	if actionName == "" {
		return fmt.Errorf("deployment_action is required in release_creation_package block")
	}

	if packageName == "" {
		return fmt.Errorf("package_reference is required in release_creation_package block")
	}

	// Validate the deployment process contains the specified action and package
	foundAction, foundPackage, stepID, err := r.findDeploymentActionAndPackage(ctx, project, actionName, packageName)
	if err != nil {
		return err
	}

	if foundAction == nil {
		return fmt.Errorf("deployment action '%s' not found in project deployment process", actionName)
	}

	if foundPackage == nil {
		return fmt.Errorf("package reference '%s' not found in deployment action '%s'", packageName, actionName)
	}

	// Set the computed step ID if not provided
	if data.ReleaseCreationPackageStepID.IsNull() || data.ReleaseCreationPackageStepID.IsUnknown() {
		data.ReleaseCreationPackageStepID = types.StringValue(stepID)
	}

	return nil
}

func (r *projectAutoCreateReleaseResource) findDeploymentActionAndPackage(ctx context.Context, project *projects.Project, actionName, packageName string) (*deployments.DeploymentAction, *packages.PackageReference, string, error) {
	// Get the deployment process
	deploymentProcess, err := deployments.GetDeploymentProcessByID(r.Client, project.SpaceID, project.DeploymentProcessID)
	if err != nil {
		return nil, nil, "", fmt.Errorf("unable to read deployment process: %w", err)
	}

	// Find the deployment action and package
	for _, step := range deploymentProcess.Steps {
		for _, action := range step.Actions {
			if action.Name == actionName {
				// Look for the package in the action's packages
				for _, pkg := range action.Packages {
					if pkg.Name == packageName {
						return action, pkg, step.ID, nil
					}
				}
				return action, nil, step.ID, nil
			}
		}
	}

	return nil, nil, "", nil
}

func (r *projectAutoCreateReleaseResource) isAutoCreateReleaseConfigured(project *projects.Project) bool {
	if !project.AutoCreateRelease {
		return false
	}

	if project.ReleaseCreationStrategy == nil {
		return false
	}

	if project.ReleaseCreationStrategy.ReleaseCreationPackage == nil {
		return false
	}

	pkg := project.ReleaseCreationStrategy.ReleaseCreationPackage
	if pkg.DeploymentAction == "" || pkg.PackageReference == "" {
		return false
	}

	return true
}
