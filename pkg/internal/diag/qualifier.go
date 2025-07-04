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
	"strconv"
)

// Qualifier holds the current package and imports for [Qualifier.Qualifier].
type Qualifier struct {
	Pkg         *types.Package
	Imports     []*ast.ImportSpec
	NeedsImport bool
}

// Qualifier returns the package name or alias to use when referring to the given package
// in the context of the currently analyzed file (whose imports are in q.Imports).
func (q *Qualifier) Qualifier(pkg *types.Package) string {
	if pkg == nil || pkg == q.Pkg {
		return ""
	}

	for _, i := range q.Imports {
		if qualifier, ok := importName(i, pkg); ok {
			return qualifier
		}
	}

	q.NeedsImport = true

	// Not imported - default to quoted path, breaking the build in an informative way.
	return strconv.Quote(pkg.Path())
}

// importName extracts the qualifier for a given package from a single import spec.
func importName(i *ast.ImportSpec, pkg *types.Package) (string, bool) {
	if path, err := strconv.Unquote(i.Path.Value); err == nil && pkg.Path() == path {
		if i.Name == nil { // Standard import: import "fmt"
			return pkg.Name(), true
		}

		switch alias := i.Name.Name; alias {
		case ".": // Dot import: import . "fmt"
			return "", true

		case "_": // Blank import: import _ "fmt"
		// This doesn't make the package name available for qualification.

		default: // Aliased import: import f "fmt"
			return alias, true
		}
	}

	return "", false
}
