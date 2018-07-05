package octopusdeploy

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/dghubble/sling"
)

// Client is a Twitter client for making Twitter API requests.
type Client struct {
	sling *sling.Sling
	// Octopus Deploy API Services
	DeploymentProcess *DeploymentProcessService
	ProjectGroup      *ProjectGroupService
	Project          *ProjectService
}

// NewClient returns a new Client.
func NewClient(httpClient *http.Client, octopusURL, octopusAPIKey string) *Client {
	baseURLWithAPI := strings.TrimRight(octopusURL, "/")
	baseURLWithAPI = fmt.Sprintf("%s/api/", baseURLWithAPI)
	fmt.Println(baseURLWithAPI)
	base := sling.New().Client(httpClient).Base(baseURLWithAPI).Set("X-Octopus-ApiKey", octopusAPIKey)
	return &Client{
		sling:             base,
		DeploymentProcess: NewDeploymentProcessService(base.New()),
		ProjectGroup:      NewProjectGroupService(base.New()),
		Project:          NewProjectService(base.New()),
	}
}

type APIError struct {
	ErrorMessage  string   `json:"ErrorMessage"`
	Errors        []string `json:"Errors"`
	FullException string   `json:"FullException"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Octopus Deploy Error Response: %v %+v %v", e.ErrorMessage, e.Errors, e.FullException)
}

func APIErrorChecker(urlPath string, resp *http.Response, wantedResponseCode int, slingError error, octopusDeployError *APIError) error {
	if octopusDeployError.Errors != nil {
		return fmt.Errorf("cannot get all projects. response from octopusdeploy %s: ", octopusDeployError.Errors)
	}

	if slingError != nil {
		return fmt.Errorf("cannot get path %s from server. failure from http client %v", urlPath, slingError)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrItemNotFound
	}

	if resp.StatusCode != wantedResponseCode {
		return fmt.Errorf("cannot get item from %s from server. response from server %s", urlPath, resp.Status)
	}

	return nil
}

var ErrItemNotFound = errors.New("cannot find the item")
