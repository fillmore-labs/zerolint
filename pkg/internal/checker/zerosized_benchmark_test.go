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

package checker_test

import (
	"go/types"
	"strings"
	"sync"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/internal/checker"
)

/*
BenchmarkIsZeroSized compares the performance of the implemented
and an “optimized” implementations of IsZeroSized.

The optimized code has some micro-optimizations, while the trade-off is a significant
increase in complexity: The logic is split between two functions, with nested loops and breaks.

The original prioritizes clarity and maintainability over these micro-optimizations.
Most importantly, the performance difference is negligible in practice,
since the caller caches the result of the zero-size analysis for each type.

Run with `go test -bench=BenchmarkIsZeroSized ./pkg/internal/checker`.
*/
func BenchmarkIsZeroSized(b *testing.B) {
	testCases := []struct {
		name    string
		typeStr string
	}{
		// Simple cases
		{name: "Simple/ZeroStruct", typeStr: "struct{}"},
		{name: "Simple/ZeroArray", typeStr: "[0]int"},
		{name: "Simple/NonZeroStruct", typeStr: "struct{f int}"},
		{name: "Simple/NonZeroArray", typeStr: "[1]int"},
		{name: "Simple/Basic", typeStr: "string"},

		// Nested "linear" types (optimized implementation is optimized for this)
		{name: "Linear/Depth1", typeStr: "struct{f struct{}}"},
		{name: "Linear/Depth3", typeStr: "struct{f1 struct{f2 struct{f3 [0]int}}}"},
		{name: "Linear/Depth5", typeStr: "struct{f1 struct{f2 struct{f3 struct{f4 struct{f5 [0]int}}}}}"},
		{name: "Linear/NonZero", typeStr: "struct{f1 struct{f2 struct{f3 struct{f4 struct{f5 [1]int}}}}}"},

		// Nested "branched" types (optimized implementation must use stack)
		{name: "Branched/Simple", typeStr: "struct{f1 struct{}; f2 [0]int}"},
		{name: "Branched/Complex", typeStr: "struct{f1 struct{f1a [0]int; f1b struct{}}; f2 [0]string; f3 struct{}}"},
		{name: "Branched/Mixed", typeStr: "struct{f1 struct{f1a [0]int}; f2 [0]string; f3 struct{f3a int}}"},

		// Deeply nested branched type
		{name: "Deep/Branched", typeStr: strings.Repeat("struct{f1 ", 20) + "struct{}" + strings.Repeat("; f2 [0]int}", 20)},
	}

	for _, tc := range testCases {
		typ := parseType(b, tc.typeStr)

		b.Run(tc.name+"/Recursive", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result = ZeroSized(typ, 0)
			}
		})
		b.Run(tc.name+"/SemiOptimized", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result = isZeroSizedSemiOptimized(typ)
			}
		})
		b.Run(tc.name+"/Optimized", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result = isZeroSizedOptimized(typ)
			}
		})
		b.Run(tc.name+"/Plain", func(b *testing.B) {
			b.ReportAllocs()

			for b.Loop() {
				result = isZeroSizedPlain(typ)
			}
		})
	}
}

// result is a package-level variable to which the result of the benchmarked
// function is stored. This prevents the compiler from optimizing away the call.
var result bool //nolint:gochecknoglobals

//nolint:gochecknoglobals
var (
	typeCache = make(map[string]types.Type)
	typeMutex sync.Mutex
)

// parseType parses a string representation of a Go type and returns the corresponding
// types.Type, using a cache to avoid re-parsing.
func parseType(tb testing.TB, typeStr string) types.Type {
	tb.Helper()

	typeMutex.Lock()
	defer typeMutex.Unlock()

	if typ, ok := typeCache[typeStr]; ok {
		return typ
	}

	src := "package main\n\ntype V " + typeStr
	pkg := parseSource(tb, "main.go", src)

	typ := getType(tb, pkg, "V")
	typeCache[typeStr] = typ

	return typ
}
