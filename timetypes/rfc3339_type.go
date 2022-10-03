package timetypes

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/attr/xattr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure implementation satisfies expected interfaces.
var (
	_ tftypes.AttributePathStepper = RFC3339Type{}
	_ attr.Type                    = RFC3339Type{}
	_ xattr.TypeWithValidate       = RFC3339Type{}
)

// RFC3339Type implements the attr.Type interface for usage in schema definitions
// and data models.
type RFC3339Type struct{}

// ApplyTerraform5AttributePathStep always returns an error as this type
// cannot be walked any further.
func (t RFC3339Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is Type.
func (t RFC3339Type) Equal(o attr.Type) bool {
	_, ok := o.(RFC3339Type)

	return ok
}

// String returns a human readable string of the type.
func (t RFC3339Type) String() string {
	return "timetypes.RFC3339Type"
}

// TerraformType always returns tftypes.String.
func (t RFC3339Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// Validate ensures the value is always RFC 3339 conformant.
func (t RFC3339Type) Validate(_ context.Context, terraformValue tftypes.Value, schemaPath path.Path) diag.Diagnostics {
	if terraformValue.IsNull() || !terraformValue.IsKnown() {
		return nil
	}

	var str string

	err := terraformValue.As(&str)

	if err != nil {
		return diag.Diagnostics{
			diag.NewAttributeErrorDiagnostic(
				schemaPath,
				"Invalid RFC 3339 Terraform Value",
				"An unexpected error occurred while attempting to read a RFC 3339 string from the Terraform value. "+
					"Please contact the provider developers with the following:\n\n"+
					"Error: "+err.Error(),
			),
		}
	}

	_, diags := RFC3339String(str, schemaPath)

	return diags
}

// ValueFromTerraform converts the tftypes.Value into a value.
func (t RFC3339Type) ValueFromTerraform(_ context.Context, terraformValue tftypes.Value) (attr.Value, error) {
	if terraformValue.IsNull() {
		return NullRFC3339(), nil
	}

	if !terraformValue.IsKnown() {
		return UnknownRFC3339(), nil
	}

	var str string

	err := terraformValue.As(&str)

	if err != nil {
		return UnknownRFC3339(), err
	}

	strTime, err := time.Parse(time.RFC3339, str)

	if err != nil {
		return UnknownRFC3339(), err
	}

	return RFC3339Time(strTime), nil
}

// ValueType returns the associated attr.Value.
func (t RFC3339Type) ValueType(_ context.Context) attr.Value {
	return RFC3339{}
}
