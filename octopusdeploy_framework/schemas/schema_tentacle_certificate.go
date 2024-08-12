package schemas

import (
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func GetTentacleCertificateSchema() resourceSchema.Schema {
	return resourceSchema.Schema{
		Description: "Generates a X.509 self-signed certificate for use with a Octopus Deploy Tentacle.",
		Attributes: map[string]resourceSchema.Attribute{
			"id": resourceSchema.StringAttribute{
				Description: "The unique ID for this resource.",
				Computed:    true,
			},
			"base64": resourceSchema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The base64 encoded pfx certificate.",
			},
			"thumbprint": resourceSchema.StringAttribute{
				Computed:    true,
				Description: "The SHA1 sum of the certificate represented in hexadecimal.",
			},
			"dependencies": resourceSchema.MapAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Optional map of dependencies that when modified will trigger a re-creation of this resource.",
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
		},
	}
}

type TentacleCertificateResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Base64       types.String `tfsdk:"base64"`
	Thumbprint   types.String `tfsdk:"thumbprint"`
	Dependencies types.Map    `tfsdk:"dependencies"`
}
