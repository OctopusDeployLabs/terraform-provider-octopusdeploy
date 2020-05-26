package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceCertificate() *schema.Resource {
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
			"certificate_data": {
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
	client := m.(*octopusdeploy.Client)

	certificateID := d.Id()
	certificate, err := client.Certificate.Get(certificateID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading certificate %s: %s", certificateID, err.Error())
	}

	d.Set("name", certificate.Name)
	d.Set("notes", certificate.Notes)
	d.Set("environment_ids", certificate.EnvironmentIds)
	d.Set("tenanted_deployment_participation", certificate.TenantedDeploymentParticipation)
	d.Set("tenant_ids", certificate.TenantIds)
	d.Set("tenant_tags", certificate.TenantTags)

	return nil
}

func buildCertificateResource(d *schema.ResourceData) *octopusdeploy.Certificate {
	certificateName := d.Get("name").(string)

	var notes string
	var certificateData string
	var password string
	var environmentIds []string
	var tenantedDeploymentParticipation string
	var tenantIds []string
	var tenantTags []string

	notesInterface, ok := d.GetOk("notes")
	if ok {
		notes = notesInterface.(string)
	}

	certificateDataInterface, ok := d.GetOk("certificate_data")
	if ok {
		certificateData = certificateDataInterface.(string)
	}

	passwordInterface, ok := d.GetOk("password")
	if ok {
		password = passwordInterface.(string)
	}

	environmentIdsInterface, ok := d.GetOk("environment_ids")
	if ok {
		environmentIds = getSliceFromTerraformTypeList(environmentIdsInterface)
	}

	if environmentIds == nil {
		environmentIds = []string{}
	}

	tenantedDeploymentParticipationInterface, ok := d.GetOk("tenanted_deployment_participation")
	if ok {
		tenantedDeploymentParticipation = tenantedDeploymentParticipationInterface.(string)
	}

	tenantIdsInterface, ok := d.GetOk("tenant_ids")
	if ok {
		tenantIds = getSliceFromTerraformTypeList(tenantIdsInterface)
	}

	if tenantIds == nil {
		tenantIds = []string{}
	}

	tenantTagsInterface, ok := d.GetOk("tenant_tags")
	if ok {
		tenantTags = getSliceFromTerraformTypeList(tenantTagsInterface)
	}

	if tenantTags == nil {
		tenantTags = []string{}
	}

	var certificate = octopusdeploy.NewCertificate(certificateName, octopusdeploy.SensitiveValue{NewValue: certificateData}, octopusdeploy.SensitiveValue{NewValue: password})
	certificate.Notes = notes
	certificate.EnvironmentIds = environmentIds
	certificate.TenantedDeploymentParticipation, _ = octopusdeploy.ParseTenantedDeploymentMode(tenantedDeploymentParticipation)
	certificate.TenantIds = tenantIds
	certificate.TenantTags = tenantTags

	return certificate
}

func resourceCertificateCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newCertificate := buildCertificateResource(d)
	certificate, err := client.Certificate.Add(newCertificate)

	if err != nil {
		return fmt.Errorf("error creating certificate %s: %s", newCertificate.Name, err.Error())
	}

	d.SetId(certificate.ID)

	return nil
}

func resourceCertificateUpdate(d *schema.ResourceData, m interface{}) error {
	certificate := buildCertificateResource(d)
	certificate.ID = d.Id()

	client := m.(*octopusdeploy.Client)

	updatedCertificate, err := client.Certificate.Replace(certificate)

	if err != nil {
		return fmt.Errorf("error updating certificate id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedCertificate.ID)
	return nil
}

func resourceCertificateDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	certificateID := d.Id()

	err := client.Certificate.Delete(certificateID)

	if err != nil {
		return fmt.Errorf("error deleting certificate id %s: %s", certificateID, err.Error())
	}

	d.SetId("")
	return nil
}
