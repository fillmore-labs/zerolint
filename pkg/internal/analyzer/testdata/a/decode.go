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

import (
	"encoding/json"

	ghodssyaml "github.com/ghodss/yaml"
	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"
)

type decodeEmpty struct{}

func Decode() {
	json.Unmarshal(nil, &decodeEmpty{})
	ghodssyaml.Unmarshal(nil, &decodeEmpty{})
	k8syaml.Unmarshal(nil, &decodeEmpty{})
	k8syaml.UnmarshalStrict(nil, &decodeEmpty{})
	yamlv2.Unmarshal(nil, &decodeEmpty{})
	yamlv2.UnmarshalStrict(nil, &decodeEmpty{})
	yamlv3.Unmarshal(nil, &decodeEmpty{})

	json.NewDecoder(nil).Decode(&decodeEmpty{})
	yamlv2.NewDecoder(nil).Decode(&decodeEmpty{})
	yamlv3.NewDecoder(nil).Decode(&decodeEmpty{})

	wrap := func(x any) any { return x }
	json.Unmarshal(nil, wrap(&decodeEmpty{}))
}
