package octopusdeploy

import (
	"context"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceFeed() *schema.Resource {
	schemaMap := map[string]*schema.Schema{
		constName: {
			Type:     schema.TypeString,
			Required: true,
		},
		constFeedURI: {
			Type:     schema.TypeString,
			Required: true,
		},
		constEnhancedMode: {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  true,
		},
		constDownloadAttempts: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  5,
		},
		constDownloadRetryBackoffSeconds: {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  10,
		},
		constUsername: {
			Type:     schema.TypeString,
			Optional: true,
		},
		constPassword: {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
		},
	}
	schemaMap[constFeedType] = getFeedTypeSchema()

	return &schema.Resource{
		CreateContext: resourceFeedCreate,
		DeleteContext: resourceFeedDelete,
		ReadContext:   resourceFeedRead,
		Schema:        schemaMap,
		UpdateContext: resourceFeedUpdate,
	}
}

func resourceFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)

	feedID := d.Id()
	feed, err := client.Feeds.GetByID(feedID)
	if err != nil {
		return diag.FromErr(err)
	}

	feedResource := feed.(*octopusdeploy.FeedResource)

	d.Set(constDownloadAttempts, feedResource.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, feedResource.DownloadRetryBackoffSeconds)
	d.Set(constEnhancedMode, feedResource.EnhancedMode)
	d.Set(constFeedType, feedResource.FeedType)
	d.Set(constFeedURI, feedResource.FeedURI)
	d.Set(constName, feedResource.Name)

	// TODO
	// d.Set(constPassword, feedResource.Password)

	d.Set(constUsername, feedResource.Username)

	return nil
}

func buildFeedResource(d *schema.ResourceData) *octopusdeploy.FeedResource {
	name := d.Get(constName).(string)

	var feedType octopusdeploy.FeedType
	feedTypeInterface, ok := d.GetOk(constFeedType)
	if ok {
		feedType = octopusdeploy.FeedType(feedTypeInterface.(string))
	}

	var feedURI string
	feedURIInterface, ok := d.GetOk(constFeedURI)
	if ok {
		feedURI = feedURIInterface.(string)
	}

	var feedResource = octopusdeploy.NewFeedResource(name, feedType)
	feedResource.FeedURI = feedURI

	enhancedModeInterface, ok := d.GetOk(constEnhancedMode)
	if ok {
		feedResource.EnhancedMode = enhancedModeInterface.(bool)
	}

	downloadAttemptsInterface, ok := d.GetOk(constDownloadAttempts)
	if ok {
		feedResource.DownloadAttempts = downloadAttemptsInterface.(int)
	}

	downloadRetryBackoffSecondsInterface, ok := d.GetOk(constDownloadRetryBackoffSeconds)
	if ok {
		feedResource.DownloadRetryBackoffSeconds = downloadRetryBackoffSecondsInterface.(int)
	}

	feedUsernameInterface, ok := d.GetOk(constUsername)
	if ok {
		feedResource.Username = feedUsernameInterface.(string)
	}

	feedPasswordInterface, ok := d.GetOk(constPassword)
	if ok {
		feedResource.Password = octopusdeploy.NewSensitiveValue(feedPasswordInterface.(string))
	}

	return feedResource
}

func resourceFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feedResource := buildFeedResource(d)

	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.Add(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(feed.GetID())

	return nil
}

func resourceFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feedResource := buildFeedResource(d)
	feedResource.ID = d.Id()

	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.Update(feedResource)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(feed.GetID())

	return nil
}

func resourceFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(id)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(constEmptyString)

	return nil
}
