package schemas

import (
	"fmt"

	"github.com/OctopusDeploy/go-octopusdeploy/v2/pkg/runbooks"
	"github.com/hashicorp/terraform-plugin-framework-validators/boolvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	resourceSchema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var runbookRetentionPeriodSchemeAttributeNames = struct {
	QuantityToKeep    string
	ShouldKeepForever string
}{
	QuantityToKeep:    "quantity_to_keep",
	ShouldKeepForever: "should_keep_forever",
}

func GetRunbookRetentionPeriodObjectType() map[string]attr.Type {
	return map[string]attr.Type{
		runbookRetentionPeriodSchemeAttributeNames.QuantityToKeep:    types.Int64Type,
		runbookRetentionPeriodSchemeAttributeNames.ShouldKeepForever: types.BoolType,
	}
}

func getRunbookRetentionPeriodSchema() map[string]resourceSchema.Attribute {
	return map[string]resourceSchema.Attribute{
		runbookRetentionPeriodSchemeAttributeNames.QuantityToKeep: resourceSchema.Int64Attribute{
			Description: "How many runs to keep per environment.",
			Computed:    true,
			Optional:    true,
			Validators: []validator.Int64{
				int64validator.AtLeast(0),
			},
		},
		runbookRetentionPeriodSchemeAttributeNames.ShouldKeepForever: resourceSchema.BoolAttribute{
			Description: "Indicates if items should never be deleted. The default value is `false`.",
			Computed:    true,
			Optional:    true,
			Default:     booldefault.StaticBool(false),
			Validators: []validator.Bool{
				boolvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName(runbookRetentionPeriodSchemeAttributeNames.QuantityToKeep)),
			},
		},
	}
}

func GetDefaultRunbookRetentionPeriod() *runbooks.RunbookRetentionPeriod {
	return &runbooks.RunbookRetentionPeriod{
		QuantityToKeep:    100,
		ShouldKeepForever: false,
	}
}

func MapFromRunbookRetentionPeriod(retentionPeriod *runbooks.RunbookRetentionPeriod) attr.Value {
	if retentionPeriod == nil {
		return MapFromRunbookRetentionPeriod(GetDefaultRunbookRetentionPeriod())
	}

	attrs := map[string]attr.Value{
		runbookRetentionPeriodSchemeAttributeNames.QuantityToKeep:    types.Int64Value(int64(retentionPeriod.QuantityToKeep)),
		runbookRetentionPeriodSchemeAttributeNames.ShouldKeepForever: types.BoolValue(retentionPeriod.ShouldKeepForever),
	}

	return types.ObjectValueMust(GetRunbookRetentionPeriodObjectType(), attrs)
}

func MapToRunbookRetentionPeriod(flattenedRunbookRetentionPeriod types.List) *runbooks.RunbookRetentionPeriod {
	if flattenedRunbookRetentionPeriod.IsNull() || len(flattenedRunbookRetentionPeriod.Elements()) == 0 {
		return GetDefaultRunbookRetentionPeriod()
	}
	obj := flattenedRunbookRetentionPeriod.Elements()[0].(types.Object)
	attrs := obj.Attributes()

	var runbookRetentionPeriod runbooks.RunbookRetentionPeriod
	if quantityToKeep, ok := attrs[runbookRetentionPeriodSchemeAttributeNames.QuantityToKeep].(types.Int64); ok && !quantityToKeep.IsNull() {
		runbookRetentionPeriod.QuantityToKeep = int32(quantityToKeep.ValueInt64())
	}
	if shouldKeepForever, ok := attrs[runbookRetentionPeriodSchemeAttributeNames.ShouldKeepForever].(types.Bool); ok && !shouldKeepForever.IsNull() {
		runbookRetentionPeriod.ShouldKeepForever = shouldKeepForever.ValueBool()
	}
	fmt.Printf("runbook retention period: %#v", runbookRetentionPeriod)
	return &runbookRetentionPeriod
}
