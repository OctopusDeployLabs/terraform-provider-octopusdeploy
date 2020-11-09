package octopusdeploy

import (
	"context"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceNuGetFeed() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNuGetFeedCreate,
		ReadContext:   resourceNuGetFeedRead,
		UpdateContext: resourceNuGetFeedUpdate,
		DeleteContext: resourceNuGetFeedDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Required: true,
				Type:     schema.TypeString,
			},
			constFeedURI: {
				Required: true,
				Type:     schema.TypeString,
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
				Optional: true,
				Type:     schema.TypeString,
			},
			constPassword: {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceNuGetFeedRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(id)
	if err != nil {
		return diag.FromErr(err)
	}
	if feedResource == nil {
		d.SetId("")
		return nil
	}

	logResource(constFeed, m)

	nuGetFeed := feedResource.(*octopusdeploy.NuGetFeed)

	d.Set(constName, nuGetFeed.Name)
	d.Set(constFeedURI, nuGetFeed.FeedURI)
	d.Set(constEnhancedMode, nuGetFeed.EnhancedMode)
	d.Set(constDownloadAttempts, nuGetFeed.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, nuGetFeed.DownloadRetryBackoffSeconds)
	d.Set(constUsername, nuGetFeed.Username)
	d.Set(constPassword, nuGetFeed.Password)

	return nil
}

func buildNuGetFeedResource(d *schema.ResourceData) *octopusdeploy.NuGetFeed {
	feedName := d.Get(constName).(string)

	var feedURI string
	var enhancedMode bool
	var downloadAttempts int
	var downloadRetryBackoffSeconds int
	var feedUsername string
	var feedPassword string

	feedURIInterface, ok := d.GetOk(constFeedURI)
	if ok {
		feedURI = feedURIInterface.(string)
	}

	enhancedModeInterface, ok := d.GetOk(constEnhancedMode)
	if ok {
		enhancedMode = enhancedModeInterface.(bool)
	}

	downloadAttemptsInterface, ok := d.GetOk(constDownloadAttempts)
	if ok {
		downloadAttempts = downloadAttemptsInterface.(int)
	}

	downloadRetryBackoffSecondsInterface, ok := d.GetOk(constDownloadRetryBackoffSeconds)
	if ok {
		downloadRetryBackoffSeconds = downloadRetryBackoffSecondsInterface.(int)
	}

	feedUsernameInterface, ok := d.GetOk(constUsername)
	if ok {
		feedUsername = feedUsernameInterface.(string)
	}

	feedPasswordInterface, ok := d.GetOk(constPassword)
	if ok {
		feedPassword = feedPasswordInterface.(string)
	}

	feed := octopusdeploy.NewNuGetFeed(feedName, feedURI)
	feed.EnhancedMode = enhancedMode
	feed.DownloadAttempts = downloadAttempts
	feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSeconds
	feed.Username = feedUsername
	feed.Password = octopusdeploy.NewSensitiveValue(feedPassword)

	return feed
}

func resourceNuGetFeedCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := buildNuGetFeedResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Feeds.Add(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceNuGetFeedUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	feed := buildNuGetFeedResource(d)
	feed.ID = d.Id() // set ID so Octopus API knows which feed to update

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(updatedFeed.GetID())

	return nil
}

func resourceNuGetFeedDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
