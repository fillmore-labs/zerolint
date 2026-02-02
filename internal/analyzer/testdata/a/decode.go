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
	yamlv2i "github.com/go-yaml/yaml"
	goccyaml "github.com/goccy/go-yaml"
	goyaml2 "go.yaml.in/yaml/v2"
	goyaml3 "go.yaml.in/yaml/v3"
	goyaml4 "go.yaml.in/yaml/v4"
	yamlv2 "gopkg.in/yaml.v2"
	yamlv3 "gopkg.in/yaml.v3"
	k8syaml "sigs.k8s.io/yaml"
)

type decodeEmpty struct{}

func Decode() {
	json.Unmarshal(nil, &decodeEmpty{})
	goccyaml.Unmarshal(nil, &decodeEmpty{})
	goccyaml.UnmarshalContext(nil, nil, &decodeEmpty{})
	goccyaml.UnmarshalWithOptions(nil, &decodeEmpty{})
	ghodssyaml.Unmarshal(nil, &decodeEmpty{})
	k8syaml.Unmarshal(nil, &decodeEmpty{})
	k8syaml.UnmarshalStrict(nil, &decodeEmpty{})
	yamlv2i.Unmarshal(nil, &decodeEmpty{})
	yamlv2i.UnmarshalStrict(nil, &decodeEmpty{})
	yamlv2.Unmarshal(nil, &decodeEmpty{})
	yamlv2.UnmarshalStrict(nil, &decodeEmpty{})
	yamlv3.Unmarshal(nil, &decodeEmpty{})
	goyaml2.Unmarshal(nil, &decodeEmpty{})
	goyaml3.Unmarshal(nil, &decodeEmpty{})
	goyaml4.Unmarshal(nil, &decodeEmpty{})

	json.NewDecoder(nil).Decode(&decodeEmpty{})
	goccyaml.NewDecoder(nil).Decode(&decodeEmpty{})
	goccyaml.NewDecoder(nil).DecodeContext(nil, &decodeEmpty{})
	goccyaml.NewDecoder(nil).DecodeFromNode(nil, &decodeEmpty{})
	goccyaml.NewDecoder(nil).DecodeFromNodeContext(nil, nil, &decodeEmpty{})
	yamlv2i.NewDecoder(nil).Decode(&decodeEmpty{})
	yamlv2.NewDecoder(nil).Decode(&decodeEmpty{})
	yamlv3.NewDecoder(nil).Decode(&decodeEmpty{})
	goyaml2.NewDecoder(nil).Decode(&decodeEmpty{})
	goyaml3.NewDecoder(nil).Decode(&decodeEmpty{})
	goyaml4.NewDecoder(nil).Decode(&decodeEmpty{})

	wrap := func(x any) any { return x }
	json.Unmarshal(nil, wrap(&decodeEmpty{}))
}
