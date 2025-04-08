// Copyright 2024 Oliver Eikemeier. All Rights Reserved.
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

// change to `go:generate go tool stringer ...` when go1.23 is dropped
//go:generate go run golang.org/x/tools/cmd/stringer -linecomment -type category

package visitor

// category represents an internal category code used to categorize
// different types of issues found by the linter.
type category int

//nolint:godot
const (
	catNone category = iota // unk
	// keep-sorted start
	catAddress             // add
	catArgumentNil         // arg
	catCast                // cst
	catCastNil             // nil
	catComparison          // cmp
	catComparisonError     // cme
	catComparisonInterface // cmi
	catDeref               // der
	catError               // err
	catMethodExpression    // mex
	catNew                 // new
	catParameter           // par
	catReceiver            // rcv
	catResult              // res
	catReturnNil           // ret
	catStarType            // typ
	catStructEmbedded      // emb
	catStructField         // fld
	catTypeAssert          // art
	catTypeDeclaration     // dcl
	catVar                 // var
	// keep-sorted end
)
