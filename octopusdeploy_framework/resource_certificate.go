package octopusdeploy_framework

import (
	"context"
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type certificateResource struct {
	*Config
}

func NewCertificateResource() resource.Resource {
	return &certificateResource{}
}

var _ resource.ResourceWithImportState = &certificateResource{}

func (r *certificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("certificate")
}

func (r *certificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.CertificateSchema{}.GetResourceSchema()
}

func (r *certificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	r.Config = ResourceConfiguration(req, resp)
}

func (r *certificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.CertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating certificate", map[string]interface{}{
		"name": plan.Name.ValueString(),
	})

	certificate := expandCertificate(ctx, plan)
	createdCertificate, err := certificates.Add(r.Config.Client, certificate)
	if err != nil {
		resp.Diagnostics.AddError("Error creating certificate`", err.Error())
		return
	}

	state := flattenCertificate(ctx, createdCertificate, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *certificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state schemas.CertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate, err := certificates.GetByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		if err := errors.ProcessApiErrorV2(ctx, resp, state, err, "certificateResource"); err != nil {
			resp.Diagnostics.AddError("unable to load certificate", err.Error())
		}
		return
	}

	newState := flattenCertificate(ctx, certificate, state)
	resp.Diagnostics.Append(resp.State.Set(ctx, newState)...)
	return
}

func (r *certificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan schemas.CertificateModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate := expandCertificate(ctx, plan)
	updatedCertificate, err := certificates.Update(r.Client, certificate)
	if err != nil {
		resp.Diagnostics.AddError("Error updating certificate", err.Error())
		return
	}

	state := flattenCertificate(ctx, updatedCertificate, plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	return
}

func (r *certificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state schemas.CertificateModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := certificates.DeleteByID(r.Client, state.SpaceID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error deleting certificate", err.Error())
		return
	}
}

func (*certificateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func expandCertificate(ctx context.Context, model schemas.CertificateModel) *certificates.CertificateResource {
	var name = model.Name.ValueString()
	var certificateData = core.NewSensitiveValue(model.CertificateData.ValueString())
	var password = core.NewSensitiveValue(model.Password.ValueString())

	certificate := certificates.NewCertificateResource(name, certificateData, password)

	certificate.Archived = model.Archived.ValueString()
	certificate.CertificateDataFormat = model.CertificateDataFormat.ValueString()
	certificate.EnvironmentIDs = expandStringList(model.EnvironmentIDs)
	certificate.HasPrivateKey = model.HasPrivateKey.ValueBool()
	certificate.IsExpired = model.IsExpired.ValueBool()
	certificate.IssuerCommonName = model.IssuerCommonName.ValueString()
	certificate.IssuerDistinguishedName = model.IssuerDistinguishedName.ValueString()
	certificate.IssuerOrganization = model.IssuerOrganization.ValueString()
	certificate.NotAfter = model.NotAfter.ValueString()
	certificate.NotBefore = model.NotBefore.ValueString()
	certificate.Notes = model.Notes.ValueString()
	certificate.ReplacedBy = model.ReplacedBy.ValueString()
	certificate.SelfSigned = model.SelfSigned.ValueBool()
	certificate.SerialNumber = model.SerialNumber.ValueString()
	certificate.SignatureAlgorithmName = model.SignatureAlgorithmName.ValueString()
	certificate.SpaceID = model.SpaceID.ValueString()
	certificate.SubjectAlternativeNames = expandStringList(model.SubjectAlternativeNames)
	certificate.SubjectCommonName = model.SubjectCommonName.ValueString()
	certificate.SubjectDistinguishedName = model.SubjectDistinguishedName.ValueString()
	certificate.SubjectOrganization = model.SubjectOrganization.ValueString()
	certificate.TenantedDeploymentMode = core.TenantedDeploymentMode(model.TenantedDeploymentMode.ValueString())
	certificate.TenantIDs = expandStringList(model.TenantIDs)

	convertedTenantTags, diags := util.SetToStringArray(ctx, model.TenantTags)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error converting tenant tags: %v\n", diags))
	}

	certificate.TenantTags = convertedTenantTags
	certificate.Thumbprint = model.Thumbprint.ValueString()
	certificate.Version = int(model.Version.ValueInt64())

	return certificate
}

func flattenCertificate(ctx context.Context, certificate *certificates.CertificateResource, model schemas.CertificateModel) schemas.CertificateModel {
	model.ID = types.StringValue(certificate.ID)
	model.Archived = types.StringValue(certificate.Archived)
	model.CertificateDataFormat = types.StringValue(certificate.CertificateDataFormat)
	model.EnvironmentIDs = flattenStringList(certificate.EnvironmentIDs, model.EnvironmentIDs)
	model.HasPrivateKey = types.BoolValue(certificate.HasPrivateKey)
	model.IsExpired = types.BoolValue(certificate.IsExpired)
	model.IssuerCommonName = types.StringValue(certificate.IssuerCommonName)
	model.IssuerDistinguishedName = types.StringValue(certificate.IssuerDistinguishedName)
	model.IssuerOrganization = types.StringValue(certificate.IssuerOrganization)
	model.NotAfter = types.StringValue(certificate.NotAfter)
	model.NotBefore = types.StringValue(certificate.NotBefore)
	model.Notes = types.StringValue(certificate.Notes)
	model.ReplacedBy = types.StringValue(certificate.ReplacedBy)
	model.SelfSigned = types.BoolValue(certificate.SelfSigned)
	model.SerialNumber = types.StringValue(certificate.SerialNumber)
	model.SignatureAlgorithmName = types.StringValue(certificate.SignatureAlgorithmName)
	model.SpaceID = types.StringValue(certificate.SpaceID)
	model.SubjectAlternativeNames = flattenStringList(certificate.SubjectAlternativeNames, model.SubjectAlternativeNames)
	model.SubjectCommonName = types.StringValue(certificate.SubjectCommonName)
	model.SubjectDistinguishedName = types.StringValue(certificate.SubjectDistinguishedName)
	model.SubjectOrganization = types.StringValue(certificate.SubjectOrganization)
	model.TenantedDeploymentMode = types.StringValue(string(certificate.TenantedDeploymentMode))
	model.TenantIDs = flattenStringList(certificate.TenantIDs, model.TenantIDs)
	model.Thumbprint = types.StringValue(certificate.Thumbprint)
	model.Version = types.Int64Value(int64(certificate.Version))

	convertedTenantTags, diags := util.SetToStringArray(ctx, model.TenantTags)
	if diags.HasError() {
		tflog.Error(ctx, fmt.Sprintf("Error converting tenant tags: %v\n", diags))
	}

	model.TenantTags = basetypes.SetValue(util.FlattenStringList(convertedTenantTags))

	// Note: We don't flatten the password or certificate data as these values are sensitive and not returned by the API

	return model
}
