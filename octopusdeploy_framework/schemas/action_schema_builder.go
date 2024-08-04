package schemas

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
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
	attributes map[string]resourceSchema.Attribute
	blocks     map[string]resourceSchema.Block
}

func NewActionResourceSchemaBuilder() *ActionResourceSchemaBuilder {
	builder := &ActionResourceSchemaBuilder{
		attributes: make(map[string]resourceSchema.Attribute),
		blocks:     make(map[string]resourceSchema.Block),
	}

	builder.attributes["can_be_used_for_project_versioning"] = resourceSchema.BoolAttribute{
		Computed: true,
		Optional: true,
	}
	builder.attributes["channels"] = resourceSchema.ListAttribute{
		Computed:    true,
		Description: "The channels associated with this deployment action.",
		ElementType: types.StringType,
		Optional:    true,
	}
	builder.attributes["condition"] = resourceSchema.StringAttribute{
		Computed:    true,
		Description: "The condition associated with this deployment action.",
		Optional:    true,
	}

	builder.blocks["container"] = resourceSchema.ListNestedBlock{
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"feed_id": schema.StringAttribute{
					Optional: true,
				},
				"image": resourceSchema.StringAttribute{
					Optional: true,
				},
			},
		},
	}
	builder.attributes["environments"] = resourceSchema.ListAttribute{
		Computed:    true,
		Optional:    true,
		Description: "The environments within which this deployment action will run.",
		ElementType: types.StringType,
	}

	builder.attributes["excluded_environments"] = resourceSchema.ListAttribute{
		Computed:    true,
		Optional:    true,
		Description: "The environments that this step will be skipped in",
		ElementType: types.StringType,
	}
	builder.attributes["features"] = resourceSchema.ListAttribute{
		Computed:    true,
		Optional:    true,
		ElementType: types.StringType,
		Description: "A list of enabled features for this action.",
	}

	builder.blocks["action_template"] = resourceSchema.ListNestedBlock{
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"id": GetIdResourceSchema(),
				"version": resourceSchema.StringAttribute{
					Optional: true,
				},
			},
		},
	}

	builder.attributes["id"] = util.GetIdResourceSchema()
	builder.blocks["git_dependency"] = resourceSchema.SetNestedBlock{
		Description: "Configuration for resource sourcing from a git repository.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"repository_uri": resourceSchema.StringAttribute{
					Description: "The Git URI for the repository where this resource is sourced from.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"default_branch": resourceSchema.StringAttribute{
					Description: "Name of the default branch of the repository.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"git_credential_type": resourceSchema.StringAttribute{
					Description: "The Git credential authentication type.",
					Required:    true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"file_path_filters": resourceSchema.ListAttribute{
					Description: "List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"git_credential_id": resourceSchema.StringAttribute{
					Description: "ID of an existing Git credential.",
					Optional:    true,
				},
			},
		},
	}

	builder.attributes["is_disabled"] = resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the disabled status of this deployment action.",
		Optional:    true,
		Computed:    true,
	}

	builder.attributes["is_required"] = resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the required status of this deployment action.",
		Optional:    true,
		Computed:    true,
	}
	builder.attributes["name"] = GetNameResourceSchema(false)
	builder.attributes["notes"] = resourceSchema.StringAttribute{
		Description: "The notes associated with this deployment action.",
		Optional:    true,
	}

	builder.WithPackages()
	builder.WithProperties("")

	builder.attributes["sort_order"] = resourceSchema.Int64Attribute{
		Description: "Order used by terraform to ensure correct ordering of actions. This property must be either omitted from all actions, or provided on all actions",
		Optional:    true,
		Default:     int64default.StaticInt64(-1),
		Computed:    true,
	}

	builder.attributes["slug"] = GetSlugResourceSchema("action")
	builder.attributes["tenant_tags"] = resourceSchema.ListAttribute{
		Computed:    true,
		Description: "A list of tenant tags associated with this resource.",
		ElementType: types.StringType,
		Optional:    true,
	}

	return builder
}

