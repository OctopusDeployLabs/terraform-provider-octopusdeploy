package octopusdeploy

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/certificates"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/octoclient"
	"github.com/OctopusSolutionsEngineering/OctopusTerraformTestFramework/test"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"path/filepath"
)

func (suite *IntegrationTestSuite) TestAccOctopusDeployCertificateBasic() {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_certificate." + localName

	certificateData := "MIIDiDCCAnACCQDXHofnqz05ITANBgkqhkiG9w0BAQsFADCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wHhcNMTkwNjE0MjExMzI1WhcNMjAwNjEzMjExMzI1WjCBhTELMAkGA1UEBhMCVVMxETAPBgNVBAgMCE9rbGFob21hMQ8wDQYDVQQHDAZOb3JtYW4xEzARBgNVBAoMCk1vb25zd2l0Y2gxGTAXBgNVBAMMEGRlbW8ub2N0b3B1cy5jb20xIjAgBgkqhkiG9w0BCQEWE2plZmZAbW9vbnN3aXRjaC5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDSTiD0OHyFDMH9O+d/h3AiqcuvpvUgRkKjf+whZ6mVlQnGkvPddRTUY48xCEaQ4QD1MAVJcGaJ2PU4NxwhrQgHqWW8TQkAZESL4wfzSwIKO2NX/I2tWqyv7a0uA/WdtlWQye+2oPV5rCnS0kM75X+gjEwOTpFh/ryS6KhMPFDb0zeNGREdg6564FdxWSvN4ppUZMqhvMpfzM7rsDWqEzYsMaQ4CNJDFdWkG89D4j5qk4b4Qb4m+l7QINdmYIXf4qO/0LE1WcfIkCpAS65tjc/hefIHmYtj/E/ijoNJbWKZDK3WLZg3zq99Ipqv/9DFvSiMQFBhZT0jO2B5d5zBUuIHAgMBAAEwDQYJKoZIhvcNAQELBQADggEBAKsa4gaW7GhByu8aq56h99DaIl1LauI5WMVH8Q9Qpapho2VLRIpfwGeI5eENFoXwuKrnJp1ncsCqGnMQnugQHS+SrruS3Yyl0Uog4Zak9GbbK6qn+olx7GNJbsckmD371lqQOaKITLqYzK6kTc7/v8Cv0BwHFCBda1OCrmeVBSaarucPxZhGxzLAielzHHdlkZFQT/oO2VR3thhURIqtni7jVQ2MoeZF1ccvmAfVbzr/QnlNe/jrcmyPYymuF2JyrezzIjlKuiDhalKqwqkCHpOOgzV4y6BFuS+0w3DS8pa07nUudZ6E0kZzvhjjiyAx/sBdX6ZDdUjP9TDJMM4f5YA="
	name := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	password := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)

	resource.Test(suite.T(), resource.TestCase{
		CheckDestroy:             testAccCertificateCheckDestroy,
		PreCheck:                 func() { testAccPreCheck(suite.T()) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Check: resource.ComposeTestCheckFunc(
					testCertificateExists(prefix),
					resource.TestCheckResourceAttr(prefix, "certificate_data", certificateData),
					resource.TestCheckResourceAttr(prefix, "name", name),
					resource.TestCheckResourceAttr(prefix, "password", password),
				),
				Config: testCertificateBasic(localName, name, certificateData, password),
			},
		},
	})
}

func testCertificateBasic(localName string, name string, certificateData string, password string) string {
	return fmt.Sprintf(`resource "octopusdeploy_certificate" "%s" {
		certificate_data = "%s"
		name             = "%s"
		password         = "%s"
	}`, localName, certificateData, name, password)
}

func testCertificateExists(prefix string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		certificateID := s.RootModule().Resources[prefix].Primary.ID
		if _, err := octoClient.Certificates.GetByID(certificateID); err != nil {
			return err
		}

		return nil
	}
}

func testAccCertificateCheckDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_certficate" {
			continue
		}

		certificate, err := octoClient.Certificates.GetByID(rs.Primary.ID)
		if err == nil && certificate != nil {
			return fmt.Errorf("certificate (%s) still exists", rs.Primary.ID)
		}
	}

	return nil
}

// TestCertificateResource verifies that a certificate can be reimported with the correct settings
func (suite *IntegrationTestSuite) TestCertificateResource() {
	testFramework := test.OctopusContainerTest{}
	t := suite.T()

	newSpaceId, err := testFramework.Act(t, octoContainer, "../terraform", "25-certificates", []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	err = testFramework.TerraformInitAndApply(t, octoContainer, filepath.Join("../terraform", "25a-certificatesds"), newSpaceId, []string{})

	if err != nil {
		t.Fatal(err.Error())
	}

	// Assert
	client, err := octoclient.CreateClient(octoContainer.URI, newSpaceId, test.ApiKey)
	query := certificates.CertificatesQuery{
		PartialName: "Test",
		Skip:        0,
		Take:        1,
	}

	resources, err := client.Certificates.Get(query)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(resources.Items) == 0 {
		t.Fatalf("Space must have a certificate called \"Test\"")
	}
	resource := resources.Items[0]

	if resource.Notes != "A test certificate" {
		t.Fatal("The tenant must be have a description of \"A test certificate\" (was \"" + resource.Notes + "\")")
	}

	if resource.TenantedDeploymentMode != "Untenanted" {
		t.Fatal("The tenant must be have a tenant participation of \"Untenanted\" (was \"" + resource.TenantedDeploymentMode + "\")")
	}

	if resource.SubjectDistinguishedName != "CN=test.com" {
		t.Fatal("The tenant must be have a subject distinguished name of \"CN=test.com\" (was \"" + resource.SubjectDistinguishedName + "\")")
	}

	if len(resource.EnvironmentIDs) != 0 {
		t.Fatal("The tenant must have one project environment")
	}

	if len(resource.TenantTags) != 0 {
		t.Fatal("The tenant must have no tenant tags")
	}

	if len(resource.TenantIDs) != 0 {
		t.Fatal("The tenant must have no tenants")
	}

	// Verify the environment data lookups work
	lookup, err := testFramework.GetOutputVariable(t, filepath.Join("..", "terraform", "25a-certificatesds"), "data_lookup")

	if err != nil {
		t.Fatal(err.Error())
	}

	if lookup != resource.ID {
		t.Fatal("The environment lookup did not succeed. Lookup value was \"" + lookup + "\" while the resource value was \"" + resource.ID + "\".")
	}
}
