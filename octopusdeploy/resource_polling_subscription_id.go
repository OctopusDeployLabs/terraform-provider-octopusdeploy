package octopusdeploy

import (
	"context"
	"github.com/OctopusDeploy/terraform-provider-octopusdeploy/internal"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourcePollingSubscriptionId() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePollingSubscriptionIdCreate,
		DeleteContext: resourcePollingSubscriptionIdDelete,
		Description:   "A unique polling subscription ID that can be used by polling tentacles.",
		ReadContext:   resourcePollingSubscriptionIdRead,
		Schema:        getPollingSubscriptionIDSchema(),
	}
}

func resourcePollingSubscriptionIdRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Don't need to do anything as all the values are already in state
	return nil
}

func resourcePollingSubscriptionIdDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}

func resourcePollingSubscriptionIdCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var generatedSubscriptionId = internal.GenerateRandomCryptoString(20)
	d.SetId(generatedSubscriptionId)
	d.Set("polling_uri", "poll://"+generatedSubscriptionId+"/")

	return nil
}
