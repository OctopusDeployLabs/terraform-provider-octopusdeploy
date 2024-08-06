package octopusdeploy_framework

import (
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"path/filepath"
	"testing"
)

func TestTentacleCertificateResource(t *testing.T) {
	testFramework := test.OctopusContainerTest{}
	_, err := testFramework.Act(t, octoContainer, "../terraform", "57-tentaclecertificate", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	thumbprintLookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "57-tentaclecertificate"), "base_certificate_thumbprint")
	if err != nil {
		t.Fatal(err.Error())
	}

	if thumbprintLookup == "" {
		t.Fatalf("Expected a thumbprint to be returned in Terraform output")
	}

}
