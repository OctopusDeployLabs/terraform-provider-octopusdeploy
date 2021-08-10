package octopusdeploy

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenantProjectVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantProjectVariableCreate,
		DeleteContext: resourceTenantProjectVariableDelete,
		Description:   "This resource manages tenant project variables in Octopus Deploy.",
		Importer:      &schema.ResourceImporter{State: resourceTenantProjectVariableImporter},
		ReadContext:   resourceTenantProjectVariableRead,
		Schema: map[string]*schema.Schema{
			"environment_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"project_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"template_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"tenant_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"value": {
				Default:   "",
				Optional:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
		},
		UpdateContext: resourceTenantProjectVariableUpdate,
	}
}

func resourceTenantProjectVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)
	value := d.Get("value").(string)

	id := tenantID + ":" + projectID + ":" + environmentID + ":" + templateID

	log.Printf("[INFO] creating tenant project variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.ProjectVariables[projectID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
		if environment, ok := projectVariable.Variables[environmentID]; ok {
			environment[templateID] = octopusdeploy.NewPropertyValue(value, isSensitive)
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			d.SetId(id)
			log.Printf("[INFO] tenant project variable created (%s)", d.Id())
			return nil
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant project variable for tenant ID, %s", tenantID)
}

func resourceTenantProjectVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)

	id := tenantID + ":" + projectID + ":" + environmentID + ":" + templateID

	log.Printf("[INFO] deleting tenant project variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] tenant (%s) not found; deleting common variable from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.ProjectVariables[projectID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
		if environment, ok := projectVariable.Variables[environmentID]; ok {
			if isSensitive {
				environment[templateID] = octopusdeploy.PropertyValue{IsSensitive: true, SensitiveValue: &octopusdeploy.SensitiveValue{HasValue: false}}
			} else {
				delete(environment, templateID)
			}
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			log.Printf("[INFO] tenant project variable deleted (%s)", d.Id())
			d.SetId("")
			return nil
		}
	}

	log.Printf("[INFO] tenant project variable not found; deleting from state: %s", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantProjectVariableImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] importing tenant project variable (%s)", d.Id())

	id := d.Id()

	importStrings := strings.Split(id, ":")
	if len(importStrings) != 4 {
		return nil, fmt.Errorf("octopusdeploy_tenant_project_variable import must be in the form of TenantID:ProjectID:EnvironmentID:TemplateID (e.g. Tenants-123:Projects-456:Environments-789:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
	}

	d.Set("tenant_id", importStrings[0])
	d.Set("project_id", importStrings[1])
	d.Set("environment_id", importStrings[2])
	d.Set("template_id", importStrings[3])
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceTenantProjectVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)

	id := tenantID + ":" + projectID + ":" + environmentID + ":" + templateID

	log.Printf("[INFO] reading tenant project variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] tenant (%s) not found; deleting tenant project variable from state", d.Id())
				d.SetId("")
				return nil
			}
		}
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.ProjectVariables[projectID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
		if templates, ok := projectVariable.Variables[environmentID]; ok {
			if !isSensitive {
				d.Set("value", templates[templateID].Value)
			}

			d.SetId(id)
			log.Printf("[INFO] tenant project variable read (%s)", d.Id())
			return nil
		}
	}

	log.Printf("[INFO] tenant project variable not found; deleting from state, %s", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantProjectVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)
	value := d.Get("value").(string)

	id := tenantID + ":" + projectID + ":" + environmentID + ":" + templateID

	log.Printf("[INFO] updating tenant project variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.ProjectVariables[projectID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if projectVariable, ok := tenantVariables.ProjectVariables[projectID]; ok {
		if environment, ok := projectVariable.Variables[environmentID]; ok {
			environment[templateID] = octopusdeploy.NewPropertyValue(value, isSensitive)
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			d.SetId(id)
			log.Printf("[INFO] tenant project variable updated (%s)", d.Id())
			return nil
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant variable for tenant ID, %s", tenantID)
}
