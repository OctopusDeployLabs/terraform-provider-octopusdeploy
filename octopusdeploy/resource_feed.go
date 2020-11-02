package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	feed, err := client.Feeds.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingFeed, id, err)
	}
	if feed == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constFeed, m)

	feedResource := feed.(*octopusdeploy.FeedResource)

	d.Set(constName, feedResource.Name)
	d.Set(constFeedType, feedResource.FeedType)
	d.Set(constFeedURI, feedResource.FeedURI)
	d.Set(constEnhancedMode, feedResource.EnhancedMode)
	d.Set(constDownloadAttempts, feedResource.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, feedResource.DownloadRetryBackoffSeconds)
	d.Set(constUsername, feedResource.Username)
	d.Set(constPassword, feedResource.Password)

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

func resourceFeedCreate(d *schema.ResourceData, m interface{}) error {
	feedResource := buildFeedResource(d)
	client := m.(*octopusdeploy.Client)

	feed, err := client.Feeds.Add(feedResource)
	if err != nil {
		return createResourceOperationError(errorCreatingFeed, feedResource.Name, err)
	}

	if isEmpty(feed.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(feed.GetID())
	}

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feedResource := buildFeedResource(d)
	client := m.(*octopusdeploy.Client)

	feedResource.ID = d.Id() // set ID so Octopus API knows which feed to update

	feed, err := client.Feeds.Update(feedResource)
	if err != nil {
		return fmt.Errorf(errorUpdatingFeed, d.Id(), err.Error())
	}

	d.SetId(feed.GetID())

	return nil
}

func resourceFeedDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingFeed, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
