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

package analyzer

import "go/ast"

// visitCallIdent processes encoding/json.Unmarshal (ignored as it requires pointer arguments),
// errors.Is and errors.As from the standard library or golang.org/x/exp/errors.
func (v *visitor) visitCallIdent(n *ast.CallExpr, fun *ast.Ident) bool { //nolint:funlen,cyclop
	if isCfunc(fun.Name) {
		return false
	}

	obj, ok := v.diag.TypesInfo().Uses[fun]
	if !ok { // should not happen
		v.diag.LogErrorf(fun, "Can't find identifier %q", fun.Name)

		return true
	}

	var path string
	if pkg := obj.Pkg(); pkg != nil {
		path = pkg.Path()
	}

	name := obj.Name()

	switch path {
	case "encoding/json", "sigs.k8s.io/yaml", "github.com/ghodss/yaml",
		"gopkg.in/yaml.v3", "gopkg.in/yaml.v2", "github.com/go-yaml/yaml":
		if name == "Unmarshal" {
			return false // Do not report pointers in json.Unmarshal(..., ...).
		}

		return true

	case "errors", "golang.org/x/exp/errors", "golang.org/x/xerrors", "github.com/pkg/errors":
		if len(n.Args) != 2 {
			return true // We are only interested in comparisons.
		}

		switch name {
		case "As":
			return false // Do not report pointers in errors.As(..., ...).

		case "Is":
			return v.visitCmp(n, n.Args[0], n.Args[1]) // Delegate analysis of errors.Is(..., ...) to visitCmp.

		default:
			return true
		}

	case "github.com/stretchr/testify/assert", "github.com/stretchr/testify/require":
		if len(n.Args) < 3 {
			return true // We are only interested in comparisons.
		}

		switch name { // assert.Equal does not compare object identity, but uses [reflect.DeepEqual].
		case "ErrorAs", "ErrorAsf", "NotErrorAs", "NotErrorAsf":
			return false // Do not report pointers in ErrorAs(t, ..., ...).

		case "ErrorIs", "ErrorIsf", "NotErrorIs", "NotErrorIsf":
			return v.visitCmp(n, n.Args[1], n.Args[2]) // Delegate analysis of ErrorIs(t, ..., ...) to visitCmp.

		default:
			return true
		}

	case "gotest.tools/v3/assert":
		if len(n.Args) < 3 {
			return true // gotest.tools comparison functions typically take at least t, expected, actual.
		}

		switch name {
		case "Equal", "ErrorIs":
			// Delegate analysis of assert.Equal(t, ..., ...) to visitCmp.
			return v.visitCmp(n, n.Args[1], n.Args[2])

		default:
			return true
		}

	default:
		return true
	}
}
