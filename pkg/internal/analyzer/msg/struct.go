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

// Struct implements [Formatter] for struct fields, providing appropriate messages
// for diagnostics related to struct fields pointing to zero-sized types.
type Struct struct{}

// ZeroMsg returns a message for embedded pointer fields to zero-sized types.
func (Struct) ZeroMsg(typ types.Type, valueMethod bool) diag.CategorizedMessage {
	return Formatf(CatStructEmbedded, valueMethod, "embedded pointer to zero-sized type %q", typ)
}

// SingularMsg returns a message for a single named struct field pointing to a zero-sized type.
func (Struct) SingularMsg(typ types.Type, valueMethod bool, name string) diag.CategorizedMessage {
	return Formatf(CatStructField, valueMethod, "field %q points to zero-sized type %q", name, typ)
}

// PluralMsg returns a message for multiple struct fields pointing to a zero-sized type.
func (Struct) PluralMsg(typ types.Type, valueMethod bool, names string) diag.CategorizedMessage {
	return Formatf(CatStructField, valueMethod, "fields %s point to zero-sized type %q", names, typ)
}
