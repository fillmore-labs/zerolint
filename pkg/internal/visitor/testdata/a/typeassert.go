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

package a

type (
	empt  = struct{}
	pempt = *empt // want "\\(zl:dcl\\)"
)

func Assert() {
	var a any = &empt{} // want "\\(zl:add\\)"
	var b pempt         // want "\\(zl:var\\)"
	b = a.(pempt)       // want "\\(zl:art\\)"
	_ = (*empt)(b)      // want "\\(zl:cst\\)"
	_ = b

	var s any = ""
	_ = s.(string)
}
