package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceNugetFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceNugetFeedCreate,
		Read:   resourceNugetFeedRead,
		Update: resourceNugetFeedUpdate,
		Delete: resourceNugetFeedDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"feed_uri": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enhanced_mode": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"download_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"download_retry_backoff_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func resourceNugetFeedRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	feedID := d.Id()
	feed, err := client.Feed.Get(feedID)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading feed %s: %s", feedID, err.Error())
	}

	d.Set("name", feed.Name)
	d.Set("feed_uri", feed.FeedUri)
	d.Set("enhanced_mode", feed.EnhancedMode)
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("username", feed.Username)
	d.Set("password", feed.Password)

	return nil
}

func buildNugetFeedResource(d *schema.ResourceData) *octopusdeploy.Feed {
	feedName := d.Get("name").(string)

	var feedURI string
	var enhancedMode bool
	var downloadAttempts int
	var downloadRetryBackoffSeconds int
	var feedUsername string
	var feedPassword string

	feedURIInterface, ok := d.GetOk("feed_uri")
	if ok {
		feedURI = feedURIInterface.(string)
	}

	enhancedModeInterface, ok := d.GetOk("enhanced_mode")
	if ok {
		enhancedMode = enhancedModeInterface.(bool)
	}

	downloadAttemptsInterface, ok := d.GetOk("download_attempts")
	if ok {
		downloadAttempts = downloadAttemptsInterface.(int)
	}

	downloadRetryBackoffSecondsInterface, ok := d.GetOk("download_retry_backoff_seconds")
	if ok {
		downloadRetryBackoffSeconds = downloadRetryBackoffSecondsInterface.(int)
	}

	feedUsernameInterface, ok := d.GetOk("username")
	if ok {
		feedUsername = feedUsernameInterface.(string)
	}

	feedPasswordInterface, ok := d.GetOk("password")
	if ok {
		feedPassword = feedPasswordInterface.(string)
	}

	feed := octopusdeploy.NewFeed(feedName, "NuGet", feedURI)
	feed.EnhancedMode = enhancedMode
	feed.DownloadAttempts = downloadAttempts
	feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSeconds
	feed.Username = feedUsername
	feed.Password = octopusdeploy.SensitiveValue{
		NewValue: feedPassword,
	}

	return feed
}

func resourceNugetFeedCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newFeed := buildNugetFeedResource(d)
	feed, err := client.Feed.Add(newFeed)

	if err != nil {
		return fmt.Errorf("error creating nuget feed %s: %s", newFeed.Name, err.Error())
	}

	d.SetId(feed.ID)

	return nil
}

func resourceNugetFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildNugetFeedResource(d)
	feed.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	updatedFeed, err := client.Feed.Update(feed)

	if err != nil {
		return fmt.Errorf("error updating nuget feed id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedFeed.ID)
	return nil
}

func resourceNugetFeedDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	feedID := d.Id()

	err := client.Feed.Delete(feedID)

	if err != nil {
		return fmt.Errorf("error deleting nuget feed id %s: %s", feedID, err.Error())
	}

	d.SetId("")
	return nil
}
