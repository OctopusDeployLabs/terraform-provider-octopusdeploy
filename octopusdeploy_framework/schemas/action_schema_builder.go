package schemas

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ planmodifier.Map = warnIfIncludesRunOnServer{}

type warnIfIncludesRunOnServer struct{}

func (lm warnIfIncludesRunOnServer) Description(_ context.Context) string {
	return "Use run_on_server on the action rather than in the properties collection."
}

func (lm warnIfIncludesRunOnServer) MarkdownDescription(ctx context.Context) string {
	return lm.Description(ctx)
}

func (lm warnIfIncludesRunOnServer) PlanModifyMap(ctx context.Context, req planmodifier.MapRequest, resp *planmodifier.MapResponse) {
	const key = "Octopus.Action.RunOnServer"
	properties := req.ConfigValue.Elements()
	if properties[key] != nil {
		resp.Diagnostics.AddWarning(fmt.Sprintf("\"%s\" is defined in properties", key), "Please update your template to specify \"run_on_server\" under the action instead.")
	}
}

type ActionResourceSchemaBuilder struct {
	attributes map[string]schema.Attribute
	blocks     map[string]schema.Block
}

func NewActionResourceSchemaBuilder() *ActionResourceSchemaBuilder {
	builder := &ActionResourceSchemaBuilder{
		attributes: make(map[string]schema.Attribute),
		blocks:     make(map[string]schema.Block),
	}

	builder.attributes["can_be_used_for_project_versioning"] = schema.BoolAttribute{
		Computed: true,
		Optional: true,
	}
	builder.attributes["channels"] = schema.ListAttribute{
		Computed:    true,
		Description: "The channels associated with this deployment action.",
		ElementType: types.StringType,
		Optional:    true,
	}
	builder.attributes["condition"] = schema.StringAttribute{
		Computed:    true,
		Description: "The condition associated with this deployment action.",
		Optional:    true,
	}

	builder.blocks["container"] = schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"feed_id": schema.StringAttribute{
					Optional: true,
				},
				"image": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
	builder.attributes["environments"] = schema.ListAttribute{
		Computed:    true,
		Optional:    true,
		Description: "The environments within which this deployment action will run.",
		ElementType: types.StringType,
	}

	builder.attributes["excluded_environments"] = schema.ListAttribute{
		Computed:    true,
		Optional:    true,
		Description: "The environments that this step will be skipped in",
		ElementType: types.StringType,
	}
	builder.attributes["features"] = schema.ListAttribute{
		Computed:    true,
		Optional:    true,
		ElementType: types.StringType,
		Description: "A list of enabled features for this action.",
	}

	builder.blocks["action_template"] = schema.ListNestedBlock{
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"id": GetIdResourceSchema(),
				"version": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}

	builder.attributes["id"] = util.GetIdResourceSchema()
	builder.blocks["git_dependency"] = schema.SetNestedBlock{
		Description: "Configuration for resource sourcing from a git repository.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"repository_uri": schema.StringAttribute{
					Description: "The Git URI for the repository where this resource is sourced from.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"default_branch": schema.StringAttribute{
					Description: "Name of the default branch of the repository.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"git_credential_type": schema.StringAttribute{
					Description: "The Git credential authentication type.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"file_path_filters": schema.ListAttribute{
					Description: "List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"git_credential_id": schema.StringAttribute{
					Description: "ID of an existing Git credential.",
					Optional:    true,
				},
			},
		},
	}

	builder.attributes["is_disabled"] = schema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the disabled status of this deployment action.",
		Optional:    true,
		Computed:    true,
	}

	builder.attributes["is_required"] = schema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the required status of this deployment action.",
		Optional:    true,
		Computed:    true,
	}
	builder.attributes["name"] = GetNameResourceSchema(false)
	builder.attributes["notes"] = schema.StringAttribute{
		Description: "The notes associated with this deployment action.",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}

	builder.WithPackages()
	builder.WithProperties("")

	builder.attributes["sort_order"] = schema.Int64Attribute{
		Description: "Order used by terraform to ensure correct ordering of actions. This property must be either omitted from all actions, or provided on all actions",
		Optional:    true,
		Computed:    true,
	}

	builder.attributes["computed_sort_order"] = schema.Int64Attribute{
		Description: "This is the final sort order for the action. This will be provided by the API.",
		Optional:    true,
		Computed:    true,
		PlanModifiers: []planmodifier.Int64{
			int64planmodifier.UseStateForUnknown(),
		},
	}

	builder.attributes["slug"] = GetSlugResourceSchema("action")
	builder.attributes["tenant_tags"] = schema.ListAttribute{
		Computed:    true,
		Description: "A list of tenant tags associated with this resource.",
		ElementType: types.StringType,
		Optional:    true,
	}

	return builder
}

