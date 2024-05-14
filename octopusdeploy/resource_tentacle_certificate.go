package octopusdeploy

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/hex"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"software.sslmate.com/src/go-pkcs12"
	"strings"
	"time"
)

func resourceTentacleCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTentacleCertificateCreate,
		DeleteContext: resourceTentacleCertificateDelete,
		Description:   "Generates a X.509 self-signed certificate for use with a Octopus Deploy Tentacle.",
		ReadContext:   resourceTentacleCertificateRead,
		Schema:        getTentacleCertificateSchema(),
	}
}

func resourceTentacleCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Don't need to do anything as all the values are already in state
	return nil
}

func resourceTentacleCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourceTentacleCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate, thumbprint, err := generateCertificate("Octopus Tentacle")

	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("base64", certificate)
	d.Set("thumbprint", thumbprint)

	d.SetId(generateRandomCryptoString(20))
	return nil
}

func generateCertificate(fullName string) (string, string, error) {
	random := rand.Reader

	privateKey, err := rsa.GenerateKey(random, 2048)
	if err != nil {
		return "", "", err
	}

	serialNumber := generateRandomSerialNumber()
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
