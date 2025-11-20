// Copyright 2024-2025 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

/*
`zerolint` is a Go static analysis tool (linter) that detects unnecessary or potentially incorrect usage of pointers to
zero-sized types.

Pointers to zero-size types (ZSTs) can be problematic:
  - They carry very little information.
  - Two pointers to distinct zero-size variables may or may not compare equal. This can lead to subtle bugs.
  - The pointers themselves are not zero-sized and might waste memory in data structures, on the stack and in the CPU
    cache.

Usage:

	zerolint [flags] [package ...]

The flags are:

	-level level
		analysis level (basic, extended, full) (default Basic)
	-fix
		apply all suggested fixes
	-excluded file
		read excluded types from this file
	-match regex
		only check types matching this regex, useful with -fix
	-generated
		check generated files
	-c int
		display offending line with this many lines of context (default -1)
	-zerotrace
		trace found zero-sized types

# Examples

To check the current package for basic problems with pointers to zero-sized types:

	zerolint .

To fix all issues across packages, using an exclude file:

	zerolint -level=full -excluded=excludes.txt -fix ./...
*/
package main
