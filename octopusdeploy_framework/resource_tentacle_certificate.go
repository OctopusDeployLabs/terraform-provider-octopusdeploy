package octopusdeploy_framework

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"software.sslmate.com/src/go-pkcs12"
	"strings"
	"time"

	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/util"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type tentacleCertificateResource struct {
	*Config
}

func NewTentacleCertificateResource() resource.Resource {
	return &tentacleCertificateResource{}
}

func (t *tentacleCertificateResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = util.GetTypeName("tentacle_certificate")
}

func (t *tentacleCertificateResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schemas.GetTentacleCertificateSchema()
}

func (t *tentacleCertificateResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	t.Config = ResourceConfiguration(req, resp)
}

func (t *tentacleCertificateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan schemas.TentacleCertificateResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	certificate, thumbprint, err := generateCertificate("Octopus Tentacle")

	if err != nil {
		util.AddDiagnosticError(resp.Diagnostics, t.Config.SystemInfo, "cannot generate tentacle", err.Error())
		return
	}

	plan.Base64 = types.StringValue(certificate)
	plan.Thumbprint = types.StringValue(thumbprint)
	plan.ID = types.StringValue(internal.GenerateRandomCryptoString(20))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

}

func (t *tentacleCertificateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	return
}

func (t *tentacleCertificateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data schemas.TentacleCertificateResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	data.ID = types.StringValue("")
	return
}

func (r *tentacleCertificateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	return
}

func generateCertificate(fullName string) (string, string, error) {
	random := rand.Reader

	privateKey, err := rsa.GenerateKey(random, 2048)
	if err != nil {
		return "", "", err
	}

	serialNumber := internal.GenerateRandomSerialNumber()
	template := x509.Certificate{
		SerialNumber: &serialNumber,
		Subject: pkix.Name{
			CommonName: fullName,
		},
		Issuer: pkix.Name{
			CommonName: fullName,
		},
		NotBefore: time.Now().AddDate(0, 0, -1),
		NotAfter:  time.Now().AddDate(100, 0, 0),
		KeyUsage:  x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
		},
		IsCA:                  false,
		BasicConstraintsValid: true,
	}

	certBytes, err := x509.CreateCertificate(random, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", "", err
	}

	parsedCert, _ := x509.ParseCertificate(certBytes)
	pkcs12Bytes, err := pkcs12.Passwordless.Encode(privateKey, parsedCert, nil, "")
	if err != nil {
		return "", "", err
	}

	pkcs12Base64 := base64.StdEncoding.EncodeToString(pkcs12Bytes)

	thumbprint := sha1.Sum(certBytes)
	thumbprintStr := strings.ToUpper(hex.EncodeToString(thumbprint[:]))

	return pkcs12Base64, thumbprintStr, nil
}
