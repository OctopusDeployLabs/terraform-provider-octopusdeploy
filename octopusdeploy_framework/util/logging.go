package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func Create(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("creating %s: %#v", resource, v))
}

func Created(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("created %s: %#v", resource, v))
}

func Delete(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("deleting %s: %#v", resource, v))
}

func Deleted(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("deleted %s: %#v", resource, v))
}

func Reading(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("reading %s: %#v", resource, v))
}

func DatasourceReading(ctx context.Context, resource string, v ...any) {
	tflog.Debug(ctx, fmt.Sprintf("reading %s data source with query: %+v", resource, v))
}

func DatasourceResultCount(ctx context.Context, resource string, count int) {
	tflog.Debug(ctx, fmt.Sprintf("reading %s returned %d items", resource, count))
}

func Read(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("read %s: %#v", resource, v))
}

func Update(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("updating %s: %#v", resource, v))
}

func Updated(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("updated %s: %#v", resource, v))
}

// AddDiagnosticError is used to wrap calls to Diagnostics.AddError with additional information about the executable and versions
func AddDiagnosticError(diagnostics *diag.Diagnostics, systemInfo *SystemInfo, message string, err string) {
	suffix := "\n\nPlease ensure these details are included in any error report you raise.\n" +
		"Executable: " + systemInfo.GetExecutableName() + "\n" +
		"Terraform Version: " + systemInfo.TerraformVersion + "\n" +
		"Octopus Version: " + systemInfo.OctopusVersion

	diagnostics.AddError(message, err+suffix)
}
