package errors

import (
	"context"
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DeleteFromState(ctx context.Context, d *schema.ResourceData, resource string) diag.Diagnostics {
	log.Printf("[INFO] %s (%s) not found; deleting from state", resource, d.Id())
	d.SetId("")
	return nil
}

func ProcessApiError(ctx context.Context, d *schema.ResourceData, err error, resource string) diag.Diagnostics {
	if err == nil {
		return nil
	}

	if apiError, ok := err.(*core.APIError); ok {
		if apiError.StatusCode == http.StatusNotFound {
			return DeleteFromState(ctx, d, resource)
		}
	}

	return diag.FromErr(err)
}
