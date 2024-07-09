# Fillmore Labs zerolint

[![Go Reference](https://pkg.go.dev/badge/fillmore-labs.com/zerolint.svg)](https://pkg.go.dev/fillmore-labs.com/zerolint)
[![Test](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml/badge.svg?branch=main)](https://github.com/fillmore-labs/zerolint/actions/workflows/test.yml)
[![Coverage](https://codecov.io/gh/fillmore-labs/zerolint/branch/main/graph/badge.svg?token=TUE1BL1QZV)](https://codecov.io/gh/fillmore-labs/zerolint)
[![Maintainability](https://api.codeclimate.com/v1/badges/baf50ad423cc30ff7790/maintainability)](https://codeclimate.com/github/fillmore-labs/zerolint/maintainability)
[![Go Report Card](https://goreportcard.com/badge/fillmore-labs.com/zerolint)](https://goreportcard.com/report/fillmore-labs.com/zerolint)
[![License](https://img.shields.io/github/license/fillmore-labs/zerolint)](https://www.apache.org/licenses/LICENSE-2.0)

The `zerolint` linter checks usage patterns of pointers to zero-sized variables in Go.

## Usage

```console
go install fillmore-labs.com/zerolint@latest
```

Usage: `zerolint [-flag] [package]`

Flags:

- **-c** int display offending line with this many lines of context (default -1)
- **-basic** basic analysis only
- **-excluded** `<filename>` read excluded types from this file
- **-zerotrace** trace found zero-sized types
- **-fix** apply all suggested fixes
