package octopusdeploy

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/tenants"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal/errors"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceTenantCommonVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceTenantCommonVariableCreate,
		DeleteContext: resourceTenantCommonVariableDelete,
		Description:   "This resource manages tenant common variables in Octopus Deploy.",
		Importer:      &schema.ResourceImporter{State: resourceTenantCommonVariableImporter},
		ReadContext:   resourceTenantCommonVariableRead,
		Schema: map[string]*schema.Schema{
			"library_variable_set_id": {
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
			"space_id": {
				Optional: true,
				Computed: true,
				Type:     schema.TypeString,
			},
			"value": {
				Default:   "",
				Optional:  true,
				Sensitive: true,
				Type:      schema.TypeString,
			},
		},
		UpdateContext: resourceTenantCommonVariableUpdate,
	}
}

func resourceTenantCommonVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	templateID := d.Get("template_id").(string)
	spaceID := d.Get("space_id").(string)
	value := d.Get("value").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + templateID

	log.Printf("[INFO] creating tenant common variable (%s)", id)

	client := m.(*client.Client)
	tenant, err := tenants.GetByID(client, spaceID, tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.LibraryVariables[libraryVariableSetID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
		libraryVariable.Variables[templateID] = core.NewPropertyValue(value, isSensitive)
		client.Tenants.UpdateVariables(tenant, tenantVariables)

		d.SetId(id)
		log.Printf("[INFO] tenant common variable created (%s)", d.Id())
		return nil
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant common variable for tenant ID, %s", tenantID)
}

func resourceTenantCommonVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	spaceID := d.Get("space_id").(string)
	templateID := d.Get("template_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + templateID

	log.Printf("[INFO] deleting tenant common variable (%s)", id)

	client := m.(*client.Client)
	tenant, err := tenants.GetByID(client, spaceID, tenantID)
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
			if apiError.StatusCode == 404 {
				log.Printf("[INFO] tenant (%s) not found; deleting tenant common variable from state", d.Id())
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
	for _, template := range tenantVariables.LibraryVariables[libraryVariableSetID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
		if _, ok := libraryVariable.Variables[templateID]; ok {
			if isSensitive {
				libraryVariable.Variables[templateID] = core.PropertyValue{IsSensitive: true, SensitiveValue: &core.SensitiveValue{HasValue: false}}
			} else {
				delete(libraryVariable.Variables, templateID)
			}
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			log.Printf("[INFO] tenant common variable deleted (%s)", d.Id())
			d.SetId("")
			return nil
		}
	}

	return errors.DeleteFromState(ctx, d, "tenant common variable")
}

func resourceTenantCommonVariableImporter(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	log.Printf("[INFO] importing tenant common variable (%s)", d.Id())

	id := d.Id()

	importStrings := strings.Split(id, ":")
	if len(importStrings) != 3 {
		return nil, fmt.Errorf("octopusdeploy_tenant_common_variable import must be in the form of TenantID:LibraryVariableSetID:VariableID (e.g. Tenants-123:LibraryVariableSets-456:6c9f2ba3-3ccd-407f-bbdf-6618e4fd0a0c")
	}

	d.Set("tenant_id", importStrings[0])
	d.Set("library_variable_set_id", importStrings[1])
	d.Set("template_id", importStrings[2])
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceTenantCommonVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	spaceID := d.Get("space_id").(string)
	templateID := d.Get("template_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + templateID

	log.Printf("[INFO] reading tenant common variable (%s)", id)

	client := m.(*client.Client)
	tenant, err := tenants.GetByID(client, spaceID, tenantID)
	if err != nil {
		if apiError, ok := err.(*core.APIError); ok {
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
	for _, template := range tenantVariables.LibraryVariables[libraryVariableSetID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
		if template, ok := libraryVariable.Variables[templateID]; ok {
			if !isSensitive {
				d.Set("value", template.Value)
			}

			d.SetId(id)
			log.Printf("[INFO] tenant common variable read (%s)", d.Id())
			return nil
		}
	}

	return errors.DeleteFromState(ctx, d, "tenant common variable")
}

func resourceTenantCommonVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	mutex.Lock()
	defer mutex.Unlock()

	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	spaceID := d.Get("space_id").(string)
	templateID := d.Get("template_id").(string)
	value := d.Get("value").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + templateID

	log.Printf("[INFO] updating tenant common variable (%s)", id)

	client := m.(*client.Client)
	tenant, err := tenants.GetByID(client, spaceID, tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	isSensitive := false
	for _, template := range tenantVariables.LibraryVariables[libraryVariableSetID].Templates {
		if template.GetID() == templateID {
			isSensitive = template.DisplaySettings["Octopus.ControlType"] == "Sensitive"
		}
	}

	if libraryVariable, ok := tenantVariables.LibraryVariables[libraryVariableSetID]; ok {
		libraryVariable.Variables[templateID] = core.NewPropertyValue(value, isSensitive)
		client.Tenants.UpdateVariables(tenant, tenantVariables)

		d.SetId(id)
		log.Printf("[INFO] tenant common variable updated (%s)", d.Id())
		return nil
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant common variable for tenant ID, %s", tenantID)
}
