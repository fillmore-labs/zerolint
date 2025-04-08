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

package base

import (
	"bytes"
	"go/types"
	"strconv"
)

// qualifier returns the package name or alias to use when referring to the given package in the Current context.
func (v *Base) qualifier(pkg *types.Package) string {
	if pkg == nil || pkg == v.pass.Pkg {
		return ""
	}

	if v.Current == nil {
		return pkg.Name()
	}

	for _, i := range v.Current.Imports {
		path, err := strconv.Unquote(i.Path.Value)
		if err != nil {
			continue
		}

		if pkg.Path() == path {
			if i.Name == nil {
				return pkg.Name()
			}

			return i.Name.Name
		}
	}

	return pkg.Name()
}

// writeType returns the string representation of the zero value for a given type suitable for go source code.
// Returns false if a suitable literal cannot be determined.
func (v *Base) writeType(buf *bytes.Buffer, typ types.Type) bool {
	switch typ := typ.(type) {
	case *types.Named:
		types.WriteType(buf, typ, v.qualifier)

	case *types.Alias:
		types.WriteType(buf, typ, v.qualifier)

	case *types.Struct:
		types.WriteType(buf, typ, nil)

	case *types.Array:
		types.WriteType(buf, typ, nil)

	default:
		// types with non-zero sizes
		return false
	}

	return true
}
