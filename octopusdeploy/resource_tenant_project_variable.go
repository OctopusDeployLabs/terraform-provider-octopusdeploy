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
			"property_value": {
				Required: true,
				Elem:     &schema.Resource{Schema: getPropertyValueSchema()},
				MaxItems: 1,
				Type:     schema.TypeList,
			},
			"template_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"tenant_id": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
		UpdateContext: resourceTenantProjectVariableUpdate,
	}
}

func resourceTenantProjectVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)

	log.Printf("[INFO] creating tenant project variable")

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range tenantVariables.ProjectVariables {
		if v.ProjectID == projectID {
			for k := range v.Variables {
				if k == environmentID {
					propertyValue := expandPropertyValue(d.Get("property_value"))

					tenantVariables.ProjectVariables[projectID].Variables[environmentID][templateID] = *propertyValue
					client.Tenants.UpdateVariables(tenant, tenantVariables)

					d.SetId(tenantID + ":" + projectID + ":" + environmentID + ":" + templateID)
					log.Printf("[INFO] tenant project variable created (%s)", d.Id())
					return nil
				}
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant variable for tenant ID, %s", tenantID)
}

func resourceTenantProjectVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
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

	for _, v := range tenantVariables.ProjectVariables {
		if v.ProjectID == projectID {
			for k := range v.Variables {
				if k == environmentID {
					delete(tenantVariables.ProjectVariables[projectID].Variables[environmentID], templateID)
					client.Tenants.UpdateVariables(tenant, tenantVariables)

					log.Printf("[INFO] tenant project variable deleted (%s)", d.Id())
					d.SetId("")
					return nil
				}
			}
		}
	}

	d.SetId("")
	log.Printf("[INFO] tenant project variable not found; deleting from state: %s", d.Id())
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

	for _, v := range tenantVariables.ProjectVariables {
		if v.ProjectID == projectID {
			for k, value := range v.Variables {
				if k == environmentID {
					if !value[environmentID].IsSensitive && value[environmentID].SensitiveValue == nil {
						propertyValue := value[environmentID]
						d.Set("property_value", flattenPropertyValue(&propertyValue))
					}
					d.SetId(id)

					log.Printf("[INFO] tenant project variable read (%s)", d.Id())
					return nil
				}
			}
		}
	}

	log.Printf("[INFO] tenant project variable not found; deleting from state, %s", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantProjectVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	environmentID := d.Get("environment_id").(string)
	projectID := d.Get("project_id").(string)
	templateID := d.Get("template_id").(string)
	tenantID := d.Get("tenant_id").(string)

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

	for _, v := range tenantVariables.ProjectVariables {
		if v.ProjectID == projectID {
			for k := range v.Variables {
				if k == environmentID {
					propertyValue := expandPropertyValue(d.Get("property_value"))
					tenantVariables.ProjectVariables[projectID].Variables[environmentID][templateID] = *propertyValue
					client.Tenants.UpdateVariables(tenant, tenantVariables)

					d.SetId(id)
					log.Printf("[INFO] tenant project variable updated (%s)", d.Id())
					return nil
				}
			}
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant variable for tenant ID, %s", tenantID)
}
