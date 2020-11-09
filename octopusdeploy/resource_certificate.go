package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceCertificate() *schema.Resource {
	validateSchema()

	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			constNotes: {
				Optional: true,
				Type:     schema.TypeString,
			},
			constCertificateData: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			constPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			constEnvironmentIDs: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenanted_deployment_participation": getTenantedDeploymentSchema(),
			"tenant_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tenant_tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceCertificateRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	resource, err := client.Certificates.GetByID(id)
	if err != nil {
		diag.FromErr(err)
	}
	if resource == nil {
		d.SetId("")
		return nil
	}

	logResource(constCertificate, m)

	d.Set("name", resource.Name)
	d.Set(constNotes, resource.Notes)
	d.Set(constEnvironmentIDs, resource.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", resource.TenantedDeploymentMode)
	d.Set("tenant_ids", resource.TenantIDs)
	d.Set("tenant_tags", resource.TenantTags)

	return nil
}

func buildCertificateResource(d *schema.ResourceData) (*octopusdeploy.CertificateResource, error) {
	name := d.Get("name").(string)
	if isEmpty(name) {
		log.Println("certificate name is empty; please specify a name for the certificate")
	}

	password := d.Get(constPassword).(string)
	if isEmpty(password) {
		log.Println("password is empty; please specify a password")
	}

	pass := octopusdeploy.NewSensitiveValue(password)
	certData := d.Get(constCertificateData).(string)
	if isEmpty(certData) {
		log.Println("certificate data is empty; please specify certificate data")
	}

	certificateData := octopusdeploy.NewSensitiveValue(certData)
	certificate := octopusdeploy.NewCertificateResource(name, certificateData, pass)

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		certificate.TenantedDeploymentMode = octopusdeploy.TenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		certificate.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return certificate, nil
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate, err := buildCertificateResource(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := m.(*octopusdeploy.Client)
	resource, err := client.Certificates.Add(certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	certificate, err := buildCertificateResource(d)
	if err != nil {
		return diag.FromErr(err)
	}
	certificate.ID = d.Id() // set ID so Octopus API knows which certificate to update

	client := m.(*octopusdeploy.Client)
	resource, err := client.Certificates.Update(*certificate)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(resource.GetID())

	return nil
}

func resourceCertificateDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Certificates.DeleteByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
