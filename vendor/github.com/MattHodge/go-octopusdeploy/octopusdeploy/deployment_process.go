package octopusdeploy

import (
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type DeploymentProcessService struct {
	sling *sling.Sling
}

func NewDeploymentProcessService(sling *sling.Sling) *DeploymentProcessService {
	return &DeploymentProcessService{
		sling: sling,
	}
}

type DeploymentProcesses struct {
	Items []DeploymentProcess `json:"Items"`
	PagedResults
}

type DeploymentProcess struct {
	ID             string                   `json:"Id"`
	ProjectID      string                   `json:"ProjectId"`
	Steps          []DeploymentStepResource `json:"Steps"`
	Version        int                      `json:"Version"`
	LastSnapshotID string                   `json:"LastSnapshotId"`
	LastModifiedOn string                   `json:"LastModifiedOn"` // date time
	LastModifiedBy string                   `json:"LastModifiedBy"`
	Links          Links                    `json:"Links"`
}

func (d *DeploymentProcessService) Get(deploymentProcessID string) (DeploymentProcess, error) {
	deploymentProcess := new(DeploymentProcess)
	octopusDeployError := new(APIError)
	path := fmt.Sprintf("deploymentprocesses/%s", deploymentProcessID)

	resp, err := d.sling.New().Get(path).Receive(deploymentProcess, octopusDeployError)

	if err != nil {
		return *deploymentProcess, fmt.Errorf("cannot get deploymentprocess id %s from server. failure from http client %v", deploymentProcessID, err)
	}

	if resp.StatusCode != http.StatusOK {
		return *deploymentProcess, fmt.Errorf("cannot get deploymentprocess id %s from server. response from server %s", deploymentProcessID, resp.Status)
	}

	return *deploymentProcess, err
}

func (d *DeploymentProcessService) GetAll() ([]DeploymentProcess, error) {
	var listOfDeploymentProcess []DeploymentProcess
	path := fmt.Sprintf("deploymentprocesses")

	for {
		deploymentProcesses := new(DeploymentProcesses)
		octopusDeployError := new(APIError)

		resp, err := d.sling.New().Get(path).Receive(deploymentProcesses, octopusDeployError)
		if err != nil {
			return nil, err
		}

		fmt.Printf("Response: %v", resp.Status)
		fmt.Printf("Total Results: %d", deploymentProcesses.NumberOfPages)

		for _, deploymentProcess := range deploymentProcesses.Items {
			listOfDeploymentProcess = append(listOfDeploymentProcess, deploymentProcess)
		}

		if deploymentProcesses.PagedResults.Links.PageNext != "" {
			fmt.Printf("More pages to go! Next link: %s", deploymentProcesses.PagedResults.Links.PageNext)
			path = deploymentProcesses.PagedResults.Links.PageNext
		} else {
			break
		}
	}

	return listOfDeploymentProcess, nil // no more pages to go through
}
