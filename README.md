# terraform-plugin-framework-type-rfc3339

[![PkgGoDev](https://pkg.go.dev/badge/github.com/bflad/terraform-plugin-framework-type-rfc3339)](https://pkg.go.dev/github.com/bflad/terraform-plugin-framework-type-rfc3339)

[RFC 3339](https://tools.ietf.org/html/rfc3339) timestamp string type and value implementations for [terraform-plugin-framework](https://github.com/hashicorp/terraform-plugin-framework). These types automatically include syntax validation and support Go standard library [`time.Time`](https://pkg.go.dev/time#Time) conversion.

## Status

_Experimental_: This project is not officially maintained by HashiCorp.

## Prerequisites

This Go module tracks recent versions of `terraform-plugin-framework` for Go version and interface compatibility.

## Getting Started

### Schema

Replace usage of `types.StringType` in schema definitions with `rfc3339type.Type{}`.

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
    Type:     rfc3339type.Type{},
}
```

### Schema Data Model

Replace usage of `string`, `*string`, or `types.String` in schema data models with `rfc3339type.Value`.

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
    Example rfc3339type.Value `tfsdk:"example"`
}
```

### Accessing Values

Similar to other value types, use the `IsNull()` and `IsUnknown()` methods to check whether the value is null or unknown. Use the `Time()` method to extract a known `time.Time` value.

### Writing Values

Create a `rfc3339type.Value` by calling one of these functions:

- `NullValue() Value`: creates a `null` value.
- `UnknownValue() Value`: creates an unknown value.
- `StringValue(string, path.Path) (Value, diag.Diagnostics)`: creates a known value using the given `string` or returns validation errors if `string` is not in the expected RFC 3339 format.
- `TimeValue(time.Time) Value` creates a known value using the given `time.Time`.

### Adding the Dependency

All functionality is located in the `github.com/bflad/terraform-plugin-framework-type-rfc3339/rfc3339type` package. Add this to relevant Go file `import` statements.

Run these Go module commands to fetch the latest version and ensure all module files are up to date.

```shell
go get github.com/bflad/terraform-plugin-framework-type-rfc3339@latest
go mod tidy
```
