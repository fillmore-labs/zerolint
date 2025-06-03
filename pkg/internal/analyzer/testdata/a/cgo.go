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

type _Ctype_struc struct{}

type _Ctype_struct_ptr *_Ctype_struc

func something() _Ctype_struct_ptr {
	return nil
}

func _Cfunc_something(e *struct{}) *struct{} {
	return e
}

var _cgoSomething *struct{}

var _ = _Cfunc_something(&struct{}{})
