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

import "encoding/json"

type (
	myDecoder1 struct {
		json.Decoder
	}

	myDecoder2 struct {
		*json.Decoder
	}

	myDecoder3 = *json.Decoder
	myDecoder4 = json.Decoder
)

func IgnoreJson() {
	empty := struct{}{}

	_ = json.Unmarshal(nil, &empty)
	_ = (*json.Decoder)(nil).Decode(&empty)
	_ = json.NewDecoder(nil).Decode(&empty)
	_ = (*myDecoder1)(nil).Decode(&empty)
	_ = (&myDecoder1{}).Decode(&empty)
	_ = myDecoder2{}.Decode(&empty)
	_ = (myDecoder3)(nil).Decode(&empty)

	_ = (*json.Decoder).Decode(nil, &empty)
	_ = (myDecoder3).Decode(nil, &empty)
	_ = (*myDecoder4).Decode(nil, &empty)

	_, _ = json.Marshal(&empty) // want " \\(zl:add\\)$"
}
