package octopusdeploy

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const ContainerFeedCategory = "container"

func resourceContainerFeedCreateReleaseTrigger() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceContainerFeedCreateReleaseTriggerCreate,
		DeleteContext: resourcePackageFeedCreateReleaseTriggerDelete,
		Importer:      getImporter(),
		ReadContext:   resourcePackageFeedCreateReleaseTriggerRead,
		Schema:        getPackageFeedCreateReleaseTriggerSchema(),
		UpdateContext: resourceContainerFeedCreateReleaseTriggerUpdate,
	}
}

func resourceContainerFeedCreateReleaseTriggerCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePackageFeedCreateReleaseTriggerCreate(ctx, d, m, ContainerFeedCategory)
}

func resourceContainerFeedCreateReleaseTriggerUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourcePackageFeedCreateReleaseTriggerUpdate(ctx, d, m, ContainerFeedCategory)
}
