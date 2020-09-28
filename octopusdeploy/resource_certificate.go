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
	id := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Certificates.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingCertificate, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constCertificate, m)

	d.Set(constName, resource.Name)
	d.Set(constNotes, resource.Notes)
	d.Set(constEnvironmentIDs, resource.EnvironmentIDs)
	d.Set(constTenantedDeploymentParticipation, resource.TenantedDeploymentParticipation)
	d.Set(constTenantIDs, resource.TenantIDs)
	d.Set(constTenantTags, resource.TenantTags)

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
	certificate, err := buildCertificateResource(d)
	if err != nil {
		return err
	}

	apiClient := m.(*client.Client)
	resource, err := apiClient.Certificates.Add(certificate)
	if err != nil {
		return createResourceOperationError(errorCreatingCertificate, certificate.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
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
		certificate.ID = d.Id() // set ID so Octopus API knows which certificate to update
	}

	apiClient := m.(*client.Client)
	resource, err := apiClient.Certificates.Update(*certificate)
	if err != nil {
		return createResourceOperationError(errorUpdatingCertificate, d.Id(), err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is empty")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Certificates.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingCertificate, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
