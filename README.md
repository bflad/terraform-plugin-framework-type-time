# terraform-plugin-framework-type-time

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bflad/terraform-plugin-framework-type-time)](https://pkg.go.dev/github.com/bflad/terraform-plugin-framework-type-time)

[RFC 3339](https://tools.ietf.org/html/rfc3339) timestamp string type and value implementations for [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework). These types automatically include syntax validation and support Go standard library [`time.Time`](https://pkg.go.dev/time#Time) conversion.

## Status

_Experimental_: This project is not officially maintained by HashiCorp.

## Prerequisites

This Go module tracks recent versions of `terraform-plugin-framework` for Go version and interface compatibility.

## Getting Started

### Schema

Replace usage of `types.StringType` in schema definitions with `timetypes.RFC3339Type{}`.

Given the previous schema attribute:

```go
tfsdk.Attribute{
    Required: true,
    Type:     types.StringType,
    // Potentially previous Validators
}
```

The updated schema attribute:

```go
tfsdk.Attribute{
    Required: true,
    Type:     timetypes.RFC3339Type{},
}
```

### Schema Data Model

Replace usage of `string`, `*string`, or `types.String` in schema data models with `timetypes.RFC3339`.

Given the previous schema data model:

```go
type ThingResourceModel struct {
    // ...
    Example types.String `tfsdk:"example"`
}
```

The updated schema data model:

```go
type ThingResourceModel struct {
    // ...
    Example timetypes.RFC3339 `tfsdk:"example"`
}
```

### Accessing Values

Similar to other value types, use the `IsNull()` and `IsUnknown()` methods to check whether the value is null or unknown. Use the `Time()` method to extract a known `time.Time` value.

### Writing Values

Create a `timetypes.RFC3339` by calling one of these functions:

- `RFC3339Null() RFC3339`: creates a `null` value.
- `RFC3339String(string, path.Path) (Value, diag.Diagnostics)`: creates a known value using the given `string` or returns validation errors if `string` is not in the expected RFC 3339 format.
- `RFC3339Time(time.Time) Value` creates a known value using the given `time.Time`.
- `RFC3339Unknown() Value`: creates an unknown value.

### Adding the Dependency

All functionality is located in the `github.com/bflad/terraform-plugin-framework-type-time/timetypes` package. Add this to relevant Go file `import` statements.

Run these Go module commands to fetch the latest version and ensure all module files are up to date.

```shell
go get github.com/bflad/terraform-plugin-framework-type-time@latest
go mod tidy
```
