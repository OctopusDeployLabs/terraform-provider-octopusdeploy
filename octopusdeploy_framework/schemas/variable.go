package schemas

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"strings"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	datasourceSchema "github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	VariableResourceDescription    = "variable"
	VariablesDataSourceDescription = "variables"
)

var VariableSchemaAttributeNames = struct {
	Prompt          string
	OwnerID         string
	ProjectID       string
	Value           string
	SensitiveValue  string
	Scope           string
	IsEditable      string
	IsSensitive     string
	Type            string
	DisplaySettings string
	ControlType     string
	SelectOption    string
	DisplayName     string
	IsRequired      string
	Label           string
}{
	Prompt:          "prompt",
	OwnerID:         "owner_id",
	ProjectID:       "project_id",
	Value:           "value",
	SensitiveValue:  "sensitive_value",
	Scope:           "scope",
	IsEditable:      "is_editable",
	IsSensitive:     "is_sensitive",
	Type:            "type",
	DisplaySettings: "display_settings",
	ControlType:     "control_type",
	SelectOption:    "select_option",
	DisplayName:     "display_name",
	IsRequired:      "is_required",
	Label:           "label",
}

var VariableTypeNames = struct {
	AmazonWebServicesAccount string
	AzureAccount             string
	GoogleCloudAccount       string
	UsernamePasswordAccount  string
	Certificate              string
	Sensitive                string
	String                   string
	WorkerPool               string
}{
	AmazonWebServicesAccount: "AmazonWebServicesAccount",
	AzureAccount:             "AzureAccount",
	GoogleCloudAccount:       "GoogleCloudAccount",
	UsernamePasswordAccount:  "UsernamePasswordAccount",
	Certificate:              "Certificate",
	Sensitive:                "Sensitive",
	String:                   "String",
	WorkerPool:               "WorkerPool",
}

var VariableTypes = []string{
	VariableTypeNames.AmazonWebServicesAccount,
	VariableTypeNames.AzureAccount,
	VariableTypeNames.GoogleCloudAccount,
	VariableTypeNames.UsernamePasswordAccount,
	VariableTypeNames.Certificate,
	VariableTypeNames.Sensitive,
	VariableTypeNames.String,
	VariableTypeNames.WorkerPool,
}

type VariableSchema struct{}

var _ EntitySchema = VariableSchema{}