func (b *ActionResourceSchemaBuilder) WithActionType() *ActionResourceSchemaBuilder {
	b.attributes["action_type"] = schema.StringAttribute{
		Description: "The type of action",
		Optional:    true,
		Validators: []validator.String{
			validators.ActionTypeHasSpecificImplementation(),
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithExecutionLocation() *ActionResourceSchemaBuilder {
	b.attributes["run_on_server"] = schema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Computed:    true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithWorkerPool() *ActionResourceSchemaBuilder {
	b.attributes["worker_pool_id"] = schema.StringAttribute{
		Description: "The worker pool associated with this deployment action.",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithWorkerPoolVariable() *ActionResourceSchemaBuilder {
	b.attributes["worker_pool_variable"] = schema.StringAttribute{
		Description: "The worker pool variable associated with this deployment action.",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithPackages() *ActionResourceSchemaBuilder {
	b.WithPrimaryPackage()

	additionalAttributes := map[string]schema.Attribute{
		"extract_during_deployment": schema.BoolAttribute{
			Optional:    true,
			Computed:    true,
			Default:     booldefault.StaticBool(true),
			Description: "Whether to extract the package during deployment",
		},
		"name": GetNameResourceSchema(false),
	}

	packageSchema := getPackageSchema(additionalAttributes)

	b.blocks["package"] = packageSchema
	return b
}

func (b *ActionResourceSchemaBuilder) WithScriptFromPackage() *ActionResourceSchemaBuilder {
	b.attributes["script_file_name"] = schema.StringAttribute{
		Description: "The script file name in the package",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}
	b.attributes["script_parameters"] = schema.StringAttribute{
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}
	b.attributes["script_source"] = schema.StringAttribute{
		Computed: true,
		Optional: true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithScript() *ActionResourceSchemaBuilder {
	b.attributes["script_body"] = schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
	}

	b.attributes["script_syntax"] = schema.StringAttribute{
		Computed: true,
		Optional: true,
		Default:  stringdefault.StaticString(""),
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithVariableSubstitutionInFiles() *ActionResourceSchemaBuilder {
	b.attributes["variable_substitution_in_files"] = schema.StringAttribute{
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
		Optional:    true,
		Computed:    true,
		Default:     stringdefault.StaticString(""),
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithNamespace() *ActionResourceSchemaBuilder {
	b.attributes["namespace"] = schema.StringAttribute{
		Optional: true,
		Computed: true,
		Default:  stringdefault.StaticString(""),
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithProperties(deprecated string) *ActionResourceSchemaBuilder {
	b.attributes["properties"] = schema.MapAttribute{
		Computed:           true,
		Description:        "The properties associated with this deployment action.",
		ElementType:        types.StringType,
		Optional:           true,
		DeprecationMessage: deprecated,
		PlanModifiers: []planmodifier.Map{
			warnIfIncludesRunOnServer{},
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithGitDependency() *ActionResourceSchemaBuilder {
	b.blocks["git_dependency"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		Description: "Configuration for resource sourcing from a git repository.",
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"repository_uri": schema.StringAttribute{
					Description: "The Git URI for the repository where this resource is sourced from.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"default_branch": schema.StringAttribute{
					Description: "Name of the default branch of the repository.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"git_credential_type": schema.StringAttribute{
					Description: "The Git credential authentication type.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"file_path_filters": schema.ListAttribute{
					Description: "List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"git_credential_id": schema.StringAttribute{
					Description: "ID of an existing Git credential.",
					Optional:    true,
				},
			},
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithPrimaryPackage() *ActionResourceSchemaBuilder {
	packageSchema := getPackageSchema(nil)
	packageSchema.Validators = append(packageSchema.Validators, listvalidator.SizeAtMost(1))
	b.blocks["primary_package"] = packageSchema
	return b
}

func (b *ActionResourceSchemaBuilder) Build() schema.ListNestedBlock {
	return schema.ListNestedBlock{
		NestedObject: schema.NestedBlockObject{
			Attributes: b.attributes,
			Blocks:     b.blocks,
		},
	}
}

func (b *ActionResourceSchemaBuilder) WithTerraform() *ActionResourceSchemaBuilder {
	b.blocks["advanced_options"] = schema.SetNestedBlock{
		Description: "Optional advanced options for Terraform",
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"allow_additional_plugin_downloads": schema.BoolAttribute{
					Default:  booldefault.StaticBool(true),
					Optional: true,
				},
				"apply_parameters": schema.StringAttribute{
					Optional: true,
				},
				"init_parameters": schema.StringAttribute{
					Optional: true,
				},
				"plugin_cache_directory": schema.StringAttribute{
					Optional: true,
				},
				"workspace": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}

	b.blocks["aws_account"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"region": schema.StringAttribute{
					Optional: true,
				},
				"variable": schema.StringAttribute{
					Optional: true,
				},
				"use_instance_role": schema.BoolAttribute{
					Optional: true,
				},
			},
			Blocks: map[string]schema.Block{
				"role": schema.SetNestedBlock{
					Validators: []validator.Set{
						setvalidator.SizeAtMost(1),
					},
					NestedObject: schema.NestedBlockObject{
						Attributes: map[string]schema.Attribute{
							"arn": schema.StringAttribute{
								Optional: true,
							},
							"external_id": schema.StringAttribute{
								Optional: true,
							},
							"role_session_name": schema.StringAttribute{
								Optional: true,
							},
							"session_duration": schema.Int64Attribute{
								Optional: true,
							},
						},
					},
				},
			},
		},
	}

	b.blocks["azure_account"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},

		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"variable": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}

	b.blocks["google_cloud_account"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"variable": schema.StringAttribute{
					Optional: true,
				},
				"use_vm_service_account": schema.BoolAttribute{
					Optional:    true,
					Description: "When running in a Compute Engine virtual machine, use the associated VM service account",
				},
				"project": schema.StringAttribute{
					Optional:    true,
					Description: "This sets GOOGLE_PROJECT environment variable",
				},
				"region": schema.StringAttribute{
					Optional:    true,
					Description: "This sets GOOGLE_REGION environment variable",
				},
				"zone": schema.StringAttribute{
					Optional:    true,
					Description: "This sets GOOGLE_ZONE environment variable",
				},
				"service_account_emails": schema.StringAttribute{
					Optional:    true,
					Description: "This sets GOOGLE_IMPERSONATE_SERVICE_ACCOUNT environment variable",
				},
				"impersonate_service_account": schema.BoolAttribute{
					Optional:    true,
					Description: "Impersonate service accounts",
				},
			},
		},
	}

	b.blocks["template"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		NestedObject: schema.NestedBlockObject{
			Attributes: map[string]schema.Attribute{
				"additional_variable_files": schema.StringAttribute{
					Optional: true,
				},
				"directory": schema.StringAttribute{
					Optional: true,
				},
				"run_automatic_file_substitution": schema.BoolAttribute{
					Optional: true,
				},
				"target_files": schema.StringAttribute{
					Optional: true,
				},
			},
		},
	}

	b.attributes["template_parameters"] = schema.StringAttribute{Optional: true}
	b.attributes["inline_template"] = schema.StringAttribute{Optional: true}

	return b
}

func (b *ActionResourceSchemaBuilder) WithKubernetesSecret() *ActionResourceSchemaBuilder {
	b.attributes["secret_name"] = schema.StringAttribute{
		Description: "The name of the secret resource",
		Required:    true,
	}

	b.attributes["secret_values"] = schema.MapAttribute{
		ElementType: types.StringType,
		Required:    true,
	}

	b.attributes["kubernetes_object_status_check_enabled"] = schema.BoolAttribute{
		Optional:    true,
		Default:     booldefault.StaticBool(true),
		Description: "Indicates the status of the Kubernetes Object Status feature",
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithWindowsServiceFeature() *ActionResourceSchemaBuilder {
	b.blocks["windows_service"] = schema.SetNestedBlock{
		Validators: []validator.Set{
			setvalidator.SizeAtMost(1),
		},
		Description: "Deploy a windows service feature",
		NestedObject: schema.NestedBlockObject{
			Attributes: getDeployWindowsServiceSchema(),
		},
	}
	return b
}

func (b *ActionResourceSchemaBuilder) WithWindowsService() *ActionResourceSchemaBuilder {
	for key, attribute := range getDeployWindowsServiceSchema() {
		b.attributes[key] = attribute
	}
	return b
}

func (b *ActionResourceSchemaBuilder) WithManualIntervention() *ActionResourceSchemaBuilder {
	b.attributes["instructions"] = schema.StringAttribute{
		Description: "The instructions for the user to follow",
		Required:    true,
	}
	b.attributes["responsible_teams"] = schema.StringAttribute{
		Description: "The teams responsible to resolve this step. If no teams are specified, all users who have permission to deploy the project can resolve it.",
		Optional:    true,
	}
	return b
}

func getDeployWindowsServiceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"arguments": schema.StringAttribute{
			Description: "The command line arguments that will be passed to the service when it starts",
			Optional:    true,
		},
		"create_or_update_service": schema.BoolAttribute{
			Computed: true,
			Optional: true,
		},
		"custom_account_name": schema.StringAttribute{
			Description: "The Windows/domain account of the custom user that the service will run under",
			Optional:    true,
		},
		"custom_account_password": schema.StringAttribute{
			Computed:    true,
			Description: "The password for the custom account",
			Optional:    true,
			Sensitive:   true,
		},
		"dependencies": schema.StringAttribute{
			Description: "Any dependencies that the service has. Separate the names using forward slashes (/).",
			Optional:    true,
		},
		"description": schema.StringAttribute{
			Description: "User-friendly description of the service (optional)",
			Optional:    true,
		},
		"display_name": schema.StringAttribute{
			Description: "The display name of the service (optional)",
			Optional:    true,
		},
		"executable_path": schema.StringAttribute{
			Description: "The path to the executable relative to the package installation directory",
			Required:    true,
		},
		"service_account": schema.StringAttribute{
			Description: "Which built-in account will the service run under. Can be LocalSystem, NT Authority\\NetworkService, NT Authority\\LocalService, _CUSTOM or an expression",
			Default:     stringdefault.StaticString("LocalSystem"),
			Optional:    true,
		},
		"service_name": schema.StringAttribute{
			Description: "The name of the service",
			Required:    true,
		},
		"start_mode": schema.StringAttribute{
			Default:     stringdefault.StaticString("auto"),
			Description: "When will the service start. Can be auto, delayed-auto, manual, unchanged or an expression",
			Optional:    true,
		},
	}
}

func getPackageSchema(additionalAttributes map[string]schema.Attribute) schema.ListNestedBlock {
	attributes := map[string]schema.Attribute{
		"acquisition_location": schema.StringAttribute{
			Default:     stringdefault.StaticString("Server"),
			Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("Server", "ExecutionTarget", "NotAcquired"),
			},
		},
		"feed_id": schema.StringAttribute{
			Default:     stringdefault.StaticString("feeds-builtin"),
			Description: "The feed ID associated with this package reference.",
			Optional:    true,
			Computed:    true,
		},
		"id": GetIdResourceSchema(),
		"package_id": schema.StringAttribute{
			Description: "The ID of the package.",
			Required:    true,
		},
		"properties": schema.MapAttribute{
			Computed:    true,
			Optional:    true,
			Description: "A list of properties associated with this package.",
			ElementType: types.StringType,
		},
	}

	if additionalAttributes != nil {
		for k, v := range additionalAttributes {
			attributes[k] = v
		}
	}

	return schema.ListNestedBlock{
		Description:  "The package associated with this action.",
		NestedObject: schema.NestedBlockObject{Attributes: attributes},
	}
}
