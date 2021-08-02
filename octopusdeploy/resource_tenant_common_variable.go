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
			"property_value": {
				Required: true,
				Elem:     &schema.Resource{Schema: getPropertyValueSchema()},
				MaxItems: 1,
				Type:     schema.TypeList,
			},
			"tenant_id": {
				Required: true,
				Type:     schema.TypeString,
			},
			"variable_id": {
				Required: true,
				Type:     schema.TypeString,
			},
		},
		UpdateContext: resourceTenantCommonVariableUpdate,
	}
}

func resourceTenantCommonVariableCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	variableID := d.Get("variable_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + variableID

	log.Printf("[INFO] creating tenant common variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range tenantVariables.LibraryVariables {
		if v.LibraryVariableSetID == libraryVariableSetID {
			propertyValue := expandPropertyValue(d.Get("property_value"))

			tenantVariables.LibraryVariables[libraryVariableSetID].Variables[variableID] = *propertyValue
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			d.SetId(id)
			log.Printf("[INFO] tenant common variable created (%s)", d.Id())
			return nil
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant variable for tenant ID, %s", tenantID)
}

func resourceTenantCommonVariableDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	variableID := d.Get("variable_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + variableID

	log.Printf("[INFO] deleting tenant common variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		if apiError, ok := err.(*octopusdeploy.APIError); ok {
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

	for _, v := range tenantVariables.LibraryVariables {
		if v.LibraryVariableSetID == libraryVariableSetID {
			for variable := range v.Variables {
				if variable == variableID {
					delete(tenantVariables.LibraryVariables[libraryVariableSetID].Variables, variableID)
					client.Tenants.UpdateVariables(tenant, tenantVariables)

					log.Printf("[INFO] tenant common variable deleted (%s)", d.Id())
					d.SetId("")
					return nil
				}
			}
		}
	}

	log.Printf("[INFO] tenant common variable not found; deleting from state: %s", d.Id())
	d.SetId("")
	return nil
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
	d.Set("variable_id", importStrings[2])
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func resourceTenantCommonVariableRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	variableID := d.Get("variable_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + variableID

	log.Printf("[INFO] reading tenant common variable (%s)", id)

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

	for _, v := range tenantVariables.LibraryVariables {
		if v.LibraryVariableSetID == libraryVariableSetID {
			for i := range v.Variables {
				if i == variableID {
					if !v.Variables[i].IsSensitive && v.Variables[i].SensitiveValue == nil {
						propertyValue := v.Variables[i]
						d.Set("property_value", flattenPropertyValue(&propertyValue))
					}

					d.SetId(i)
					log.Printf("[INFO] tenant common variable read (%s)", d.Id())
					return nil
				}
			}
		}
	}

	log.Printf("[INFO] tenant common variable not found; deleting from state, %s", d.Id())
	d.SetId("")
	return nil
}

func resourceTenantCommonVariableUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	libraryVariableSetID := d.Get("library_variable_set_id").(string)
	tenantID := d.Get("tenant_id").(string)
	variableID := d.Get("variable_id").(string)

	id := tenantID + ":" + libraryVariableSetID + ":" + variableID

	log.Printf("[INFO] updating tenant common variable (%s)", id)

	client := m.(*octopusdeploy.Client)
	tenant, err := client.Tenants.GetByID(tenantID)
	if err != nil {
		return diag.FromErr(err)
	}

	tenantVariables, err := client.Tenants.GetVariables(tenant)
	if err != nil {
		return diag.FromErr(err)
	}

	for _, v := range tenantVariables.LibraryVariables {
		if v.LibraryVariableSetID == libraryVariableSetID {
			propertyValue := expandPropertyValue(d.Get("property_value"))
			tenantVariables.LibraryVariables[libraryVariableSetID].Variables[variableID] = *propertyValue
			client.Tenants.UpdateVariables(tenant, tenantVariables)

			d.SetId(id)
			log.Printf("[INFO] tenant common variable updated (%s)", d.Id())
			return nil
		}
	}

	d.SetId("")
	return diag.Errorf("unable to locate tenant common variable for tenant ID, %s", tenantID)
}
