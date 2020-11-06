package octopusdeploy

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestTokenAccountBasic(t *testing.T) {
	localName := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	prefix := constOctopusDeployTokenAccount + "." + localName

	name := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	token := acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	resource.Test(t, resource.TestCase{
		CheckDestroy: testAccountDestroy,
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testTokenAccountBasic(localName, name, token),
				Check: resource.ComposeTestCheckFunc(
					testAccountExists(prefix),
					resource.TestCheckResourceAttr(prefix, constName, name),
					resource.TestCheckResourceAttr(prefix, constToken, token),
				),
			},
		},
	})
}

func testTokenAccountBasic(localName string, name string, token string) string {
	return fmt.Sprintf(`resource "%s" "%s" {
		name  = "%s"
		token = "%s"
	}`, constOctopusDeployTokenAccount, localName, name, token)
}
