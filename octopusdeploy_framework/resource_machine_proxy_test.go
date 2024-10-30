package octopusdeploy_framework

import (
	"fmt"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/core"
	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/proxies"
	"github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"testing"
)

func TestAccMachineProxyBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(20, acctest.CharSetAlpha)
	prefix := "octopusdeploy_machine_proxy." + localName

	data := &proxies.Proxy{
		Name:      "Test Proxy",
		Host:      "127.0.0.1",
		Port:      8080,
		Username:  "admin",
		Password:  core.NewSensitiveValue("safepass"),
		ProxyType: "HTTP",
	}

	resource.Test(t, resource.TestCase{
		CheckDestroy:             testMachineProxyDestroy,
		PreCheck:                 func() { TestAccPreCheck(t) },
		ProtoV6ProviderFactories: ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testMachineProxyBasic(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.Name),
				),
			},
			{
				Config: testMachineProxyUpdate(data),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(prefix, "name", data.Name+"-updated"),
				),
			},
		},
	})
}

func testMachineProxyBasic(data *proxies.Proxy) string {
	return fmt.Sprintf(`
	resource "octopusdeploy_machine_proxy" "test_proxy" {
		name = %s
		host = %s
		username = %s
		password = %s
		port = %d
	}
`,
		data.Name,
		data.Host,
		data.Username,
		&data.Password.NewValue,
		data.Port,
	)
}

func testMachineProxyUpdate(data *proxies.Proxy) string {
	data.Name = data.Name + "-updated"

	return testMachineProxyBasic(data)
}

func testMachineProxyDestroy(s *terraform.State) error {
	var machineProxyID string

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "octopusdeploy_machine_proxy" {
			continue
		}

		machineProxyID = rs.Primary.ID
		break
	}
	if machineProxyID == "" {
		return fmt.Errorf("no octopusdeploy_machine_proxy resource found")
	}

	machineProxy, err := proxies.GetByID(octoClient, octoClient.GetSpaceID(), machineProxyID)
	if err == nil {
		if machineProxy != nil {
			return fmt.Errorf("machine proxy (%s) still exists", machineProxy.Name)
		}
	}

	return nil
}
