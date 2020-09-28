package octopusdeploy

import (
	"fmt"
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/client"
	"github.com/OctopusDeploy/go-octopusdeploy/enum"
	"github.com/OctopusDeploy/go-octopusdeploy/model"
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

	apiClient := m.(*client.Client)
	resource, err := apiClient.Feeds.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingFeed, id, err)
	}
	if resource == nil {
		d.SetId(constEmptyString)
		return nil
	}

	logResource(constFeed, m)

	d.Set(constName, resource.Name)
	d.Set(constFeedType, resource.FeedType)
	d.Set(constFeedURI, resource.FeedURI)
	d.Set(constEnhancedMode, resource.EnhancedMode)
	d.Set(constDownloadAttempts, resource.DownloadAttempts)
	d.Set(constDownloadRetryBackoffSeconds, resource.DownloadRetryBackoffSeconds)
	d.Set(constUsername, resource.Username)
	d.Set(constPassword, resource.Password)

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
	feed := buildFeedResource(d)

	apiClient := m.(*client.Client)
	resource, err := apiClient.Feeds.Add(*feed)
	if err != nil {
		return createResourceOperationError(errorCreatingFeed, feed.Name, err)
	}

	if isEmpty(resource.ID) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.ID)
	}

	return nil
}

func resourceFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildFeedResource(d)
	feed.ID = d.Id() // set ID so Octopus API knows which feed to update

	apiClient := m.(*client.Client)
	resource, err := apiClient.Feeds.Update(*feed)
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
		return createResourceOperationError(errorDeletingFeed, id, err)
	}

	d.SetId(constEmptyString)
	return nil
}
