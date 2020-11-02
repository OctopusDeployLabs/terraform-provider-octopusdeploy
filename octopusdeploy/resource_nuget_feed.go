package octopusdeploy

import (
	"log"

	"github.com/OctopusDeploy/go-octopusdeploy/octopusdeploy"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	feedResource, err := client.Feeds.GetByID(id)
	if err != nil {
		return createResourceOperationError(errorReadingFeed, id, err)
	}
	if feedResource == nil {
		d.SetId(constEmptyString)
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

func buildNugetFeedResource(d *schema.ResourceData) *octopusdeploy.NuGetFeed {
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

func resourceNugetFeedCreate(d *schema.ResourceData, m interface{}) error {
	feed := buildNugetFeedResource(d)

	client := m.(*octopusdeploy.Client)
	resource, err := client.Feeds.Add(feed)
	if err != nil {
		return createResourceOperationError(errorCreatingNuGetFeed, feed.Name, err)
	}

	if isEmpty(resource.GetID()) {
		log.Println("ID is nil")
	} else {
		d.SetId(resource.GetID())
	}

	return nil
}

func resourceNugetFeedUpdate(d *schema.ResourceData, m interface{}) error {
	feed := buildNugetFeedResource(d)
	feed.ID = d.Id() // set ID so Octopus API knows which feed to update

	client := m.(*octopusdeploy.Client)
	updatedFeed, err := client.Feeds.Update(feed)
	if err != nil {
		return createResourceOperationError(errorUpdatingNuGetFeed, d.Id(), err)
	}

	d.SetId(updatedFeed.GetID())

	return nil
}

func resourceNugetFeedDelete(d *schema.ResourceData, m interface{}) error {
	id := d.Id()

	client := m.(*octopusdeploy.Client)
	err := client.Feeds.DeleteByID(id)
	if err != nil {
		return createResourceOperationError(errorDeletingNuGetFeed, id, err)
	}

	d.SetId(constEmptyString)

	return nil
}
