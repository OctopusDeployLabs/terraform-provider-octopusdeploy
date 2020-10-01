package octopusdeploy

import (
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
)

func diagValidate() diag.Diagnostics {
	var diags diag.Diagnostics

	if diags == nil {
		log.Println("diag package is empty")
	}

	return diags
}
