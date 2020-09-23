package octopusdeploy

import (
	"errors"
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceCertificate() *schema.Resource {
	validateSchema()

	return &schema.Resource{
		Create: resourceCertificateCreate,
		Read:   resourceCertificateRead,
		Update: resourceCertificateUpdate,
		Delete: resourceCertificateDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"notes": {
				Type:     schema.TypeString,
				Optional: true,
			},
			CertificateData: {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"environment_ids": {
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

func resourceCertificateRead(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceCertificateRead", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceCertificateRead", "m")
	}

	apiClient := m.(*client.Client)

	certificateID := d.Id()
	certificate, err := apiClient.Certificates.Get(certificateID)

	if certificate.Validate() == nil {
		d.SetId("")
		return nil
	}

	err1 := errors.New("Validation on certificate struct: unsucessful")
	log.Println(err1)

	if err != nil {
		return fmt.Errorf("error reading certificate %s: %s", certificateID, err.Error())
	}

	d.Set("name", certificate.Name)
	d.Set("notes", certificate.Notes)
	d.Set("environment_ids", certificate.EnvironmentIDs)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentParticipation)
	d.Set("tenant_ids", certificate.TenantIds)
	d.Set("tenant_tags", certificate.TenantTags)

	return nil
}

func buildCertificateResource(d *schema.ResourceData) (*model.Certificate, error) {
	if d == nil {
		return nil, createInvalidParameterError("buildCertificateResource", "d")
	}

	certificateStruct := model.Certificate{}
	if certificateStruct.Name == "" {
		log.Println("Name struct is nil")
	}

	certificateName := d.Get("name").(string)

	password := d.Get("password").(string)
	if password == "" {
		log.Println("Password is nil. Must add in password")
	}

	pass := model.NewSensitiveValue(password)

	certData := d.Get(CertificateData).(string)
	if certData == "" {
		log.Println("Certificate data is nil. Must add cert data")
	}

	certificateData := model.NewSensitiveValue(certData)

	certificate, err := model.NewCertificate(certificateName, certificateData, pass)
	if err != nil {
		log.Println(err)
	}

	if v, ok := d.GetOk("tenanted_deployment_participation"); ok {
		certificate.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk("tenant_tags"); ok {
		certificate.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return certificate, nil
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceCertificateCreate", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceCertificateCreate", "m")
	}

	apiClient := m.(*client.Client)

	newCertificate, err := buildCertificateResource(d)
	if err != nil {
		log.Println(err)
		return err
	}

	certificate, err := apiClient.Certificates.Add(newCertificate)

	if err != nil {
		return fmt.Errorf("error creating certificate %s: %s", newCertificate.Name, err.Error())
	}

	if certificate.ID == "" {
		log.Println("ID is nil")
	} else {
		d.SetId(certificate.ID)
	}

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceCertificateUpdate", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceCertificateUpdate", "m")
	}

	certificate, err := buildCertificateResource(d)
	if err != nil {
		return err
	}

	if certificate.ID == "" {
		log.Println("ID is nil")
	} else {
		certificate.ID = d.Id()
	}

	apiClient := m.(*client.Client)

	updatedCertificate, err := apiClient.Certificates.Update(*certificate)

	if err != nil {
		return fmt.Errorf("error updating certificate id %s: %s", d.Id(), err.Error())
	}
	if updatedCertificate.ID == "" {
		log.Println("ID is nil")
	} else {
		d.SetId(updatedCertificate.ID)
	}

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	if d == nil {
		return createInvalidParameterError("resourceCertificateDelete", "d")
	}

	if m == nil {
		return createInvalidParameterError("resourceCertificateDelete", "m")
	}

	apiClient := m.(*client.Client)

	certificateID := d.Id()

	err := apiClient.Certificates.Delete(certificateID)

	if err != nil {
		return fmt.Errorf("error deleting certificate id %s: %s", certificateID, err.Error())
	}

	d.SetId("")
	return nil
}
