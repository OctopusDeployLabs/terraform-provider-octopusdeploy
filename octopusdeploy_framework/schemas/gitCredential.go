package schemas

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

func GetGitCredentialResourceSchema() map[string]resourceschema.Attribute {
	return map[string]resourceschema.Attribute{
		"id":          util.GetIdResourceSchema(),
		"space_id":    util.GetSpaceIdResourceSchema("Git Credentials"),
		"name":        util.GetNameResourceSchema(true),
		"description": util.GetDescriptionResourceSchema("Git Credentials"),
		"type": resourceschema.StringAttribute{
			Computed:    true,
			Description: "The Git credential authentication type.",
		},
		"username": resourceschema.StringAttribute{
			Required:    true,
			Description: "The username for the Git credential.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"password": resourceschema.StringAttribute{
			Required:    true,
			Sensitive:   true,
			Description: "The password for the Git credential.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
	}
}

func GetGitCredentialDataSourceSchema() map[string]datasourceSchema.Attribute {
	return map[string]datasourceSchema.Attribute{
		"id":       util.GetIdDatasourceSchema(),
		"space_id": util.GetSpaceIdDatasourceSchema("Git credentials"),
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
		"space_id":    util.GetSpaceIdDatasourceSchema("Git credentials"),
		"name":        util.GetQueryNameDatasourceSchema(),
		"description": util.GetDescriptionDatasourceSchema("Git credentials"),
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
