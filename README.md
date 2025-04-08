# Fillmore Labs zerolint

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/zerolint.svg)](https://pkg.go.dev/fillmore-labs.com/zerolint)
[![Test](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml)
[![Coverage](https://codecov.io/gh/fillmore-labs/zerolint/branch/main/graph/badge.svg?token=TUE1BL1QZV)](https://codecov.io/gh/fillmore-labs/zerolint)
[![Maintainability](https://api.codeclimate.com/v1/badges/baf50ad423cc30ff7790/maintainability)](https://codeclimate.com/github/fillmore-labs/zerolint/maintainability)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/zerolint)](https://goreportcard.com/report/fillmore-labs.com/zerolint)
[![License](https://img.shields.io/github/license/fillmore-labs/zerolint)](https://www.apache.org/licenses/LICENSE-2.0)

`zerolint` is a Go static analysis tool (linter) that detects potentially wrong or unnecessary usage of pointers to
zero-sized types.

## Motivation

Go's zero-sized types, like `struct{}` or `[0]byte`, occupy no memory. While useful in certain contexts (e.g.,
signaling on channels, map keys/values for set semantics), taking the address of a zero-sized variable (`&struct{}{}`)
or allocating them (`new(struct{})`) is often redundant.

All values of a zero-sized type are indistinguishable, so pointers to them generally don't convey unique state or
identity. Using pointers to zero-sized types can obscure intent, as readers might mistakenly assume the pointer implies
state or identity management. Also, since pointers are not zero-sized, it introduces minor performance overhead and
waste of memory.

`zerolint` helps identify these patterns, promoting cleaner and potentially more efficient code.

## Usage

```console
# Install the linter
go install fillmore-labs.com/zerolint@latest
```

Usage: `zerolint [-flag] [package]`

Flags:

- **-c** int display offending line with this many lines of context (default -1)
- **-level** Perform a more comprehensive analysis (0=default, 1=extended, 2=full)
- **-generated** Analyze files that contain code generation markers (e.g., `// Code generated ... DO NOT EDIT.`). By
default, these files are skipped.
- **-excluded** `<filename>` Read types to be excluded from analysis from the specified file. The file should contain
fully qualified type names, one per line. See Excluding Types.
- **-zerotrace** Enable verbose logging which types zerolint identifies as zero-sized. Useful for building a list of
excluded types.
- **-fix** Apply all suggested fixes automatically. Use with caution and always review the changes made by `-fix`.

## Examples

Consider the following Go code:

```go
// main.go
package main

type myError struct{}

func (*myError) Error() string {
	return "my error"
}

func processing() error {
	return &myError{}
}

func main() {
	if err := processing(); err != nil {
		panic(err)
	}
}
```

Running `zerolint % zerolint -level 2 ./...` would produce output similar to:

```text
/path/to/your/project/main.go:6:7: error interface implemented on pointer to zero-sized type "example.com/project.myError" (zl:err)
/path/to/your/project/main.go:11:9: address of zero-size variable of type "example.com/project.myError" (zl:add)
```

## Excluding Types

If you have specific zero-sized types where pointer usage is intentional or required (e.g., due to external library
constraints), you can exclude them using the `-excluded` flag with a file path.

Example `excludes.txt`:

```text
# zerolint excludes for my project
example.com/some/pkg.RequiredPointerType
vendor.org/library.MarkerInterfaceImplementation
```

Then run: `zerolint -excluded=excludes.txt ./...`

## Integration

See `zerolint-golangci-plugin`

## Known Bugs

We are aware of a number of minor bugs in the analyzer's fixes. For example, it may sometimes cause a type not
implement an interface or check a non-pointer type for nil. The known bugs are low-risk and easy to fix, as they result
in a broken build or are obvious during a code review; none cause latent behavior changes. Please report any additional
problems you encounter.
