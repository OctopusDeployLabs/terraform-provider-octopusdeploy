package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceFeedCreate,
		Read:   resourceFeedRead,
		Update: resourceFeedUpdate,
		Delete: resourceFeedDelete,

		Schema: map[string]*schema.Schema{
			constName: {
				Type:     schema.TypeString,
				Required: true,
			},
			constFeedType: {
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
		},
	}
}

func resourceFeedRead(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	feedID := d.Id()
	feed, err := apiClient.Feeds.GetByID(feedID)

	if err != nil {
		d.SetId(constEmptyString)
		return fmt.Errorf(errorReadingFeed, feedID, err.Error())
	}

	d.Set(constName, feed.Name)
	d.Set(constFeedType, feed.FeedType)
	d.Set(constFeedURI, feed.FeedURI)
	d.Set(constEnhancedMode, feed.EnhancedMode)
	d.Set(constDownloadAttempts, feed.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, feed.DownloadRetryBackoffSeconds)
	d.Set(constUsername, feed.Username)
	d.Set(constPassword, feed.Password)

	return nil
}

func buildFeedResource(d *schema.ResourceData) *model.Feed {
	name := d.Get(constName).(string)

	var feedType enum.FeedType
	feedTypeInterface, ok := d.GetOk(constFeedType)
	if ok {
		feedType = feedTypeInterface.(enum.FeedType)
	}

	var feedURI string
	feedURIInterface, ok := d.GetOk(constFeedURI)
	if ok {
		feedURI = feedURIInterface.(string)
	}

	var feed = model.NewFeed(name, feedType, feedURI)

	enhancedModeInterface, ok := d.GetOk(constEnhancedMode)
	if ok {
		feed.EnhancedMode = enhancedModeInterface.(bool)
	}

	downloadAttemptsInterface, ok := d.GetOk(constDownloadAttempts)
	if ok {
		feed.DownloadAttempts = downloadAttemptsInterface.(int)
	}

	downloadRetryBackoffSecondsInterface, ok := d.GetOk(constDownloadRetryBackoffSeconds)
	if ok {
		feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSecondsInterface.(int)
	}

	feedUsernameInterface, ok := d.GetOk(constUsername)
	if ok {
		feed.Username = feedUsernameInterface.(string)
	}

	feedPasswordInterface, ok := d.GetOk(constPassword)
	if ok {
		feed.Password = model.NewSensitiveValue(feedPasswordInterface.(string))
	}

	return feed
}

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	resource := buildFeedResource(d)

	apiClient := m.(*client.Client)
	feed, err := apiClient.Feeds.Add(*resource)
	if err != nil {
		return fmt.Errorf(errorCreatingFeed, resource.Name, err.Error())
	}

	d.SetId(feed.ID)

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	resource := buildFeedResource(d)

	// set ID to inform Octopus API which feed to update
	resource.ID = d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Feeds.Update(*resource)
	if err != nil {
		return fmt.Errorf(errorUpdatingFeed, d.Id(), err.Error())
	}

	d.SetId(resource.ID)
	return nil
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	apiClient := m.(*client.Client)
	err := apiClient.Feeds.DeleteByID(id)
	if err != nil {
		return fmt.Errorf(errorDeletingFeed, id, err.Error())
	}

	d.SetId(constEmptyString)
	return nil
}
