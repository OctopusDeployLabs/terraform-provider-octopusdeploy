package octopusdeploy

import (
	"fmt"

	"github.com/dghubble/sling"
	"gopkg.in/go-playground/validator.v9"
)

type ProjectService struct {
	sling *sling.Sling
}

func NewProjectService(sling *sling.Sling) *ProjectService {
	return &ProjectService{
		sling: sling,
	}
}

type Projects struct {
	Items []Project `json:"Items"`
	PagedResults
}

type Project struct {
	AutoCreateRelease               bool                        `json:"AutoCreateRelease"`
	AutoDeployReleaseOverrides      []AutoDeployReleaseOverride `json:"AutoDeployReleaseOverrides"`
	DefaultGuidedFailureMode        string                      `json:"DefaultGuidedFailureMode,omitempty" validate:"oneof=EnvironmentDefault Off On"`
	DefaultToSkipIfAlreadyInstalled bool                        `json:"DefaultToSkipIfAlreadyInstalled"`
	DeploymentProcessID             string                      `json:"DeploymentProcessId"`
	Description                     string                      `json:"Description"`
	DiscreteChannelRelease          bool                        `json:"DiscreteChannelRelease"`
	ID                              string                      `json:"Id,omitempty"`
	IncludedLibraryVariableSetIds   []string                    `json:"IncludedLibraryVariableSetIds"`
	IsDisabled                      bool                        `json:"IsDisabled"`
	LifecycleID                     string                      `json:"LifecycleId" validate:"required"`
	Name                            string                      `json:"Name" validate:"required"`
	ProjectConnectivityPolicy       ProjectConnectivityPolicy   `json:"ProjectConnectivityPolicy"`
	ProjectGroupID                  string                      `json:"ProjectGroupId" validate:"required"`
	ReleaseCreationStrategy         ReleaseCreationStrategy     `json:"ReleaseCreationStrategy"`
	Slug                            string                      `json:"Slug"`
	Templates                       []ActionTemplateParameter   `json:"Templates,omitempty"`
	TenantedDeploymentMode          string                      `json:"TenantedDeploymentMode,omitempty"`
	VariableSetID                   string                      `json:"VariableSetId"`
	VersioningStrategy              VersioningStrategy          `json:"VersioningStrategy"`
}

func (t *Project) Validate() error {
	validate := validator.New()

	err := validate.Struct(t)

	if err != nil {
		return err
	}

	return nil
}

func NewProject(name, lifeCycleID, projectGroupID string) *Project {
	return &Project{
		Name: name,
		DefaultGuidedFailureMode: "EnvironmentDefault",
		LifecycleID:              lifeCycleID,
		ProjectGroupID:           projectGroupID,
		VersioningStrategy: VersioningStrategy{
			Template: "#{Octopus.Version.LastMajor}.#{Octopus.Version.LastMinor}.#{Octopus.Version.NextPatch}",
		},
		ProjectConnectivityPolicy: ProjectConnectivityPolicy{
			AllowDeploymentsToNoTargets: false,
			SkipMachineBehavior:         "None",
		},
	}
}

func (s *ProjectService) Get(projectid string) (*Project, error) {
	path := fmt.Sprintf("projects/%s", projectid)
	resp, err := apiGet(s.sling, new(Project), path)

	if err != nil {
		return nil, err
	}

	return resp.(*Project), nil
}

func (s *ProjectService) GetAll() (*[]Project, error) {
	var p []Project

	path := "projects"

	loadNextPage := true

	for loadNextPage {
		resp, err := apiGet(s.sling, new(Projects), path)

		if err != nil {
			return nil, err
		}

		r := resp.(*Projects)

		for _, item := range r.Items {
			p = append(p, item)
		}

		path, loadNextPage = LoadNextPage(r.PagedResults)
	}

	return &p, nil
}

func (s *ProjectService) GetByName(projectName string) (*Project, error) {
	var foundProject Project
	projects, err := s.GetAll()

	if err != nil {
		return nil, err
	}

	for _, project := range *projects {
		if project.Name == projectName {
			return &project, nil
		}
	}

	return &foundProject, fmt.Errorf("no project found with project name %s", projectName)
}

func (s *ProjectService) Add(project *Project) (*Project, error) {
	resp, err := apiAdd(s.sling, project, new(Project), "projects")

	if err != nil {
		return nil, err
	}

	return resp.(*Project), nil
}

func (s *ProjectService) Delete(projectid string) error {
	path := fmt.Sprintf("projects/%s", projectid)
	err := apiDelete(s.sling, path)

	if err != nil {
		return err
	}

	return nil
}

func (s *ProjectService) Update(project *Project) (*Project, error) {
	path := fmt.Sprintf("projects/%s", project.ID)
	resp, err := apiUpdate(s.sling, project, new(Project), path)

	if err != nil {
		return nil, err
	}

	return resp.(*Project), nil
}
