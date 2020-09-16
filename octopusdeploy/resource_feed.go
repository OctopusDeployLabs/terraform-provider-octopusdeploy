package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedCreate,
		Read:   resourceFeedRead,
		Update: resourceFeedUpdate,
		Delete: resourceFeedDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"feed_type": {
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

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	feedID := d.Id()
	feed, err := apiClient.Feeds.Get(feedID)

	if err == client.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading feed %s: %s", feedID, err.Error())
	}

	d.Set("name", feed.Name)
	d.Set("feed_type", feed.FeedType)
	d.Set("feed_uri", feed.FeedURI)
	d.Set("enhanced_mode", feed.EnhancedMode)
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("username", feed.Username)
	d.Set("password", feed.Password)

	return nil
}

func buildFeedResource(d *schema.ResourceData) *model.Feed {
	feedName := d.Get("name").(string)

	var feedType enum.FeedType
	var feedURI string
	var enhancedMode bool
	var downloadAttempts int
	var downloadRetryBackoffSeconds int
	var feedUsername string
	var feedPassword string

	feedTypeInterface, ok := d.GetOk("feed_type")
	if ok {
		feedType = feedTypeInterface.(enum.FeedType)
	}

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

	var feed = model.NewFeed(feedName, feedType, feedURI)
	feed.EnhancedMode = enhancedMode
	feed.DownloadAttempts = downloadAttempts
	feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSeconds
	feed.Username = feedUsername
	feed.Password = model.SensitiveValue{
		NewValue: &feedPassword,
	}

	return feed
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newFeed := buildFeedResource(d)
	feed, err := apiClient.Feeds.Add(*newFeed)

	if err != nil {
		return fmt.Errorf("error creating feed %s: %s", newFeed.Name, err.Error())
	}

	d.SetId(feed.ID)

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildFeedResource(d)
	feed.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)

	updatedFeed, err := apiClient.Feeds.Update(*feed)

	if err != nil {
		return fmt.Errorf("error updating feed id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedFeed.ID)
	return nil
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	feedID := d.Id()

	err := apiClient.Feeds.Delete(feedID)

	if err != nil {
		return fmt.Errorf("error deleting feed id %s: %s", feedID, err.Error())
	}

	d.SetId("")
	return nil
}
