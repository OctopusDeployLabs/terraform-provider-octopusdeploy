package errors

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"log"
	"net/http"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/octopusdeploy_framework/schemas"
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

func DeleteFromStateV2(ctx context.Context, state *tfsdk.State, resource schemas.IResourceModel, resourceDescription string) error {
	log.Printf("[INFO] %s (%s) not found; deleting from state", resourceDescription, resource.GetID())
	state.RemoveResource(ctx)
	return nil
}

func ProcessApiErrorV2(ctx context.Context, state *tfsdk.State, resource schemas.IResourceModel, err error, resourceDescription string) error {
	if err == nil {
		return nil
	}

	if apiError, ok := err.(*core.APIError); ok {
		if apiError.StatusCode == http.StatusNotFound {
			return DeleteFromStateV2(ctx, state, resource, resourceDescription)
		}
	}

	return nil
}
