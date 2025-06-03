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

package e

//zerolint:exclude
type Excluded struct{}

type NotExcluded struct{}

type Excluded2 = NotExcluded

//zerolint:exclude
var (
	_ Excluded  // want "IS excluded"
	_ Excluded2 // want "IS excluded"
)

//zerolint:exclude
var (
	e NotExcluded     // want " \\(zl:com\\)$" "IS NOT excluded"
	_ = NotExcluded{} // want " \\(zl:com\\)$"
	_ struct{}        // want " \\(zl:com\\)$"
)

var (
	//zerolint:exclude
	_ NotExcluded // want " \\(zl:com\\)$" "IS NOT excluded"
	_ NotExcluded //zerolint:exclude // want " \\(zl:com\\)$" "IS NOT excluded"
)

type (
	//zerolint:exclude
	NotExcluded2 struct{} // want " \\(zl:com\\)$"

	NotExcluded3 struct{} //zerolint:exclude // want " \\(zl:com\\)$"
)
