package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

const (
	GitCredentialResourceDescription = "Git Credential"
	GitCredentialResourceName        = "git_credential"
	GitCredentialDatasourceName      = "git_credentials"
)

func GetGitCredentialResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Manages a Git credential in Octopus Deploy.",
		Attributes: map[string]resourceSchema.Attribute{
			"id":          util.GetIdResourceSchema(),
			"space_id":    util.GetSpaceIdResourceSchema(GitCredentialResourceDescription),
			"name":        util.GetNameResourceSchema(true),
			"description": util.GetDescriptionResourceSchema(GitCredentialResourceDescription),
			"type": resourceSchema.StringAttribute{
				Optional:    true,
				Description: "The Git credential authentication type.",
			},
			"username": resourceSchema.StringAttribute{
				Required:    true,
				Description: "The username for the Git credential.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"password": resourceSchema.StringAttribute{
				Required:    true,
				Sensitive:   true,
				Description: "The password for the Git credential.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
		},
	}
}

func GetGitCredentialDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":       util.GetIdDatasourceSchema(),
		"space_id": util.GetSpaceIdDatasourceSchema(GitCredentialResourceDescription),
		"name":     util.GetQueryNameDatasourceSchema(),
		"skip":     util.GetQuerySkipDatasourceSchema(),
		"take":     util.GetQueryTakeDatasourceSchema(),
		"git_credentials": datasourceSchema.ListNestedAttribute{
			Computed:    true,
			Description: "A list of Git Credentials that match the filter(s).",
			NestedObject: datasourceSchema.NestedAttributeObject{
				Attributes: GetGitCredentialAttributes(),
			},
		},
	}
}

func GetGitCredentialAttributes() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":          util.GetIdDatasourceSchema(),
		"space_id":    util.GetSpaceIdDatasourceSchema(GitCredentialResourceDescription),
		"name":        util.GetQueryNameDatasourceSchema(),
		"description": util.GetDescriptionDatasourceSchema(GitCredentialResourceDescription),
		"type": datasourceSchema.StringAttribute{
			Computed:    true,
			Description: "The Git credential authentication type.",
		},
		"username": datasourceSchema.StringAttribute{
			Computed:    true,
			Description: "The username for the Git credential.",
		},
	}
}
