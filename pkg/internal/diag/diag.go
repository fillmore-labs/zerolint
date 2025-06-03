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

package diag

import (
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/analysis"
)

// Diag provides helper functions for reporting and fixing pointers to zero-sized types.
type Diag struct {
	pass *analysis.Pass

	// Currently processed file, used by [Diag.Qualifier] for imports.
	CurrentFile *ast.File
}

// New creates and initializes a [Diag] instance using the provided [analysis.Pass].
func New(pass *analysis.Pass) *Diag {
	d := &Diag{}
	d.Prepare(pass)

	return d
}

// Prepare initializes the [Diag] with the provided [analysis.Pass], preparing for new analysis.
func (d *Diag) Prepare(pass *analysis.Pass) {
	d.pass = pass
}

// TypesInfo returns the type information for the current analysis pass.
func (d *Diag) TypesInfo() *types.Info {
	return d.pass.TypesInfo
}
