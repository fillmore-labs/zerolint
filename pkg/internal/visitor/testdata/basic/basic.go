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

package basic

import (
	"encoding/json"
	"errors"
	"fmt"

	. "golang.org/x/exp/errors"
)

type myError struct{}

type myErrorPtr = *myError // want "\\(zl:dcl\\+\\)"

func (myErrorPtr) Error() string { // want "\\(zl:err\\+\\)"
	return "my error"
}

func (myErrorPtr) String() string {
	return "my error description"
}

func (myError) GoString() string {
	return "go string"
}

type embeddedError struct {
	*myError // want "\\(zl:emb\\+\\)"
	_        *myError
}

var (
	ErrOne = &myError{}
	ErrTwo = new(myError)
)

func Exported() {
	if errors.Is(nil, ErrOne) {
		fmt.Println("nil")
	}

	if errors.Is(func() error { // want "\\(zl:cme\\+\\)"
		return ErrOne
	}(), ErrTwo) {
		fmt.Println("equal")
	}

	if ErrOne == ErrTwo { // want "\\(zl:cmp\\+\\)"
		fmt.Println("equal")
	}

	if ErrOne != ErrTwo { // want "\\(zl:cmp\\+\\)"
		fmt.Println("not equal")
	}

	if (Is)(ErrOne, ErrTwo) { // want "\\(zl:cmp\\+\\)"
		fmt.Println("equal")
	}

	empty := struct{}{}
	json := (*json.Decoder)(nil)
	_ = json.Decode(&empty)

	_ = (*myError).Error(ErrOne)
}
