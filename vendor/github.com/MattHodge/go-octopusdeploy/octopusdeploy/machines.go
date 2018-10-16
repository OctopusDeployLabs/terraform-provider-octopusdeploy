package octopusdeploy

import (
	"fmt"
	"strings"

	"github.com/dghubble/sling"
)

type MachineService struct {
	sling *sling.Sling
}

func NewMachineService(sling *sling.Sling) *MachineService {
	return &MachineService{
		sling: sling,
	}
}

type Machines struct {
	Items []Machine `json:"Items"`
	PagedResults
}

type Machine struct {
	ID                              string           `json:"Id"`
	Name                            string           `json:"Name"`
	Thumbprint                      string           `json:"Thumbprint"`
	URI                             string           `json:"Uri"`
	IsDisabled                      bool             `json:"IsDisabled"`
	EnvironmentIDs                  []string         `json:"EnvironmentIds"`
	Roles                           []string         `json:"Roles"`
	MachinePolicyID                 string           `json:"MachinePolicyId"`
	TenantedDeploymentParticipation string           `json:"TenantedDeploymentParticipation"`
	TenantIDs                       []string         `json:"TenantIDs"`
	TenantTags                      []string         `json:"TenantTags"`
	Status                          string           `json:"Status"`
	HasLatestCalamari               bool             `json:"HasLatestCalamari"`
	StatusSummary                   string           `json:"StatusSummary"`
	IsInProcess                     bool             `json:"IsInProcess"`
	Endpoint                        *MachineEndpoint `json:"Endpoint,omitempty"`
	LastModifiedOn                  *string          `json:"LastModifiedOn,omitempty"`
	LastModifiedBy                  *string          `json:"LastModifiedBy,omitempty"`
}

type MachineEndpoint struct {
	ID                     string                        `json:"Id"`
	CommunicationStyle     string                        `json:"CommunicationStyle"`
	ProxyID                *string                       `json:"ProxyId"`
	Thumbprint             string                        `json:"Thumbprint"`
	TentacleVersionDetails MachineTentacleVersionDetails `json:"TentacleVersionDetails"`
	LastModifiedOn         *string                       `json:"LastModifiedOn,omitempty"`
	LastModifiedBy         *string                       `json:"LastModifiedBy,omitempty"`
	URI                    string                        `json:"Uri"` //This is not in the spec doc, but it shows up and needs to be kept in sync
}

type MachineTentacleVersionDetails struct {
	UpgradeLocked    bool   `json:"UpgradeLocked"`
	Version          string `json:"Version"`
	UpgradeSuggested bool   `json:"UpgradeSuggested"`
	UpgradeRequired  bool   `json:"UpgradeRequired"`
}

func NewMachine(Name string, Disabled bool, EnvironmentIDs []string, Roles []string, MachinePolicyId string, TenantedDeploymentParticipation string, TenantIDs, TenantTags []string) *Machine {
	return &Machine{
		Name:                            Name,
		IsDisabled:                      Disabled,
		EnvironmentIDs:                  EnvironmentIDs,
		Roles:                           Roles,
		MachinePolicyID:                 MachinePolicyId,
		TenantedDeploymentParticipation: TenantedDeploymentParticipation,
		TenantIDs:                       TenantIDs,
		TenantTags:                      TenantTags,
		Status:                          "Unknown",
		Thumbprint:                      "0123456789ABCDEF",
		URI:                             "https://localhost/",
	}
}

// ValidateMachineValues checks the values of a Machine object to see if they are suitable for
// sending to Octopus Deploy. Used when adding or updating machines.
func ValidateMachineValues(Machine *Machine) error {
	if Machine.Endpoint != nil {
		if Machine.Endpoint.URI != Machine.URI {
			return fmt.Errorf("URI and Endpoint.URI must match. They are currently %s and %s", Machine.URI, Machine.Endpoint.URI)
		}
		if Machine.Endpoint.Thumbprint != Machine.Thumbprint {
			return fmt.Errorf("Thumbprint and Endpoint.Thumbprint must match. They are currently %s and %s", Machine.Thumbprint, Machine.Endpoint.Thumbprint)
		}
	}

	// Check TenantedDeploymentParticipation
	ValidTenantedDeploymentModes := []string{
		"Untenanted", "TenantedOrUntenanted", "Tenanted",
	}
	var TenantedDeploymentParticipationOK bool
	for _, a := range ValidTenantedDeploymentModes {
		if a == Machine.TenantedDeploymentParticipation {
			TenantedDeploymentParticipationOK = true
		}
	}
	if !TenantedDeploymentParticipationOK {
		return fmt.Errorf("TenantedDeploymentParticipation must be one of \"%v\"", strings.Join(ValidTenantedDeploymentModes, ","))
	}

	// Check Status
	Statuses := []string{
		"Online", "Offline", "Unknown", "NeedsUpgrade", "CalamariNeedsUpgrade", "Disabled",
	}
	var StatusOK bool
	for _, a := range Statuses {
		if a == Machine.Status {
			StatusOK = true
		}
	}
	if !StatusOK {
		return fmt.Errorf("Status must be one of \"%v\"", strings.Join(Statuses, ","))
	}

	return nil
}

// Get returns a single machine with a given MachineID
func (s *MachineService) Get(MachineID string) (*Machine, error) {
	path := fmt.Sprintf("machines/%s", MachineID)
	resp, err := apiGet(s.sling, new(Machine), path)

	if err != nil {
		return nil, err
	}

	return resp.(*Machine), nil
}

// GetAll returns all registered machines
func (s *MachineService) GetAll() (*[]Machine, error) {
	var p []Machine
	path := "machines"
	loadNextPage := true

	for loadNextPage {
		resp, err := apiGet(s.sling, new(Machines), path)
		if err != nil {
			return nil, err
		}

		r := resp.(*Machines)
		for _, item := range r.Items {
			p = append(p, item)
		}

		path, loadNextPage = LoadNextPage(r.PagedResults)
	}
	return &p, nil
}

// Add creates a new machine in Octopus Deploy
func (s *MachineService) Add(machine *Machine) (*Machine, error) {
	err := ValidateMachineValues(machine)
	if err != nil {
		return nil, err
	}

	resp, err := apiAdd(s.sling, machine, new(Machine), "machines")
	if err != nil {
		return nil, err
	}

	return resp.(*Machine), nil
}

// Delete deletes an existing machine in Octopus Deploy
func (s *MachineService) Delete(MachineID string) error {
	path := fmt.Sprintf("machines/%s", MachineID)
	err := apiDelete(s.sling, path)

	if err != nil {
		return err
	}

	return nil
}

// Delete updates an existing machine in Octopus Deploy
func (s *MachineService) Update(machine *Machine) (*Machine, error) {
	err := ValidateMachineValues(machine)
	if err != nil {
		return nil, err
	}

	path := fmt.Sprintf("machines/%s", machine.ID)
	resp, err := apiUpdate(s.sling, machine, new(Machine), path)

	if err != nil {
		return nil, err
	}

	return resp.(*Machine), nil
}
