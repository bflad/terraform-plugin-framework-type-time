package rfc3339type_test

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/bflad/terraform-plugin-framework-type-rfc3339/rfc3339type"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestTypeApplyTerraform5AttributePathStep(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ           rfc3339type.Type
		step          tftypes.AttributePathStep
		expected      any
		expectedError error
	}{
		"AttributeName": {
			typ:           rfc3339type.Type{},
			step:          tftypes.AttributeName("test"),
			expectedError: fmt.Errorf("cannot apply AttributePathStep tftypes.AttributeName to rfc3339type.Type"),
		},
		"ElementKeyInt": {
			typ:           rfc3339type.Type{},
			step:          tftypes.ElementKeyInt(1),
			expectedError: fmt.Errorf("cannot apply AttributePathStep tftypes.ElementKeyInt to rfc3339type.Type"),
		},
		"ElementKeyString": {
			typ:           rfc3339type.Type{},
			step:          tftypes.ElementKeyString("test"),
			expectedError: fmt.Errorf("cannot apply AttributePathStep tftypes.ElementKeyString to rfc3339type.Type"),
		},
		"ElementKeyValue": {
			typ:           rfc3339type.Type{},
			step:          tftypes.ElementKeyValue{},
			expectedError: fmt.Errorf("cannot apply AttributePathStep tftypes.ElementKeyValue to rfc3339type.Type"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.typ.ApplyTerraform5AttributePathStep(testCase.step)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("expected no error, got: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, tfType: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTypeEqual(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ      rfc3339type.Type
		other    attr.Type
		expected bool
	}{
		"nil": {
			typ:      rfc3339type.Type{},
			other:    nil,
			expected: false,
		},
		"rfc3339type.Type": {
			typ:      rfc3339type.Type{},
			other:    rfc3339type.Type{},
			expected: true,
		},
		"types.StringType": {
			typ:      rfc3339type.Type{},
			other:    types.StringType,
			expected: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.typ.Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTypeString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ      rfc3339type.Type
		expected string
	}{
		"any": {
			typ:      rfc3339type.Type{},
			expected: "rfc3339type.Type",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.typ.String()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTypeTerraformType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ      rfc3339type.Type
		expected tftypes.Type
	}{
		"any": {
			typ:      rfc3339type.Type{},
			expected: tftypes.String,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.typ.TerraformType(context.Background())

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestTypeValidate(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ            rfc3339type.Type
		terraformValue tftypes.Value
		schemaPath     path.Path
		expectedDiags  diag.Diagnostics
	}{
		"not-string": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.Bool, true),
			schemaPath:     path.Root("test"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid RFC 3339 Terraform Value",
					"An unexpected error occurred while attempting to read a RFC 3339 string from the Terraform value. "+
						"Please contact the provider developers with the following:\n\n"+
						"Error: can't unmarshal tftypes.Bool into *string, expected string",
				),
			},
		},
		"string-null": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, nil),
			schemaPath:     path.Root("test"),
		},
		"string-unknown": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			schemaPath:     path.Root("test"),
		},
		"string-value-invalid": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "not-rfc3339-format"),
			schemaPath:     path.Root("test"),
			expectedDiags: diag.Diagnostics{
				diag.NewAttributeErrorDiagnostic(
					path.Root("test"),
					"Invalid RFC 3339 String Value",
					"An unexpected error occurred while converting a string value that was expected to be RFC 3339 format. "+
						"The RFC 3339 string format is YYYY-MM-DDTHH:MM:SSZ, such as 2006-01-02T15:04:05Z or 2006-01-02T15:04:05+07:00.\n\n"+
						"Error: parsing time \"not-rfc3339-format\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"not-rfc3339-format\" as \"2006\"",
				),
			},
		},
		"string-value-valid-offset-negative": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05-07:00"),
			schemaPath:     path.Root("test"),
		},
		"string-value-valid-offset-positive": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05+07:00"),
			schemaPath:     path.Root("test"),
		},
		"string-value-valid-z": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05Z"),
			schemaPath:     path.Root("test"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			diags := testCase.typ.Validate(context.Background(), testCase.terraformValue, testCase.schemaPath)

			if diff := cmp.Diff(diags, testCase.expectedDiags); diff != "" {
				t.Errorf("unexpected diagnostics difference: %s", diff)
			}
		})
	}
}

func TestTypeValueFromTerraform(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		typ            rfc3339type.Type
		terraformValue tftypes.Value
		expected       attr.Value
		expectedError  error
	}{
		"not-string": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.Bool, true),
			expected:       rfc3339type.UnknownValue(),
			expectedError:  fmt.Errorf("can't unmarshal tftypes.Bool into *string, expected string"),
		},
		"string-null": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, nil),
			expected:       rfc3339type.NullValue(),
		},
		"string-unknown": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
			expected:       rfc3339type.UnknownValue(),
		},
		"string-value-invalid": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "not-rfc3339-format"),
			expected:       rfc3339type.UnknownValue(),
			expectedError:  fmt.Errorf("parsing time \"not-rfc3339-format\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"not-rfc3339-format\" as \"2006\""),
		},
		"string-value-valid-offset-negative": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05-07:00"),
			expected:       rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
		},
		"string-value-valid-offset-positive": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05+07:00"),
			expected:       rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
		},
		"string-value-valid-z": {
			typ:            rfc3339type.Type{},
			terraformValue: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05Z"),
			expected:       rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.typ.ValueFromTerraform(context.Background(), testCase.terraformValue)

			if err != nil {
				if testCase.expectedError == nil {
					t.Fatalf("expected no error, got: %s", err)
				}

				if !strings.Contains(err.Error(), testCase.expectedError.Error()) {
					t.Fatalf("expected error %q, got: %s", testCase.expectedError, err)
				}
			}

			if err == nil && testCase.expectedError != nil {
				t.Fatalf("got no error, tfType: %s", testCase.expectedError)
			}

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
