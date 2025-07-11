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

package a

import "test/a/b"

var (
	_ b.Excluded1     // want "IS excluded"
	_ b.Excluded2     // want "IS excluded"
	_ b.Excluded2a    // want "IS excluded"
	_ b.NotExcluded1  // want "IS NOT excluded"
	_ b.NotExcluded2  // want "IS NOT excluded"
	_ b.NotExcluded2a // want "IS NOT excluded"

	_ Excluded1    // want "IS excluded"
	_ NotExcluded1 // want "IS NOT excluded"
)

//zerolint:exclude
type Excluded1 struct{}

// zerolint:
type NotExcluded1 struct{}
