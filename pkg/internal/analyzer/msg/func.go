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

// Param implements [Formatter] for function parameters, providing appropriate diagnostic
// messages for parameters that are pointers to zero-sized types.
type Param struct{}

// ZeroMsg returns a message for unnamed function parameters pointing to zero-sized types.
func (Param) ZeroMsg(typ types.Type, valueMethod bool) diag.CategorizedMessage {
	return Formatf(CatParameter, valueMethod, "function has pointer parameter to zero-sized type %q", typ)
}

// SingularMsg returns a message for a single named function parameter pointing to a zero-sized type.
func (Param) SingularMsg(typ types.Type, valueMethod bool, name string) diag.CategorizedMessage {
	return Formatf(CatParameter, valueMethod, "function parameter %q points to zero-sized type %q", name, typ)
}

// PluralMsg returns a message for multiple function parameters pointing to a zero-sized type.
func (Param) PluralMsg(typ types.Type, valueMethod bool, names string) diag.CategorizedMessage {
	return Formatf(CatParameter, valueMethod, "function parameters %s point to zero-sized type %q", names, typ)
}

// Result implements [Formatter] for function results, providing appropriate diagnostic
// messages for results that are pointers to zero-sized types.
type Result struct{}

// ZeroMsg returns a message for unnamed function results pointing to zero-sized types.
func (Result) ZeroMsg(typ types.Type, valueMethod bool) diag.CategorizedMessage {
	return Formatf(CatResult, valueMethod, "function has pointer result to zero-sized type %q", typ)
}

// SingularMsg returns a message for a single named function result pointing to a zero-sized type.
func (Result) SingularMsg(typ types.Type, valueMethod bool, name string) diag.CategorizedMessage {
	return Formatf(CatResult, valueMethod, "function result %q points to zero-sized type %q", name, typ)
}

// PluralMsg returns a message for multiple function results pointing to a zero-sized type.
func (Result) PluralMsg(typ types.Type, valueMethod bool, names string) diag.CategorizedMessage {
	return Formatf(CatResult, valueMethod, "function results %s point to zero-sized type %q", names, typ)
}
