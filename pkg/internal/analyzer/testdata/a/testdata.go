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

import "fmt"

func Exported() {
	var x [0]string
	var y [0]string

	_ = (new)(struct{ _ [0]func() }) // want " \\(zl:new\\)$"

	type empty struct{}
	_ = new(empty) // want " \\(zl:new\\)$"

	xp, yp := &x, &y // want " \\(zl:add\\)$" " \\(zl:add\\)$"

	_ = *xp // want " \\(zl:der\\)$"

	if xp == yp { // want " \\(zl:cmp\\)$"
		fmt.Println("equal")
	}

	if xp != yp { // want " \\(zl:cmp\\)$"
		fmt.Println("not equal")
	}

	_, _ = any(xp).((*[0]string)) // want " \\(zl:art\\)$"

	switch any(xp).(type) {
	case (*[0]string): // want " \\(zl:art\\)$"
	case string:
	}
}
