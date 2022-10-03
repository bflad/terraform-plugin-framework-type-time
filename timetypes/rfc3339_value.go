package timetypes

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// Ensure implementation satisfies expected interfaces.
var (
	_ attr.Value = RFC3339{}
)

// NullRFC3339 returns a null Value.
func NullRFC3339() RFC3339 {
	return RFC3339{
		null: true,
	}
}

// RFC3339String returns a known Value or any errors while attempting
// to parse the string as RFC 3339 format.
func RFC3339String(s string, schemaPath path.Path) (RFC3339, diag.Diagnostics) {
	t, err := time.Parse(time.RFC3339, s)

	if err != nil {
		return RFC3339{
				unknown: true,
			}, diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					schemaPath,
					"Invalid RFC 3339 String Value",
					"An unexpected error occurred while converting a string value that was expected to be RFC 3339 format. "+
						"The RFC 3339 string format is YYYY-MM-DDTHH:MM:SSZ, such as 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00.\n\n"+
						"Error: "+err.Error(),
				),
			}
	}

	return RFC3339{
		value: t,
	}, nil
}

// RFC3339Time returns a known Value with the given time.
func RFC3339Time(t time.Time) RFC3339 {
	return RFC3339{
		value: t,
	}
}

// UnknownRFC3339 returns an unknown Value.
func UnknownRFC3339() RFC3339 {
	return RFC3339{
		unknown: true,
	}
}

// RFC3339 implements the attr.Value interface for usage in logic.
type RFC3339 struct {
	null    bool
	unknown bool
	value   time.Time
}

// Equal returns true if the given attr.Value matches the following:
//   - Is a Value type
//   - Has the same null, unknown, and value data
func (v RFC3339) Equal(o attr.Value) bool {
	otherValue, ok := o.(RFC3339)

	if !ok {
		return false
	}

	if otherValue.null != v.null {
		return false
	}

	if otherValue.unknown != v.unknown {
		return false
	}

	return otherValue.value.Equal(v.value)
}

// IsNull returns true if the Value represents a null Value.
func (v RFC3339) IsNull() bool {
	return v.null
}

// IsUnknown returns true if the Value represents an unknown Value.
func (v RFC3339) IsUnknown() bool {
	return v.unknown
}

// String returns a human readable string of the Value.
func (v RFC3339) String() string {
	if v.null {
		return attr.NullValueString
	}

	if v.unknown {
		return attr.UnknownValueString
	}

	return `"` + v.value.Format(time.RFC3339) + `"`
}

// Time returns the time.Time of a Value.
func (v RFC3339) Time() time.Time {
	return v.value
}

// ToTerraformValue converts the Value to a tftypes.String.
func (v RFC3339) ToTerraformValue(_ context.Context) (tftypes.Value, error) {
	if v.null {
		return tftypes.NewValue(tftypes.String, nil), nil
	}

	if v.unknown {
		return tftypes.NewValue(tftypes.String, tftypes.UnknownValue), nil
	}

	return tftypes.NewValue(tftypes.String, v.value.Format(time.RFC3339)), nil
}

// Type returns Type.
func (v RFC3339) Type(_ context.Context) attr.Type {
	return RFC3339Type{}
}
