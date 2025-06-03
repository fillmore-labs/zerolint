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

// ComparisonMessage generates a diagnostic message for comparing two pointers to zero-sized types.
func ComparisonMessage(left, right types.Type, valueMethod bool) diag.CategorizedMessage {
	leftTypeString := types.TypeString(left, nil)
	rightTypeString := types.TypeString(right, nil)

	cat := CatComparison

	if leftTypeString == rightTypeString { // types.Identical ignores aliases
		return Formatf(cat, valueMethod, "comparison of pointers to zero-size type %q", leftTypeString)
	}

	return Formatf(cat, valueMethod,
		"comparison of pointers to zero-size types %q and %q", leftTypeString, rightTypeString)
}

// ComparisonMessagePointerInterface generates a diagnostic message for pointer-to-interface comparison.
func ComparisonMessagePointerInterface(
	elemOp, interfaceOp types.Type, valueMethod bool,
) diag.CategorizedMessage {
	elemTypeString := types.TypeString(elemOp, nil)
	interfaceTypeString := types.TypeString(interfaceOp, nil)

	if interfaceTypeString == "error" {
		return Formatf(CatComparisonError, valueMethod,
			"comparison of pointer to zero-size type %q with error interface", elemTypeString)
	}

	return Formatf(CatComparisonInterface, valueMethod,
		"comparison of pointer to zero-size type %q with interface of type %q", elemTypeString, interfaceTypeString)
}
