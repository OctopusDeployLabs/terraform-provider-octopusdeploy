package octopusdeploy

import (
	"fmt"

	"github.com/dghubble/sling"
	"gopkg.in/go-playground/validator.v9"
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
	ID             string           `json:"Id,omitempty"`
	LastModifiedBy string           `json:"LastModifiedBy,omitempty"`
	LastModifiedOn string           `json:"LastModifiedOn,omitempty"`
	LastSnapshotID string           `json:"LastSnapshotId,omitempty"`
	Links          Links            `json:"Links,omitempty"`
	ProjectID      string           `json:"ProjectId,omitempty"`
	Steps          []DeploymentStep `json:"Steps"`
	Version        *int32           `json:"Version"`
}

type DeploymentStep struct {
	ID                 string             `json:"Id"`
	Name               string             `json:"Name"`
	PackageRequirement string             `json:"PackageRequirement,omitempty"`                                         // may need its own model / enum
	Properties         map[string]string  `json:"Properties"`                                                           // TODO: refactor to use the PropertyValueResource for handling sensitive values - https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
	Condition          string             `json:"Condition,omitempty" validate:"oneof=Success Failure Always Variable"` // variable option adds a Property "Octopus.Action.ConditionVariableExpression"
	StartTrigger       string             `json:"StartTrigger,omitempty" validate:"oneof=StartAfterPrevious StartWithPrevious"`
	Actions            []DeploymentAction `json:"Actions"`
}

type DeploymentAction struct {
	ID                            string            `json:"Id"`
	Name                          string            `json:"Name"`
	ActionType                    string            `json:"ActionType"`
	IsDisabled                    bool              `json:"IsDisabled"`
	CanBeUsedForProjectVersioning bool              `json:"CanBeUsedForProjectVersioning"`
	Environments                  []string          `json:"Environments"`
	ExcludedEnvironments          []string          `json:"ExcludedEnvironments"`
	Channels                      []string          `json:"Channels"`
	TenantTags                    []string          `json:"TenantTags"`
	Properties                    map[string]string `json:"Properties"`     // TODO: refactor to use the PropertyValueResource for handling sensitive values - https://blog.gopheracademy.com/advent-2016/advanced-encoding-decoding/
	LastModifiedOn                string            `json:"LastModifiedOn"` // datetime
	LastModifiedBy                string            `json:"LastModifiedBy"`
	Links                         Links             `json:"Links"` // may be wrong
}

func (d *DeploymentProcess) Validate() error {
	validate := validator.New()

	err := validate.Struct(d)

	if err != nil {
		return err
	}

	return nil
}

func (s *DeploymentProcessService) Get(deploymentProcessID string) (*DeploymentProcess, error) {
	path := fmt.Sprintf("deploymentprocesses/%s", deploymentProcessID)
	resp, err := apiGet(s.sling, new(DeploymentProcess), path)

	if err != nil {
		return nil, err
	}

	return resp.(*DeploymentProcess), nil
}

func (s *DeploymentProcessService) GetAll() (*[]DeploymentProcess, error) {
	var dp []DeploymentProcess

	path := "deploymentprocesses"

	loadNextPage := true

	for loadNextPage {
		resp, err := apiGet(s.sling, new(DeploymentProcesses), path)

		if err != nil {
			return nil, err
		}

		r := resp.(*DeploymentProcesses)

		for _, item := range r.Items {
			dp = append(dp, item)
		}

		path, loadNextPage = LoadNextPage(r.PagedResults)
	}

	return &dp, nil
}

func (s *DeploymentProcessService) Update(deploymentProcess *DeploymentProcess) (*DeploymentProcess, error) {
	path := fmt.Sprintf("deploymentprocesses/%s", deploymentProcess.ID)
	resp, err := apiUpdate(s.sling, deploymentProcess, new(DeploymentProcess), path)

	if err != nil {
		return nil, err
	}

	return resp.(*DeploymentProcess), nil
}
