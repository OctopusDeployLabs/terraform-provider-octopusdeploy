package schemas

import (
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	VariableResourceDescription = "variable"
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
	EncryptedValue  string
	KeyFingerprint  string
	PgpKey          string
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
	EncryptedValue:  "encrypted_value",
	KeyFingerprint:  "key_fingerprint",
	PgpKey:          "pgp_key",
}

var VariableTypeNames = struct {
	AmazonWebServicesAccount string
	AzureAccount             string
	GoogleCloudAccount       string
	Certificate              string
	Sensitive                string
	String                   string
	WorkerPool               string
	UsernamePasswordAccount  string
}{
	AmazonWebServicesAccount: "AmazonWebServicesAccount",
	AzureAccount:             "AzureAccount",
	GoogleCloudAccount:       "GoogleCloudAccount",
	Certificate:              "Certificate",
	Sensitive:                "Sensitive",
	String:                   "String",
	WorkerPool:               "WorkerPool",
	UsernamePasswordAccount:  "UsernamePasswordAccount",
}

var variableTypes = []string{
	VariableTypeNames.AmazonWebServicesAccount,
	VariableTypeNames.AzureAccount,
	VariableTypeNames.GoogleCloudAccount,
	VariableTypeNames.Certificate,
	VariableTypeNames.Sensitive,
	VariableTypeNames.String,
	VariableTypeNames.WorkerPool,
	VariableTypeNames.UsernamePasswordAccount,
}

func GetVariableResourceSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Attributes: map[string]resourceSchema.Attribute{
			SchemaAttributeNames.ID:          GetIdResourceSchema(),
			SchemaAttributeNames.Name:        GetNameResourceSchema(true),
			SchemaAttributeNames.Description: GetDescriptionResourceSchema(VariableResourceDescription),
			SchemaAttributeNames.SpaceID:     GetSpaceIdResourceSchema(VariableResourceDescription),
			VariableSchemaAttributeNames.OwnerID: resourceSchema.StringAttribute{
				Optional: true,
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
			VariableSchemaAttributeNames.IsEditable: GetBooleanResourceAttribute(
				"Indicates whether or not this variable is considered editable.",
				true,
				true),
			VariableSchemaAttributeNames.IsSensitive: GetBooleanResourceAttribute(
				"Indicates whether or not this resource is considered sensitive and should be kept secret.",
				false,
				true),
			VariableSchemaAttributeNames.KeyFingerprint: resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			VariableSchemaAttributeNames.PgpKey: resourceSchema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			VariableSchemaAttributeNames.EncryptedValue: resourceSchema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			VariableSchemaAttributeNames.SensitiveValue: resourceSchema.StringAttribute{
				Optional:  true,
				Sensitive: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(VariableSchemaAttributeNames.Value)),
				},
			},
			VariableSchemaAttributeNames.Type: resourceSchema.StringAttribute{
				Required:    true,
				Description: fmt.Sprintf("The type of variable represented by this resource. Valid types are %s", strings.Join(variableTypes, ", ")),
				Validators: []validator.String{
					stringvalidator.OneOf(
						variableTypes...,
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
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	OwnerID        types.String `tfsdk:"owner_id"`
	ProjectID      types.String `tfsdk:"project_id"`
	IsEditable     types.Bool   `tfsdk:"is_editable"`
	IsSensitive    types.Bool   `tfsdk:"is_sensitive"`
	Type           types.String `tfsdk:"type"`
	SensitiveValue types.String `tfsdk:"sensitive_value"`
	Value          types.String `tfsdk:"value"`
	PgpKey         types.String `tfsdk:"pgp_key"`
	KeyFingerprint types.String `tfsdk:"key_fingerprint"`
	EncryptedValue types.String `tfsdk:"encrypted_value"`
	Prompt         types.List   `tfsdk:"prompt"`
	Scope          types.List   `tfsdk:"scope"`
	SpaceID        types.String `tfsdk:"space_id"`
}
