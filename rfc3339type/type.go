package rfc3339type

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
	_ tftypes.AttributePathStepper = Type{}
	_ attr.Type                    = Type{}
	_ xattr.TypeWithValidate       = Type{}
)

// Type implements the attr.Type interface for usage in schema definitions
// and data models.
type Type struct{}

// ApplyTerraform5AttributePathStep always returns an error as this type
// cannot be walked any further.
func (t Type) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (any, error) {
	return nil, fmt.Errorf("cannot apply AttributePathStep %T to %s", step, t.String())
}

// Equal returns true if the given type is Type.
func (t Type) Equal(o attr.Type) bool {
	_, ok := o.(Type)

	return ok
}

// String returns a human readable string of the type.
func (t Type) String() string {
	return "rfc3339type.Type"
}

// TerraformType always returns tftypes.String.
func (t Type) TerraformType(_ context.Context) tftypes.Type {
	return tftypes.String
}

// Validate ensures the value is always RFC 3339 conformant.
func (t Type) Validate(_ context.Context, terraformValue tftypes.Value, schemaPath path.Path) diag.Diagnostics {
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

	_, diags := StringValue(str, schemaPath)

	return diags
}

// ValueFromTerraform converts the tftypes.Value into a value.
func (t Type) ValueFromTerraform(_ context.Context, terraformValue tftypes.Value) (attr.Value, error) {
	if terraformValue.IsNull() {
		return NullValue(), nil
	}

	if !terraformValue.IsKnown() {
		return UnknownValue(), nil
	}

	var str string

	err := terraformValue.As(&str)

	if err != nil {
		return UnknownValue(), err
	}

	strTime, err := time.Parse(time.RFC3339, str)

	if err != nil {
		return UnknownValue(), err
	}

	return TimeValue(strTime), nil
}
