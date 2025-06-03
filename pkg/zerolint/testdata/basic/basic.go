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

package basic

import (
	"errors"
	"fmt"
)

type myError struct{}

func (*myError) Error() string { // want " \\(zl:err\\)$"
	return "my error"
}

type aliasError = myError

type noError struct{}

func (*noError) Error() string {
	return "no error"
}

var (
	ErrOne   = &myError{}
	ErrTwo   = new(aliasError)
	ErrThree = &noError{}
	ErrFour  = new(noError)
)

func Exported() {
	if errors.Is(nil, ErrOne) {
		fmt.Println("nil")
	}

	if errors.Is(func() error { // want " \\(zl:cme\\)$"
		return ErrOne
	}(), ErrTwo) {
		fmt.Println("equal")
	}

	if ErrOne == ErrTwo { // want " \\(zl:cmp\\)$"
		fmt.Println("equal")
	}

	if ErrOne != ErrTwo { // want " \\(zl:cmp\\)$"
		fmt.Println("not equal")
	}

	if ErrThree != ErrFour {
		fmt.Println("not equal")
	}
}
