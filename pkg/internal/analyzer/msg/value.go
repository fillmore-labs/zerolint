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

import (
	"go/types"

	"fillmore-labs.com/zerolint/pkg/internal/diag"
)

// Value implements [Formatter] for value declarations, providing appropriate messages
// for diagnostics related to variables pointing to zero-sized types.
type Value struct{}

// ZeroMsg would return a message for an unnamed variable pointing to a zero-sized type. This is not possible in Go.
func (Value) ZeroMsg(typ types.Type, valueMethod bool) diag.CategorizedMessage {
	return Formatf(CatVar, valueMethod, "variable is pointer to zero-sized type %q", typ)
}

// SingularMsg returns a message for a single named variable pointing to a zero-sized type.
func (Value) SingularMsg(typ types.Type, valueMethod bool, name string) diag.CategorizedMessage {
	return Formatf(CatVar, valueMethod, "variable %q is pointer to zero-sized type %q", name, typ)
}

// PluralMsg returns a message for multiple variables pointing to a zero-sized type.
func (Value) PluralMsg(typ types.Type, valueMethod bool, names string) diag.CategorizedMessage {
	return Formatf(CatVar, valueMethod, "variables %s point to zero-sized type %q", names, typ)
}
