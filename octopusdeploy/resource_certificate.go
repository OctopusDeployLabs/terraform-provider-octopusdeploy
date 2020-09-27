package octopusdeploy

import (
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
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constNotes: {
				Type:     schema.TypeString,
				Optional: true,
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
			constTenantedDeploymentParticipation: getTenantedDeploymentSchema(),
			constTenantIDs: {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			constTenantTags: {
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
	certificateID := d.Id()

	apiClient := m.(*client.Client)
	certificate, err := apiClient.Certificates.GetByID(certificateID)

	if err != nil {
		return createResourceOperationError(errorReadingCertificate, certificateID, err)
	}
	if certificate == nil {
		d.SetId(constEmptyString)
		return nil
	}

	d.Set(constName, certificate.Name)
	d.Set(constNotes, certificate.Notes)
	d.Set(constEnvironmentIDs, certificate.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, certificate.TenantedDeploymentParticipation)
	d.Set(constTenantIDs, certificate.TenantIds)
	d.Set(constTenantTags, certificate.TenantTags)

	return nil
}

func buildCertificateResource(d *schema.ResourceData) (*model.Certificate, error) {
	name := d.Get(constName).(string)
	if isEmpty(name) {
		log.Println("certificate name is empty; please specify a name for the certificate")
	}

	password := d.Get(constPassword).(string)
	if isEmpty(password) {
		log.Println("password is empty; please specify a password")
	}

	pass := model.NewSensitiveValue(password)
	certData := d.Get(constCertificateData).(string)
	if isEmpty(certData) {
		log.Println("certificate data is empty; please specify certificate data")
	}

	certificateData := model.NewSensitiveValue(certData)
	certificate, err := model.NewCertificate(name, certificateData, pass)
	if err != nil {
		log.Println(err)
	}

	if v, ok := d.GetOk(constTenantedDeploymentParticipation); ok {
		certificate.TenantedDeploymentParticipation, _ = enum.ParseTenantedDeploymentMode(v.(string))
	}

	if v, ok := d.GetOk(constTenantTags); ok {
		certificate.TenantTags = getSliceFromTerraformTypeList(v)
	}

	return certificate, nil
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	newCertificate, err := buildCertificateResource(d)
	if err != nil {
		return err
	}

	apiClient := m.(*client.Client)
	certificate, err := apiClient.Certificates.Add(newCertificate)
	if err != nil {
		return createResourceOperationError(errorCreatingCertificate, newCertificate.Name, err)
	}

	if isEmpty(certificate.ID) {
		log.Println("ID is empty")
	} else {
		d.SetId(certificate.ID)
	}

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	certificate, err := buildCertificateResource(d)
	if err != nil {
		return err
	}

	if isEmpty(certificate.ID) {
		log.Println("ID is empty")
	} else {
		certificate.ID = d.Id()
	}

	apiClient := m.(*client.Client)
	updatedCertificate, err := apiClient.Certificates.Update(*certificate)
	if err != nil {
		return createResourceOperationError(errorUpdatingCertificate, d.Id(), err)
	}

	if isEmpty(updatedCertificate.ID) {
		log.Println("ID is empty")
	} else {
		d.SetId(updatedCertificate.ID)
	}

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	certificateID := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Certificates.DeleteByID(certificateID)
	if err != nil {
		return createResourceOperationError(errorDeletingCertificate, certificateID, err)
	}

	d.SetId(constEmptyString)
	return nil
}
