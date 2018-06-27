package octopusdeploy

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dghubble/sling"
)

type ProjectsService struct {
	sling *sling.Sling
}

func NewProjectService(sling *sling.Sling) *ProjectsService {
	return &ProjectsService{
		sling: sling,
	}
}

type Projects struct {
	Items []Project `json:"Items"`
	PagedResults
}

type Project struct {
	AutoCreateRelease               bool                                `json:"AutoCreateRelease"`
	AutoDeployReleaseOverrides      []AutoDeployReleaseOverrideResource `json:"AutoDeployReleaseOverrides"`
	DefaultGuidedFailureMode        string                              `json:"DefaultGuidedFailureMode,omitempty"`
	DefaultToSkipIfAlreadyInstalled bool                                `json:"DefaultToSkipIfAlreadyInstalled"`
	DeploymentProcessID             string                              `json:"DeploymentProcessId"`
	Description                     string                              `json:"Description"`
	DiscreteChannelRelease          bool                                `json:"DiscreteChannelRelease"`
	ID                              string                              `json:"Id,omitempty"`
	IncludedLibraryVariableSetIds   []string                            `json:"IncludedLibraryVariableSetIds"`
	IsDisabled                      bool                                `json:"IsDisabled"`
	LifecycleID                     string                              `json:"LifecycleId"`
	Name                            string                              `json:"Name"`
	ProjectConnectivityPolicy       ProjectConnectivityPolicy           `json:"ProjectConnectivityPolicy"`
	ProjectGroupID                  string                              `json:"ProjectGroupId"`
	ReleaseCreationStrategy         ReleaseCreationStrategyResource     `json:"ReleaseCreationStrategy"`
	Slug                            string                              `json:"Slug"`
	Templates                       []ActionTemplateParameterResource   `json:"Templates,omitempty"`
	TenantedDeploymentMode          string                              `json:"TenantedDeploymentMode,omitempty"`
	VariableSetID                   string                              `json:"VariableSetId"`
	VersioningStrategy              VersioningStrategyResource          `json:"VersioningStrategy"`
}

func NewProject(name, lifeCycleID, projectGroupID string) *Project {
	return &Project{
		Name:           name,
		LifecycleID:    lifeCycleID,
		ProjectGroupID: projectGroupID,
		VersioningStrategy: VersioningStrategyResource{
			Template: "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}",
		},
	}
}

func (s *ProjectsService) Get(projectid string) (*Project, error) {
	var project Project
	octopusDeployError := new(APIError)
	path := fmt.Sprintf("projects/%s", projectid)

	resp, err := s.sling.New().Get(path).Receive(&project, &octopusDeployError)

	if err != nil {
		return nil, fmt.Errorf("cannot get project id %s from server. failure from http client %v", projectid, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrItemNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cannot get project id %s from server. response from server %s", projectid, resp.Status)
	}

	return &project, err
}

func (s *ProjectsService) GetAll() (*[]Project, error) {
	var listOfProjects []Project
	path := fmt.Sprintf("projects")

	for {
		var projects Projects
		var octopusDeployError APIError

		resp, err := s.sling.New().Get(path).Receive(&projects, &octopusDeployError)

		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if octopusDeployError.Errors != nil {
			return nil, fmt.Errorf("cannot add project. response from octopus deploy %s: ", octopusDeployError.Errors)
		}

		if resp.StatusCode != http.StatusOK {
			return &listOfProjects, fmt.Errorf("cannot get all projects. response from server %s", resp.Status)
		}

		for _, project := range projects.Items {
			listOfProjects = append(listOfProjects, project)
		}

		if projects.PagedResults.Links.PageNext != "" {
			path = projects.PagedResults.Links.PageNext
		} else {
			break
		}
	}

	return &listOfProjects, nil // no more pages to go through
}

func (s *ProjectsService) GetByName(projectName string) (*Project, error) {
	var foundProject Project
	projects, err := s.GetAll()

	if err != nil {
		return &foundProject, err
	}

	for _, project := range *projects {
		if project.Name == projectName {
			return &project, nil
		}
	}

	return &foundProject, fmt.Errorf("no project found with project name %s", projectName)
}

func (s *ProjectsService) Add(project *Project) (*Project, error) {
	var created Project
	var octopusDeployError APIError
	resp, err := s.sling.New().Post("projects").BodyJSON(project).Receive(&created, &octopusDeployError)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if octopusDeployError.Errors != nil {
		return nil, fmt.Errorf("cannot add project. response from octopus deploy %s: ", octopusDeployError.Errors)
	}

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("cannot add project. response from server %s, req %s", resp.Status, resp.Request.URL)
	}

	return &created, nil
}

func (s *ProjectsService) Delete(projectid string) error {
	path := fmt.Sprintf("projects/%s", projectid)
	req, err := s.sling.New().Delete(path).Request()

	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return ErrItemNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot delete project. response from server %s", resp.Status)
	}

	return nil
}

func (s *ProjectsService) Update(project Project) (Project, error) {
	var updated Project
	path := fmt.Sprintf("projects/%s", project.ID)
	resp, err := s.sling.New().Put(path).BodyJSON(project).ReceiveSuccess(&updated)

	if err != nil {
		return updated, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return updated, fmt.Errorf("cannot update project at url %s. response from server %s", resp.Request.URL, resp.Status)
	}

	return updated, nil
}

var ErrItemNotFound = errors.New("cannot find the item")
