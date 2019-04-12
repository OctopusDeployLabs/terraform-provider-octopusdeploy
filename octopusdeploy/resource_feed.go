package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedCreate,
		Read:   resourceFeedRead,
		Update: resourceFeedUpdate,
		Delete: resourceFeedDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"feed_type": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"feed_uri": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"enhanced_mode": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"download_attempts": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  5,
			},
			"download_retry_backoff_seconds": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  10,
			},
			"username": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	feedId := d.Id()
	feed, err := client.Feed.Get(feedId)

	if err == octopusdeploy.ErrItemNotFound {
		d.SetId("")
		return nil
	}

	if err != nil {
		return fmt.Errorf("error reading feed %s: %s", feedId, err.Error())
	}

	d.Set("name", feed.Name)
	d.Set("feed_type", feed.FeedType)
	d.Set("feed_uri", feed.FeedUri)
	d.Set("enhanced_mode", feed.EnhancedMode)
	d.Set("download_attempts", feed.DownloadAttempts)
	d.Set("download_retry_backoff_seconds", feed.DownloadRetryBackoffSeconds)
	d.Set("username", feed.Username)
	d.Set("password", feed.Password)

	return nil
}

func buildFeedResource(d *schema.ResourceData) *octopusdeploy.Feed {
	feedName := d.Get("name").(string)

	var feedType string
	var feedUri string
	var enhancedMode bool
	var downloadAttempts int
	var downloadRetryBackoffSeconds int
	var feedUsername string
	var feedPassword string

	feedTypeInterface, ok := d.GetOk("feed_type")
	if ok {
		feedType = feedTypeInterface.(string)
	}

	feedUriInterface, ok := d.GetOk("feed_uri")
	if ok {
		feedUri = feedUriInterface.(string)
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

	var feed = octopusdeploy.NewFeed(feedName, feedType, feedUri)
	feed.EnhancedMode = enhancedMode
	feed.DownloadAttempts = downloadAttempts
	feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSeconds
	feed.Username = feedUsername
	feed.Password = octopusdeploy.SensitiveValue{
		NewValue: feedPassword,
	}

	return feed
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	newFeed := buildFeedResource(d)
	feed, err := client.Feed.Add(newFeed)

	if err != nil {
		return fmt.Errorf("error creating feed %s: %s", newFeed.Name, err.Error())
	}

	d.SetId(feed.ID)

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildFeedResource(d)
	feed.ID = d.Id() // set project struct ID so octopus knows which project to update

	client := m.(*octopusdeploy.Client)

	updatedFeed, err := client.Feed.Update(feed)

	if err != nil {
		return fmt.Errorf("error updating feed id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedFeed.ID)
	return nil
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	client := m.(*octopusdeploy.Client)

	feedId := d.Id()

	err := client.Feed.Delete(feedId)

	if err != nil {
		return fmt.Errorf("error deleting feed id %s: %s", feedId, err.Error())
	}

	d.SetId("")
	return nil
}
