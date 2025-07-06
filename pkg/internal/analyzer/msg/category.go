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

package msg

import "fillmore-labs.com/zerolint/pkg/internal/diag"

//nolint:godot,revive
const (
	// keep-sorted start
	CatAddress             diag.Category = "add"
	CatArgumentNil         diag.Category = "arg"
	CatCast                diag.Category = "cst"
	CatCastNil             diag.Category = "nil"
	CatCastUnsafe          diag.Category = "cup"
	CatComparison          diag.Category = "cmp"
	CatComparisonError     diag.Category = "cme"
	CatComparisonInterface diag.Category = "cmi"
	CatDeref               diag.Category = "der"
	CatError               diag.Category = "err"
	CatMethodExpression    diag.Category = "mex"
	CatNew                 diag.Category = "new"
	CatParameter           diag.Category = "par"
	CatReceiver            diag.Category = "rcv"
	CatResult              diag.Category = "res"
	CatReturnNil           diag.Category = "ret"
	CatStarType            diag.Category = "typ"
	CatStructEmbedded      diag.Category = "emb"
	CatStructField         diag.Category = "fld"
	CatTypeAssert          diag.Category = "ast"
	CatTypeDeclaration     diag.Category = "dcl"
	CatVar                 diag.Category = "var"
	// keep-sorted end
)
