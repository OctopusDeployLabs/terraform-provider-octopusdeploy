package validators

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = actionTypeHasImplementationValidator{}

type actionTypeHasImplementationValidator struct {
}

func (a actionTypeHasImplementationValidator) Description(ctx context.Context) string {
	return "This action type has its own explicit action type"
}

func (a actionTypeHasImplementationValidator) MarkdownDescription(ctx context.Context) string {
	return a.Description(ctx)
}

func (a actionTypeHasImplementationValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	var actionType string
	switch req.ConfigValue.String() {
	case "Octopus.KubernetesDeploySecret":
		actionType = "deploy_kubernetes_secret_action"
	case "Octopus.KubernetesRunScript":
		actionType = "run_kubectl_script_action"
	case "Octopus.Manual":
		actionType = "manual_intervention_action"
	case "Octopus.Script":
		actionType = "run_script_action"
	case "Octopus.TentaclePackage":
		actionType = "deploy_package_action"
	case "Octopus.TerraformApply":
		actionType = "apply_terraform_template_action"
	case "Octopus.WindowsService":
		actionType = "deploy_windows_service_action"
	}

	if actionType != "" {
		resp.Diagnostics.Append(
			validatordiag.InvalidAttributeValueDiagnostic(
				req.Path,
				fmt.Sprintf("Please use the new \"%s\" instead of the generic \"action\".", actionType),
				req.ConfigValue.String()))
	}
}

func ActionTypeHasSpecificImplementation() validator.String {
	return actionTypeHasImplementationValidator{}
}
