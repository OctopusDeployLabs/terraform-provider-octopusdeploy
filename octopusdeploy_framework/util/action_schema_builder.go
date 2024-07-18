package util

import (
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/validators"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
	builder.attributes["id"] = GetIdResourceSchema()
	builder.blocks["git_dependency"] = resourceSchema.SingleNestedBlock{
		Description: "Configuration for resource sourcing from a git repository.",
		Attributes: map[string]resourceSchema.Attribute{
			"repository_uri": resourceSchema.StringAttribute {
				Description:      "The Git URI for the repository where this resource is sourced from.",
				Required:         true,
				Validators: []validator.String {
					stringvalidator.LengthAtLeast(1),
				},
			},
			"default_branch": resourceSchema.StringAttribute{
				Description:      "Name of the default branch of the repository.",
				Required:         true,
				Validators: []validator.String {
					stringvalidator.LengthAtLeast(1),
				},
			},
			"git_credential_type": resourceSchema.StringAttribute{
				Description:      "The Git credential authentication type.",
				Required:         true,
				Validators: []validator.String {
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
	}

	builder.attributes["is_disabled"] = resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the disabled status of this deployment action.",
		Optional:    true,
	}

	builder.attributes["is_required"] = resourceSchema.BoolAttribute{
		Default:     booldefault.StaticBool(false),
		Description: "Indicates the required status of this deployment action.",
		Optional:    true,
	}
	builder.attributes["name"] = GetNameResourceSchema(true)
	builder.attributes["notes"] = resourceSchema.StringAttribute{
		Description: "The notes associated with this deployment action.",
		Optional:    true,
	}

	builder.attributes["packages"] =

	return &ActionResourceSchemaBuilder{}
}

func (b *ActionResourceSchemaBuilder) WithActionType() *ActionResourceSchemaBuilder {
	b.attributes["action_type"] = resourceSchema.StringAttribute{
		Description: "The type of action",
		Required:    true,
		Validators: []validator.String{
			validators.ActionTypeHasSpecificImplementation(),
		},
	}

	return b
}

func (b *ActionResourceSchemaBuilder) Build() resourceSchema.SingleNestedBlock {
	return resourceSchema.SingleNestedBlock{
		Attributes: b.attributes,
		Blocks:     b.blocks,
	}
}