func (b *ActionResourceSchemaBuilder) WithActionType() *ActionResourceSchemaBuilder {
	b.attributes["action_type"] = resourceSchema.StringAttribute{
		Description: "The type of action",
		Optional:    true,
		Validators: []validator.String{
			validators.ActionTypeHasSpecificImplementation(),
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithExecutionLocation() *ActionResourceSchemaBuilder {
	b.attributes["run_on_server"] = resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Whether this step runs on a worker or on the target",
		Optional:    true,
		Computed:    true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithWorkerPool() *ActionResourceSchemaBuilder {
	b.attributes["worker_pool_id"] = resourceSchema.StringAttribute{
		Description: "The worker pool associated with this deployment action.",
		Optional:    true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithWorkerPoolVariable() *ActionResourceSchemaBuilder {
	b.attributes["worker_pool_variable"] = resourceSchema.StringAttribute{
		Description: "The worker pool variable associated with this deployment action.",
		Optional:    true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithPackages() *ActionResourceSchemaBuilder {
	b.WithPrimaryPackage()

	additionalAttributes := map[string]resourceSchema.Attribute{
		"extract_during_deployment": resourceSchema.BoolAttribute{
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
	b.attributes["script_file_name"] = resourceSchema.StringAttribute{
		Description: "The script file name in the package",
		Optional:    true,
	}
	b.attributes["script_parameters"] = resourceSchema.StringAttribute{
		Description: "Parameters expected by the script. Use platform specific calling convention. e.g. -Path #{VariableStoringPath} for PowerShell or -- #{VariableStoringPath} for ScriptCS",
		Optional:    true,
	}
	b.attributes["script_source"] = resourceSchema.StringAttribute{
		Computed: true,
		Optional: true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithScript() *ActionResourceSchemaBuilder {
	b.attributes["script_body"] = resourceSchema.StringAttribute{
		Optional: true,
	}

	b.attributes["script_syntax"] = resourceSchema.StringAttribute{
		Computed: true,
		Optional: true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithVariableSubstitutionInFiles() *ActionResourceSchemaBuilder {
	b.attributes["variable_substitution_in_files"] = resourceSchema.StringAttribute{
		Description: "A newline-separated list of file names to transform, relative to the package contents. Extended wildcard syntax is supported.",
		Optional:    true,
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithProperties(deprecated string) *ActionResourceSchemaBuilder {
	b.attributes["properties"] = resourceSchema.MapAttribute{
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
	b.blocks["git_dependency"] = resourceSchema.ListNestedBlock{
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		Description: "Configuration for resource sourcing from a git repository.",
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: map[string]resourceSchema.Attribute{
				"repository_uri": resourceSchema.StringAttribute{
					Description: "The Git URI for the repository where this resource is sourced from.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"default_branch": resourceSchema.StringAttribute{
					Description: "Name of the default branch of the repository.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"git_credential_type": resourceSchema.StringAttribute{
					Description: "The Git credential authentication type.",
					Required:    true,
					Validators: []validator.String{
						validators.NotWhitespace(),
					},
				},
				"file_path_filters": resourceSchema.ListAttribute{
					Description: "List of file path filters used to narrow down the directory where files are to be sourced from. Supports glob patten syntax.",
					Optional:    true,
					ElementType: types.StringType,
				},
				"git_credential_id": resourceSchema.StringAttribute{
					Description: "ID of an existing Git credential.",
					Optional:    true,
				},
			},
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) WithPrimaryPackage() *ActionResourceSchemaBuilder {
	b.blocks["primary_package"] = getPackageSchema(nil)

	return b
}

func (b *ActionResourceSchemaBuilder) Build() resourceSchema.ListNestedBlock {
	return resourceSchema.ListNestedBlock{
		Validators: []validator.List{
			listvalidator.SizeAtMost(1),
		},
		NestedObject: resourceSchema.NestedBlockObject{
			Attributes: b.attributes,
			Blocks:     b.blocks,
		},
	}
}

func getPackageSchema(additionalAttributes map[string]resourceSchema.Attribute) resourceSchema.ListNestedBlock {
	attributes := map[string]resourceSchema.Attribute{
		"acquisition_location": resourceSchema.StringAttribute{
			Default:     stringdefault.StaticString("Server"),
			Description: "Whether to acquire this package on the server ('Server'), target ('ExecutionTarget') or not at all ('NotAcquired'). Can be an expression",
			Optional:    true,
			Computed:    true,
			Validators: []validator.String{
				stringvalidator.OneOf("Server", "ExecutionTarget", "NotAcquired"),
			},
		},
		"feed_id": resourceSchema.StringAttribute{
			Default:     stringdefault.StaticString("feeds-builtin"),
			Description: "The feed ID associated with this package reference.",
			Optional:    true,
			Computed:    true,
		},
		"id": GetIdResourceSchema(),
		"package_id": resourceSchema.StringAttribute{
			Description: "The ID of the package.",
			Required:    true,
		},
		"properties": resourceSchema.MapAttribute{
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

	return resourceSchema.ListNestedBlock{
		Description:  "The package associated with this action.",
		NestedObject: resourceSchema.NestedBlockObject{Attributes: attributes},
	}
}
