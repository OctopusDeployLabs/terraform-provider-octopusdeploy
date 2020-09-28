package octopusdeploy

import (
	"log"

	"github.com/go-playground/validator"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func validateSchema() {
	val := validator.New()
	Schema := &schema.Resource{}
	validate := val.Struct(Schema)

	if validate != nil {
		log.Println("Ensure that the schema map is valid: https://www.terraform.io/docs/extend/schemas/schema-types.html")
	}
}
