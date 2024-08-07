package validators

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"
)

var _ validator.String = notWhiteSpace{}

type notWhiteSpace struct {
}

func (n notWhiteSpace) Description(ctx context.Context) string {
	return "Value must not be whitespace"
}

func (n notWhiteSpace) MarkdownDescription(ctx context.Context) string {
	return n.Description(ctx)
}

func (n notWhiteSpace) ValidateString(ctx context.Context, request validator.StringRequest, response *validator.StringResponse) {
	if strings.TrimSpace(request.ConfigValue.ValueString()) == "" {
		response.Diagnostics.Append(validatordiag.InvalidAttributeValueDiagnostic(request.Path, "must not be empty or whitespace", request.ConfigValue.ValueString()))
	}
}

func NotWhitespace() validator.String {
	return notWhiteSpace{}
}
