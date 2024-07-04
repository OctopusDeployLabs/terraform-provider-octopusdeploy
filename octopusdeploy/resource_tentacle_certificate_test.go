package octopusdeploy

import (
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/client"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestTentacleCertificateResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	testFramework.ArrangeTest(t, func(t *testing.T, container *test.OctopusContainer, spaceClient *client.Client) error {
		// Act
		_, err := testFramework.Act(t, container, "../terraform", "57-tentaclecertificate", []string{})

		if err != nil {
			return err
		}

		thumbprintLookup, err := testFramework.GetOutputVariable(t, filepath.Join("terraform", "57-tentaclecertificate"), "base_certificate_thumbprint")
		if err != nil {
			return err
		}

		if thumbprintLookup == "" {
			t.Fatalf("Expected a thumbprint to be returned in Terraform output")
		}

		return nil
	})
}
