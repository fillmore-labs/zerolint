# Zerolint

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/zerolint.svg)](https://pkg.go.dev/fillmore-labs.com/zerolint)
[![Test](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml/badge.svg?branch=dev)](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml?query=branch%3Adev)
[![CodeQL](https://github.com/fillmore-labs/zerolint/actions/workflows/github-code-scanning/codeql/badge.svg?branch=dev)](https://github.com/fillmore-labs/zerolint/actions/workflows/github-code-scanning/codeql?query=branch%3Adev)
[![Coverage](https://codecov.io/gh/fillmore-labs/zerolint/branch/dev/graph/badge.svg?token=TUE1BL1QZV)](https://codecov.io/gh/fillmore-labs/zerolint)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/zerolint)](https://goreportcard.com/report/fillmore-labs.com/zerolint)
[![License](https://img.shields.io/github/license/fillmore-labs/zerolint)](https://www.apache.org/licenses/LICENSE-2.0)

`zerolint` is a Go static analysis tool (linter) that detects unnecessary or potentially incorrect usage of pointers to
zero-sized types.

## Motivation

Go's zero-sized types (such as `struct{}` or `[0]byte`) occupy no memory and are useful in scenarios like channel
signaling or as map keys. However, creating pointers to these types (e.g., `&struct{}` or `new(struct{})`) is almost
always unnecessary and can introduce subtle bugs and overhead.

Since all values of a zero-sized type are identical, pointers to them rarely convey unique state or identity. This can
make code less clear, as readers might incorrectly assume state or identity is being managed. Furthermore, pointers
themselves are not zero-sized, leading to minor memory and performance overhead.

`zerolint` helps identify these patterns, promoting cleaner and potentially more efficient code.

## Quickstart

### Installation

Install the linter:

### Homebrew

```console
brew install fillmore-labs/tap/zerolint
```

### Go

```console
go install fillmore-labs.com/zerolint@latest
```

### Eget

[Install `eget`](https://github.com/zyedidia/eget?tab=readme-ov-file#how-to-get-eget), then

```console
eget fillmore-labs/zerolint
```

## Usage

Run the linter on your project:

```console
zerolint ./...
```

See below for descriptions of available command-line flags.

## Optional Flags

Usage: `zerolint [-flag] [package]`

Flags:

- **-level** `<level>`: Set analysis depth:
  - **Basic**: Basic detection of pointer issues (Default)
  - **Extended**: Additional checks for more complex patterns
  - **Full**: Most comprehensive analysis, recommended with `-fix`
- **-match** `<regex>`: Limit zero-sized type detection to types matching the regex. Useful with `-fix`.
- **-excluded** `<filename>`: Read types to be excluded from analysis from the specified file. The file should contain
  fully qualified type names, one per line. See the [“Exclusion File”](#exclusion-file) section for more details.
- **-generated**: Analyze files that contain code generation markers (e.g., `// Code generated ... DO NOT EDIT.`). By
  default, these files are skipped.
- **-zerotrace**: Enable verbose logging of which types `zerolint` identifies as zero-sized. Useful for building a list
  of excluded types.
- **-c** `<N>`: Display N lines of context around the offending line (default: -1 for no context, 0 for only the
  offending line).
- **-test**: Indicates whether test files should be analyzed, too. (default: true).
- **-fix**: Apply all suggested fixes automatically. Use with caution and always review the changes made by `-fix`.
- **-diff**: With `-fix`, don't update the files, but print a unified diff.

## Example

Consider the following Go code:

```go
package main

import (
	"errors"
	"testing"
)

type DivisionByZeroError struct{}

func (*DivisionByZeroError) Error() string {
	return "division by zero"
}

func Reciprocal(x float64) (float64, error) {
	if x == 0 {
		return 0, &DivisionByZeroError{}
	}

	return 1 / x, nil
}

func TestDivisionByZero(t *testing.T) {
	_, err := Reciprocal(0)

	if !errors.Is(err, &DivisionByZeroError{}) {
		t.Errorf("Expected division by zero error, but got: %v", err)
	}
}
```

The test passes ([Go Playground](https://go.dev/play/p/7Zyi1SsiSqI)):

```console
=== RUN   TestDivisionByZero
--- PASS: TestDivisionByZero (0.00s)
PASS
```

Running `zerolint .` would produce output similar to:

```text
/path/to/your/project/main_test.go:10:7: error interface implemented on pointer to zero-sized type "example.com/project.DivisionByZeroError" (zl:err)
/path/to/your/project/main_test.go:25:6: comparison of pointer to zero-size type "example.com/project.DivisionByZeroError" with error interface (zl:cme)
```

### Understanding the Output and Zero-Sized Diagnostics

The `main_test.go` example and its `zerolint` output highlight a common pitfall with zero-sized types in Go. In the
`TestDivisionByZero` function:

```go
	if !errors.Is(err, &DivisionByZeroError{}) {
```

the expression `&DivisionByZeroError{}` creates a new pointer to a zero-sized struct. Similarly, the `Reciprocal`
function, when `x` is 0, returns another `&DivisionByZeroError{}`. The critical point is how these pointers are
compared.

The check `errors.Is(err, &DivisionByZeroError{})` might not behave as intuitively expected. When the target error
passed to `errors.Is` is a pointer type (`*DivisionByZeroError` in this case), `errors.Is` first performs a direct
pointer address comparison, before even checking whether the error implements an `Is(error) bool` method.

To illustrate that these are treated as distinct pointers for comparison purposes, we can modify `DivisionByZeroError`
to be non-zero-sized:

```go
type DivisionByZeroError struct{ _ int } // Make it non-zero-sized
```

With this change, the test `TestDivisionByZero` fails, confirming that `errors.Is` was indeed comparing distinct
instances based on their pointer values.

#### Pitfalls of Zero-Sized Pointer Comparisons

Internally, Go's runtime optimizes allocations of zero-sized types. It achieves this by
[returning a pointer to a common static variable](https://cs.opensource.google/go/go/+/refs/tags/go1.24.6:src/runtime/malloc.go;l=1017-1020)
(known as `runtime.zerobase`) rather than allocating new memory on the heap for each instance. A consequence of this
optimization is that different pointers to zero-sized types (e.g., multiple uses of `&DivisionByZeroError{}` or
`new(DivisionByZeroError)`) end up pointing to the same memory address. This can create the illusion that such pointers
will always compare as equal.

Despite this common runtime behavior, the Go language specification
[_explicitly_ states](https://go.dev/ref/spec#Comparison_operators) that the equality of pointers to distinct zero-size
variables is unspecified:

> _“pointers to distinct zero-size variables may or may not be equal.”_

This means the observed consistency in pointer comparisons is an internal implementation detail of the Go runtime, not a
guaranteed language feature. Relying on this behavior is a classic example of [Hyrum's Law](https://www.hyrumslaw.com)
in action:

> _“With a sufficient number of users of an API, it does not matter what you promise in the contract: all observable
> behaviors of your system will be depended on by somebody.”_

Consequently, code that tests the equality (or inequality) of pointers to zero-sized types might compile and appear to
function correctly under current Go versions. However, its logic is fundamentally unsound because it depends on an
implementation detail not guaranteed by the language specification. Such code is at risk of breaking unexpectedly with
future Go updates or in different compilation environments. `zerolint` identifies and flags these problematic usage
patterns, helping developers write more robust code that avoids this undefined behavior.

### Potential Fixes

When `zerolint` flags an issue, consider these approaches:

#### Use a sentinel error variable (most idiomatic for errors)

This is often the clearest and most common Go practice for handling specific error conditions.

```go
// ErrDivisionByZero is returned when attempting to divide by zero.
var ErrDivisionByZero = errors.New("division by zero")
```

This approach is preferred because comparisons like `errors.Is(err, ErrDivisionByZero)` work reliably with sentinel
error values, avoiding the pitfalls of comparing pointers to zero-sized types.

#### Applying Fixes with `zerolint` (automatic refactoring)

For many common issues identified by `zerolint`, you can attempt an automatic fix:

```console
zerolint -level=full -fix ./...
```

The `-fix` flag will try to apply corrections, such as changing pointer receivers to value receivers where appropriate
or modifying how zero-sized types are instantiated or compared. For the most comprehensive automatic fixing, using
`-level=full` with `-fix` is recommended. This combination helps ensure that `zerolint` addresses all detected issues
related to a specific zero-sized type, promoting consistency across its usages once the fixes are applied.

> **Caution:** Always review changes made by `-fix` carefully before committing them, as automatic refactoring can
> sometimes have unintended consequences, especially in complex codebases.

For instance, running `zerolint -level=full -fix .` on the example above transforms the code as follows:

```go
package main

import (
	"errors"
	"testing"
)

type DivisionByZeroError struct{}

func (DivisionByZeroError) Error() string {
	return "division by zero"
}

func Reciprocal(x float64) (float64, error) {
	if x == 0 {
		return 0, DivisionByZeroError{}
	}

	return 1 / x, nil
}

func TestDivisionByZero(t *testing.T) {
	_, err := Reciprocal(0)

	if !errors.Is(err, DivisionByZeroError{}) {
		t.Errorf("Expected division by zero error, but got: %v", err)
	}
}
```

This program is correct, since the errors are compared by value, and two zero-sized variables of the same type always
compare equal.

#### Make the Type Non-Zero-Sized

If you need to maintain the custom error type structure for specific reasons (e.g., backward compatibility), or if it's
not an error type but another zero-sized struct, you can make it non-zero-sized. For errors, you can optionally provide
an `Is` method to restore the previous behavior of `errors.Is` when comparing against this error type:

```go
type DivisionByZeroError struct{ _ int } // Add a non-zero field

func (*DivisionByZeroError) Is(err error) bool { // Optional for error types
	_, ok := err.(*DivisionByZeroError)

	return ok
}
```

While this approach is more verbose than using `errors.New` (for errors) or the original pointer-based zero-sized error
implementation, it ensures correct, defined behavior for comparisons, making it valid Go code. This might be considered
if backward compatibility with an existing pointer-based error API is a concern, though migrating away from
pointer-based zero-sized errors is generally preferable.

#### Exclude the Type from Analysis

If pointer usage for a specific zero-sized type is intentional, unavoidable (e.g., due to external library constraints),
or you've assessed the risk and accept it, you can exclude the type from `zerolint`'s analysis. See the next section
[“Excluding Types”](#excluding-types) for details.

## Excluding Types

You can instruct `zerolint` to ignore specific zero-sized types in several ways:

### Exclusion File

If you have specific zero-sized types where pointer usage is intentional or required (e.g., due to external library
constraints), you can exclude them using the `-excluded` flag with a file path. The file should contain fully qualified
type names, one per line.

Example `excludes.txt`:

```text
# zerolint excludes for my project
company.example/service/client.RequestOptions
example.com/project.DivisionByZeroError
```

Then run: `zerolint -excluded=excludes.txt ./...`

This is especially useful when running with the `-fix` flag and dealing with types from external libraries you don't
control.

### Source Code Comment

If you control the source code where the zero-sized type is defined, you can add a special comment directly above the
type definition:

```go
//zerolint:exclude
type MyIntentionalZeroSizedType struct{}
```

This comment will tell `zerolint` to ignore any issues for `MyIntentionalZeroSizedType`.

To exclude a type defined in an external package, you can declare the exclusion in your own package using a `var`
declaration with the blank identifier (`_`):

```go
//zerolint:exclude
var _ external.ZeroSizedType
```

Using these exclusion methods allows you to tailor `zerolint`'s behavior to your project's specific needs.

## Linter Scope and External Types

By default, `zerolint` analyzes all types encountered, not just those declared within your current package or module.
This includes types imported from external packages (dependencies).

While `zerolint` (especially at its default analysis level) aims to flag only genuinely problematic patterns, there
might be situations where a zero-sized type from an external package is used in a way that, while flagged, is legitimate
or required by that external package's API. For example, an external library might require you to pass a pointer to a
zero-sized option structure or for interface satisfaction in a way that cannot be altered.

`zerolint` itself cannot automatically determine if such a flagged usage of an external type is intentional or
unavoidable within the constraints of that external library. It reports based on the general principle of avoiding
unnecessary pointers to zero-sized types.

If you encounter such a scenario with an external type you cannot modify with a `//zerolint:exclude` comment, the
recommended approach to manage these legitimate cases is:

1. Run `zerolint` with the `-zerotrace` flag. This will provide a detailed log of all types that `zerolint` identifies
   as zero-sized during its analysis.
2. Inspect this log to find the fully qualified names of the specific external types that are being flagged but whose
   usage you've determined is valid.
3. Manually add these fully qualified type names to an exclusion file, as described in the
   [“Excluding Types”](#excluding-types) section. This will instruct `zerolint` to ignore these specific types in future
   analyses.

This approach allows you to maintain the benefits of `zerolint` for your own codebase and other dependencies while
selectively bypassing checks for specific external types where pointer usage is justified.

## Diagnostic Codes

`zerolint` output includes diagnostic codes to help categorize the type of issue found. In the examples for each
diagnostic code, `zst` is used as a placeholder for a zero-sized type definition (e.g., `type zst struct{}`), and `zsv`
represents a variable of that zero-sized type (e.g., `var zsv zst`).

### Basic Level

- **zl:cme**: Comparison of pointer to zero-size type with an error interface (`errors.Is(err, &zsv)`)
- **zl:cmp**: Comparison of pointers to zero-size type (`&zsv == &zsv`)
- **zl:cmi**: Comparison of pointer to zero-size type with interface (`&zsv == any(&zst{})`)
- **zl:err**: Error interface implemented on pointer to zero-sized type (`func (*zst) Error() string`)
- **zl:emb**: Embedded pointer to zero-sized type (`struct{ *zst }`)
- **zl:der**: Dereferencing pointer to zero-size variable (`zsp := &zsv; _ = *zsp`)
- **zl:dcl**: Type declaration to pointer to zero-sized type (`type zstPtr *zst`)

### Extended Level

- **zl:new**: `new` called on zero-sized type (`new(zst)`)
- **zl:nil**: Cast of nil to pointer to zero-size type (`(*zst)(nil)`)
- **zl:ret**: Explicitly returning nil as pointer to zero-sized type (`func f() *zst { return nil }`)
- **zl:cst**: Cast to pointer to zero-size type (`(*zst)(&struct{}{})`)
- **zl:var**: Variable is pointer to zero-sized type (`var _ *zst`)
- **zl:fld**: Field points to zero-sized type (`struct{ f *zst }`)
- **zl:rcv**: Method has pointer receiver to zero-sized type (`func (*zst) f()`)
- **zl:mex**: Method expression receiver is pointer to zero-size type (`(*zst).Error(nil)`)

### Full Level

- **zl:add**: Address of zero-size variable (`&zsv`)
- **zl:ast**: Type assert to pointer to zero-size variable (`var a any; a.(*zst)`)
- **zl:typ**: Pointer to zero-sized type (`map[string]*zst`)
- **zl:arg**: Passing explicit nil as parameter pointing to zero-sized type (`func f(*zst); f(nil)`)
- **zl:par**: Function parameter points to zero-sized type (`func f(*zst)`)
- **zl:res**: Function has pointer result to zero-sized type (`func f() *zst`)

## Integration

See [`zerolint-golangci-plugin`](https://github.com/fillmore-labs/zerolint-golangci-plugin).

## Known Bugs

We are aware of a number of minor bugs in the analyzer's fixes. For example, it may sometimes prevent a type from
correctly implementing an interface or cause a non-pointer type to be checked for nil. The known bugs are low-risk and
easy to fix, as they result in a broken build or are obvious during a code review; none cause latent behavior changes.
Please report any additional problems you encounter.

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.
