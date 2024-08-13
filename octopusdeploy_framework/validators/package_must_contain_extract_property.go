package validators

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.Map = packageMustContainExtractProperty{}

type packageMustContainExtractProperty struct {
}

func (p packageMustContainExtractProperty) Description(_ context.Context) string {
	return "Package property collection must contain the Extract "
}

func (p packageMustContainExtractProperty) MarkdownDescription(ctx context.Context) string {
	return p.Description(ctx)
}

func (p packageMustContainExtractProperty) ValidateMap(_ context.Context, req validator.MapRequest, resp *validator.MapResponse) {
	current := req.ConfigValue.Elements()
	if _, ok := current["Extract"]; !ok {
		resp.Diagnostics.Append(validatordiag.InvalidBlockDiagnostic(req.Path, "package properties must contain Extract"))
	}
}

func PackageMustContainExtractProperty() validator.Map {
	return packageMustContainExtractProperty{}
}
