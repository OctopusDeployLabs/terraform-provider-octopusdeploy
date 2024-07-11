package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func Create(ctx context.Context, resource string, v ...any) {
	tflog.Info(ctx, fmt.Sprintf("creating %s: %#v", resource, v))
	defer tflog.Info(ctx, fmt.Sprintf("%s created"))
}
