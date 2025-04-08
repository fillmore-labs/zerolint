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

package visitor

// category represents an internal category code used to categorize
// different types of issues found by the linter.
type category int

const (
	catNone category = iota
	catComparison
	catInterfaceComparison
	catError
	catReceiver
	catParameter
	catResult
	catStructEmbedded
	catStructField
	catVar
	catStarType
	catDeref
	catAddress
	catNew
	catCast
	catTypeAssert
	catMethodExpression
	catNilReturn
	catNilParameter
)

// String returns a short string representation of the category code (e.g., "cmp", "err").
func (d category) String() string { //nolint:cyclop
	switch d {
	case catComparison:
		return "cmp"

	case catInterfaceComparison:
		return "cmi"

	case catError:
		return "err"

	case catReceiver:
		return "rcv"

	case catParameter:
		return "arg"

	case catResult:
		return "res"

	case catStructEmbedded:
		return "emb"

	case catStructField:
		return "fld"

	case catVar:
		return "var"

	case catStarType:
		return "typ"

	case catDeref:
		return "der"

	case catAddress:
		return "add"

	case catNew:
		return "new"

	case catCast:
		return "cst"

	case catTypeAssert:
		return "art"

	case catMethodExpression:
		return "mex"

	case catNilReturn:
		return "ret"

	case catNilParameter:
		return "par"

	case catNone:
	}

	return "unk"
}
