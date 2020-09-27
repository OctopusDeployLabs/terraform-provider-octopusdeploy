package octopusdeploy

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceNugetFeed() *schema.Resource {
	return &schema.Resource{
		Create: resourceNugetFeedCreate,
		Read:   resourceNugetFeedRead,
		Update: resourceNugetFeedUpdate,
		Delete: resourceNugetFeedDelete,

		Schema: map[string]*schema.Schema{
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
		},
	}
}

func resourceNugetFeedRead(d *schema.ResourceData, m interface{}) error {
	feedID := d.Id()

	apiClient := m.(*client.Client)
	resource, err := apiClient.Feeds.GetByID(feedID)
	if err != nil {
		return fmt.Errorf(errorReadingFeed, feedID, err.Error())
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	d.Set(constName, resource.Name)
	d.Set(constFeedURI, resource.FeedURI)
	d.Set(constEnhancedMode, resource.EnhancedMode)
	d.Set(constDownloadAttempts, resource.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, resource.DownloadRetryBackoffSeconds)
	d.Set(constUsername, resource.Username)
	d.Set(constPassword, resource.Password)

	return nil
}

func buildNugetFeedResource(d *schema.ResourceData) *model.Feed {
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

	feed := model.NewNuGetFeed(feedName)
	feed.FeedURI = feedURI
	feed.EnhancedMode = enhancedMode
	feed.DownloadAttempts = downloadAttempts
	feed.DownloadRetryBackoffSeconds = downloadRetryBackoffSeconds
	feed.Username = feedUsername
	feed.Password = model.SensitiveValue{
		NewValue: &feedPassword,
	}

	return feed
}

func resourceNugetFeedCreate(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	newFeed := buildNugetFeedResource(d)
	feed, err := apiClient.Feeds.Add(*newFeed)

	if err != nil {
		return fmt.Errorf("error creating nuget feed %s: %s", newFeed.Name, err.Error())
	}

	d.SetId(feed.ID)

	return nil
}

func resourceNugetFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildNugetFeedResource(d)
	feed.ID = d.Id() // set project struct ID so octopus knows which project to update

	apiClient := m.(*client.Client)

	updatedFeed, err := apiClient.Feeds.Update(*feed)

	if err != nil {
		return fmt.Errorf("error updating nuget feed id %s: %s", d.Id(), err.Error())
	}

	d.SetId(updatedFeed.ID)
	return nil
}

func resourceNugetFeedDelete(d *schema.ResourceData, m interface{}) error {
	apiClient := m.(*client.Client)

	feedID := d.Id()

	err := apiClient.Feeds.DeleteByID(feedID)

	if err != nil {
		return fmt.Errorf("error deleting nuget feed id %s: %s", feedID, err.Error())
	}

	d.SetId(constEmptyString)
	return nil
}
