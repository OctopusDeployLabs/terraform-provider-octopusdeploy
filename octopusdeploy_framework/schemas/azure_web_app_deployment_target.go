package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const AzureWebAppDeploymentTargetDescription = "Azure web app deployment target"

type AzureWebAppDeploymentTargetSchema struct{}

type AzureWebAppDeploymentTargetModel struct {
	AccountId         types.String `tfsdk:"account_id"`
	ResourceGroupName types.String `tfsdk:"resource_group_name"`
	WebAppName        types.String `tfsdk:"web_app_name"`
	WebAppSlotName    types.String `tfsdk:"web_app_slot_name"`

	DeploymentTargetModel
}

func (a AzureWebAppDeploymentTargetSchema) GetResourceSchema() resourceSchema.Schema {
	deploymentTargetSchema := DeploymentTargetSchema{}.GetResourceSchema()

	deploymentTargetSchema.Description = util.GetResourceSchemaDescription(AzureWebAppDeploymentTargetDescription)

	deploymentTargetSchema.Attributes["account_id"] = resourceSchema.StringAttribute{
		Required: true,
	}
	deploymentTargetSchema.Attributes["resource_group_name"] = resourceSchema.StringAttribute{
		Required: true,
	}
	deploymentTargetSchema.Attributes["web_app_name"] = resourceSchema.StringAttribute{
		Required: true,
	}
	deploymentTargetSchema.Attributes["web_app_slot_name"] = resourceSchema.StringAttribute{
		Required: true,
	}

	return deploymentTargetSchema
}
