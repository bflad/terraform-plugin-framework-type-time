package timetypes_test

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/bflad/terraform-plugin-framework-type-time/timetypes"
	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TestRFC3339Equal(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		other    attr.Value
		expected bool
	}{
		"nil": {
			value:    timetypes.RFC3339Null(),
			other:    nil,
			expected: false,
		},
		"not-timetypes.RFC3339": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    types.String{Value: "2006-01-02T15:04:05Z"},
			expected: false,
		},
		"null-null": {
			value:    timetypes.RFC3339Null(),
			other:    timetypes.RFC3339Null(),
			expected: true,
		},
		"null-unknown": {
			value:    timetypes.RFC3339Null(),
			other:    timetypes.RFC3339Unknown(),
			expected: false,
		},
		"null-value": {
			value:    timetypes.RFC3339Null(),
			other:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
		"unknown-null": {
			value:    timetypes.RFC3339Unknown(),
			other:    timetypes.RFC3339Null(),
			expected: false,
		},
		"unknown-unknown": {
			value:    timetypes.RFC3339Unknown(),
			other:    timetypes.RFC3339Unknown(),
			expected: true,
		},
		"unknown-value": {
			value:    timetypes.RFC3339Unknown(),
			other:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: false,
		},
		"value-null": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    timetypes.RFC3339Null(),
			expected: false,
		},
		"value-unknown": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    timetypes.RFC3339Unknown(),
			expected: false,
		},
		"value-value-different": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    timetypes.RFC3339Time(time.Date(2007, 2, 3, 16, 5, 6, 1, time.UTC)),
			expected: false,
		},
		"value-value-equal": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			other:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

func TestRFC3339IsNull(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		expected bool
	}{
		"null": {
			value:    timetypes.RFC3339Null(),
			expected: true,
		},
		"unknown": {
			value:    timetypes.RFC3339Unknown(),
			expected: false,
		},
		"value": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

func TestRFC3339IsUnknown(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		expected bool
	}{
		"null": {
			value:    timetypes.RFC3339Null(),
			expected: false,
		},
		"unknown": {
			value:    timetypes.RFC3339Unknown(),
			expected: true,
		},
		"value": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

func TestRFC3339String(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		expected string
	}{
		"null": {
			value:    timetypes.RFC3339Null(),
			expected: "<null>",
		},
		"unknown": {
			value:    timetypes.RFC3339Unknown(),
			expected: "<unknown>",
		},
		"value-offset-negative": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: "\"2006-01-02T15:04:05-07:00\"",
		},
		"value-offset-positive": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: "\"2006-01-02T15:04:05+07:00\"",
		},
		"value-z": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

func TestRFC3339Time(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		expected time.Time
	}{
		"null": {
			value:    timetypes.RFC3339Null(),
			expected: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"unknown": {
			value:    timetypes.RFC3339Unknown(),
			expected: time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC),
		},
		"value-offset-negative": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60)),
		},
		"value-offset-positive": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60)),
		},
		"value-z": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
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

func TestRFC3339ToTerraformValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value         timetypes.RFC3339
		expected      tftypes.Value
		expectedError error
	}{
		"null": {
			value:    timetypes.RFC3339Null(),
			expected: tftypes.NewValue(tftypes.String, nil),
		},
		"unknown": {
			value:    timetypes.RFC3339Unknown(),
			expected: tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		},
		"value-offset-negative": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", -7*60*60))),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05-07:00"),
		},
		"value-offset-positive": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.FixedZone("", 7*60*60))),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05+07:00"),
		},
		"value-z": {
			value:    timetypes.RFC3339Time(time.Date(2006, 1, 2, 15, 4, 5, 0, time.UTC)),
			expected: tftypes.NewValue(tftypes.String, "2006-01-02T15:04:05Z"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := testCase.value.ToTerraformValue(context.Background())

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

func TestRFC3339Type(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		value    timetypes.RFC3339
		expected attr.Type
	}{
		"any": {
			value:    timetypes.RFC3339Null(),
			expected: timetypes.RFC3339Type{},
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