func (v VariableSchema) GetDatasourceSchema() datasourceSchema.Schema {
	return datasourceSchema.Schema{
		Description: util.GetDataSourceDescription(VariablesDataSourceDescription),
		Attributes: map[string]datasourceSchema.Attribute{
			//request
			SchemaAttributeNames.Name: datasourceSchema.StringAttribute{
				Required:    true,
				Description: "The name of variable to find.",
			},
			VariableSchemaAttributeNames.OwnerID: datasourceSchema.StringAttribute{
				Required:    true,
				Description: "Owner ID for the variable to find.",
			},
			SchemaAttributeNames.SpaceID: GetSpaceIdDatasourceSchema(VariableResourceDescription, false),

			//response
			SchemaAttributeNames.ID: datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "The identifier of the variable to find.",
			},
			SchemaAttributeNames.Description: datasourceSchema.StringAttribute{
				Computed:    true,
				Description: "The description of this variable.",
			},
			VariableSchemaAttributeNames.IsEditable: datasourceSchema.BoolAttribute{
				Description: "Indicates whether or not this variable is considered editable.",
				Computed:    true,
			},
			VariableSchemaAttributeNames.IsSensitive: datasourceSchema.BoolAttribute{
				Description: "Indicates whether or not this resource is considered sensitive and should be kept secret.",
				Computed:    true,
			},
			VariableSchemaAttributeNames.Type: datasourceSchema.StringAttribute{
				Computed:    true,
				Description: fmt.Sprintf("The type of variable represented by this resource. Valid types are %s.", strings.Join(util.Map(VariableTypes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
			},
			VariableSchemaAttributeNames.SensitiveValue: datasourceSchema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			VariableSchemaAttributeNames.Value: datasourceSchema.StringAttribute{
				Computed: true,
			},
			VariableSchemaAttributeNames.Prompt: getVariablePromptDatasourceSchema(),
			VariableSchemaAttributeNames.Scope:  getVariableScopeDatasourceSchema(),
		},
	}
}

func (v VariableSchema) GetResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: util.GetResourceSchemaDescription(VariableResourceDescription),
		Attributes: map[string]resourceSchema.Attribute{
			SchemaAttributeNames.ID:          GetIdResourceSchema(),
			SchemaAttributeNames.Name:        GetNameResourceSchema(true),
			SchemaAttributeNames.Description: GetDescriptionResourceSchema(VariableResourceDescription),
			SchemaAttributeNames.SpaceID:     GetSpaceIdResourceSchema(VariableResourceDescription),
			VariableSchemaAttributeNames.OwnerID: resourceSchema.StringAttribute{
				Description: "Owner ID for the variable(e.g., project ID or library variable set ID)",
				Optional:    true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(VariableSchemaAttributeNames.ProjectID)),
				},
			},
			VariableSchemaAttributeNames.ProjectID: resourceSchema.StringAttribute{
				DeprecationMessage: fmt.Sprintf("This attribute is deprecated; please use %s instead.", VariableSchemaAttributeNames.OwnerID),
				Optional:           true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(VariableSchemaAttributeNames.OwnerID)),
				},
			},
			VariableSchemaAttributeNames.IsEditable: resourceSchema.BoolAttribute{
				Default:            booldefault.StaticBool(true),
				Description:        "Indicates whether or not this variable is considered editable.",
				Optional:           true,
				Computed:           true,
				DeprecationMessage: "This attribute will change to readonly in the future; please do not manually provide this value as it is not intended to be user managed, any value set will be ignored.",
			},
			VariableSchemaAttributeNames.IsSensitive: GetOptionalBooleanResourceAttribute("Indicates whether or not this resource is considered sensitive and should be kept secret.", false),
			VariableSchemaAttributeNames.SensitiveValue: resourceSchema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(VariableSchemaAttributeNames.Value)),
				},
			},
			VariableSchemaAttributeNames.Type: resourceSchema.StringAttribute{
				Required:    true,
				Description: fmt.Sprintf("The type of variable represented by this resource. Valid types are %s.", strings.Join(util.Map(VariableTypes, func(item string) string { return fmt.Sprintf("`%s`", item) }), ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(
						VariableTypes...,
					),
				},
			},
			VariableSchemaAttributeNames.Value: resourceSchema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(VariableSchemaAttributeNames.SensitiveValue)),
				},
			},
		},
		Blocks: map[string]resourceSchema.Block{
			VariableSchemaAttributeNames.Prompt: getVariablePromptResourceSchema(),
			VariableSchemaAttributeNames.Scope:  getVariableScopeResourceSchema(),
		},
	}
}

type VariableTypeResourceModel struct {
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	OwnerID        types.String `tfsdk:"owner_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	IsEditable     types.Bool   `tfsdk:"is_editable"`
	IsSensitive    types.Bool   `tfsdk:"is_sensitive"`
	Type           types.String `tfsdk:"type"`
	SensitiveValue types.String `tfsdk:"sensitive_value"`
	Value          types.String `tfsdk:"value"`
	Prompt         types.List   `tfsdk:"prompt"`
	Scope          types.List   `tfsdk:"scope"`
	SpaceID        types.String `tfsdk:"space_id"`

	ResourceModel
}

type VariablesDataSourceModel struct {
	OwnerID        types.String `tfsdk:"owner_id"`
	Name           types.String `tfsdk:"name"`
	Scope          types.List   `tfsdk:"scope"`
	SpaceID        types.String `tfsdk:"space_id"`
	Description    types.String `tfsdk:"description"`
	IsEditable     types.Bool   `tfsdk:"is_editable"`
	IsSensitive    types.Bool   `tfsdk:"is_sensitive"`
	Prompt         types.List   `tfsdk:"prompt"`
	SensitiveValue types.String `tfsdk:"sensitive_value"`
	Type           types.String `tfsdk:"type"`
	Value          types.String `tfsdk:"value"`

	ResourceModel
}
