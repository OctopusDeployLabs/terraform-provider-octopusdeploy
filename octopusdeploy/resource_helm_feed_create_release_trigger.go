package octopusdeploy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const HelmFeedCategory = "helm"

func resourceHelmFeedCreateReleaseTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceHelmFeedCreateReleaseTriggerCreate,
		DeleteContext: resourcePackageFeedCreateReleaseTriggerDelete,
		Importer:      getImporter(),
		ReadContext:   resourcePackageFeedCreateReleaseTriggerRead,
		Schema:        getPackageFeedCreateReleaseTriggerSchema(),
		UpdateContext: resourceHelmFeedCreateReleaseTriggerUpdate,
	}
}

func resourceHelmFeedCreateReleaseTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePackageFeedCreateReleaseTriggerCreate(ctx, d, m, HelmFeedCategory)
}

func resourceHelmFeedCreateReleaseTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePackageFeedCreateReleaseTriggerUpdate(ctx, d, m, HelmFeedCategory)
}
