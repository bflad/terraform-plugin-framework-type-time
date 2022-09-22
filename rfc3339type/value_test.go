package rfc3339type_test

import (
	"context"
	"testing"
	"time"

	"github.com/bflad/terraform-plugin-framework-type-rfc3339/rfc3339type"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestValueEqual(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		other    attr.Value
		expected bool
	}{
		"nil": {
			value:    rfc3339type.NullValue(),
			other:    nil,
			expected: false,
		},
		"not-rfc3339type.Value": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    types.String{Value: "2006-01-02T15:04:05Z"},
			expected: false,
		},
		"null-null": {
			value:    rfc3339type.NullValue(),
			other:    rfc3339type.NullValue(),
			expected: true,
		},
		"null-unknown": {
			value:    rfc3339type.NullValue(),
			other:    rfc3339type.UnknownValue(),
			expected: false,
		},
		"null-value": {
			value:    rfc3339type.NullValue(),
			other:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
		"unknown-null": {
			value:    rfc3339type.UnknownValue(),
			other:    rfc3339type.NullValue(),
			expected: false,
		},
		"unknown-unknown": {
			value:    rfc3339type.UnknownValue(),
			other:    rfc3339type.UnknownValue(),
			expected: true,
		},
		"unknown-value": {
			value:    rfc3339type.UnknownValue(),
			other:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
		"value-null": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    rfc3339type.NullValue(),
			expected: false,
		},
		"value-unknown": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    rfc3339type.UnknownValue(),
			expected: false,
		},
		"value-value-different": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    rfc3339type.TimeValue(time.Date(2007, 2, 3, 16, 5, 6, 1, time.UTC)),
			expected: false,
		},
		"value-value-equal": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: true,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.Equal(testCase.other)

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueIsNull(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected bool
	}{
		"null": {
			value:    rfc3339type.NullValue(),
			expected: true,
		},
		"unknown": {
			value:    rfc3339type.UnknownValue(),
			expected: false,
		},
		"value": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.IsNull()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueIsUnknown(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected bool
	}{
		"null": {
			value:    rfc3339type.NullValue(),
			expected: false,
		},
		"unknown": {
			value:    rfc3339type.UnknownValue(),
			expected: true,
		},
		"value": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.IsUnknown()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueString(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected string
	}{
		"null": {
			value:    rfc3339type.NullValue(),
			expected: "<null>",
		},
		"unknown": {
			value:    rfc3339type.UnknownValue(),
			expected: "<unknown>",
		},
		"value-offset-negative": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: "\"2006-01-02T15:04:05-07:00\"",
		},
		"value-offset-positive": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: "\"2006-01-02T15:04:05+07:00\"",
		},
		"value-z": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: "\"2006-01-02T15:04:05Z\"",
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.String()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueTime(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected time.Time
	}{
		"null": {
			value:    rfc3339type.NullValue(),
			expected: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"unknown": {
			value:    rfc3339type.UnknownValue(),
			expected: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"value-offset-negative": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
		},
		"value-offset-positive": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60)),
		},
		"value-z": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.Time()

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueToTerraformValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected tftypes.Value
	}{
		"null": {
			value:    rfc3339type.NullValue(),
			expected: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			value:    rfc3339type.UnknownValue(),
			expected: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"value-offset-negative": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05-07:00"),
		},
		"value-offset-positive": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05+07:00"),
		},
		"value-z": {
			value:    rfc3339type.TimeValue(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05Z"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.Type(context.Background())

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestValueType(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    rfc3339type.Value
		expected attr.Type
	}{
		"any": {
			value:    rfc3339type.NullValue(),
			expected: rfc3339type.Type{},
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.value.Type(context.Background())

			if diff := cmp.Diff(got, testCase.expected); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}
