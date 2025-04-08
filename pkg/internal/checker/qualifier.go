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

package checker

import (
	"go/ast"
	"go/types"
	"strconv"
)

// Qualifier returns the package name or alias to use when referring to the given package
// in the context of the currently analyzed file (whose imports are in c.CurrentImports).
func (c *Checker) Qualifier(pkg *types.Package) string {
	if pkg == nil || pkg == c.pass.Pkg {
		return ""
	}

	for _, i := range c.CurrentImports {
		if qualifier, ok := importName(i, pkg); ok {
			return qualifier
		}
	}

	return pkg.Name()
}

// importName extracts the qualifier for a given package from a single import spec.
func importName(i *ast.ImportSpec, pkg *types.Package) (string, bool) {
	if path, err := strconv.Unquote(i.Path.Value); err == nil && pkg.Path() == path {
		if i.Name == nil {
			return pkg.Name(), true
		}

		switch i.Name.Name {
		case ".":
			return "", true

		case "_":

		default:
			return i.Name.Name, true
		}
	}

	return "", false
}
