package internal

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/newclient"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/projects"
)

func CheckRunbookInGit(client newclient.Client, spaceID string, projectID string) (bool, error) {
	// We need both to be able to check the git runbook setting
	// If project ID is not set, assume the runbook is not in git to maintain backwards compatibility
	if projectID == "" {
		return false, nil
	}

	project, err := projects.GetByID(client, spaceID, projectID)

	if err != nil {
		return false, err
	}

	if project.PersistenceSettings != nil && project.PersistenceSettings.Type() == projects.PersistenceSettingsTypeVersionControlled && project.PersistenceSettings.(projects.GitPersistenceSettings).RunbooksAreInGit() {
		return true, nil
	}

	return false, nil
}
