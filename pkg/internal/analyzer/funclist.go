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

package analyzer

import "fillmore-labs.com/zerolint/pkg/internal/typeutil"

type funcType int

const (
	funcNone funcType = iota
	funcDecode
	funcCmp0
	funcCmp1
)

// Since we have a lot of hardcoded libraries here, a check by signature might be a better heuristic.
var functions = map[typeutil.FuncName]funcType{ //nolint:gochecknoglobals
	{Path: "errors", Name: "Is"}:                                                                          funcCmp0,
	{Path: "golang.org/x/exp/errors", Name: "Is"}:                                                         funcCmp0,
	{Path: "golang.org/x/xerrors", Name: "Is"}:                                                            funcCmp0,
	{Path: "github.com/pkg/errors", Name: "Is"}:                                                           funcCmp0,
	{Path: "gotest.tools/v3/assert", Name: "Equal"}:                                                       funcCmp1,
	{Path: "gotest.tools/v3/assert", Name: "ErrorIs"}:                                                     funcCmp1,
	{Path: "gotest.tools/v3/assert/cmp", Name: "Equal"}:                                                   funcCmp0,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorIs"}:                                         funcCmp1,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorIsf"}:                                        funcCmp1,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorIs"}:                                      funcCmp1,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorIsf"}:                                     funcCmp1,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorIs"}:                                        funcCmp1,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorIsf"}:                                       funcCmp1,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorIs"}:                                     funcCmp1,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorIsf"}:                                    funcCmp1,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorIs", Ptr: true}:      funcCmp0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorIsf", Ptr: true}:     funcCmp0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorIs", Ptr: true}:   funcCmp0,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorIsf", Ptr: true}:  funcCmp0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorIs", Ptr: true}:     funcCmp0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorIsf", Ptr: true}:    funcCmp0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorIs", Ptr: true}:  funcCmp0,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorIsf", Ptr: true}: funcCmp0,
	{Path: "errors", Name: "As"}:                                                                          funcDecode,
	{Path: "golang.org/x/exp/errors", Name: "As"}:                                                         funcDecode,
	{Path: "golang.org/x/xerrors", Name: "As"}:                                                            funcDecode,
	{Path: "github.com/pkg/errors", Name: "As"}:                                                           funcDecode,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorAs"}:                                         funcDecode,
	{Path: "github.com/stretchr/testify/assert", Name: "ErrorAsf"}:                                        funcDecode,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorAs"}:                                      funcDecode,
	{Path: "github.com/stretchr/testify/assert", Name: "NotErrorAsf"}:                                     funcDecode,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorAs"}:                                        funcDecode,
	{Path: "github.com/stretchr/testify/require", Name: "ErrorAsf"}:                                       funcDecode,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorAs"}:                                     funcDecode,
	{Path: "github.com/stretchr/testify/require", Name: "NotErrorAsf"}:                                    funcDecode,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorAs", Ptr: true}:      funcDecode,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "ErrorAsf", Ptr: true}:     funcDecode,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorAs", Ptr: true}:   funcDecode,
	{Path: "github.com/stretchr/testify/assert", Receiver: "Assertions", Name: "NotErrorAsf", Ptr: true}:  funcDecode,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorAs", Ptr: true}:     funcDecode,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "ErrorAsf", Ptr: true}:    funcDecode,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorAs", Ptr: true}:  funcDecode,
	{Path: "github.com/stretchr/testify/require", Receiver: "Assertions", Name: "NotErrorAsf", Ptr: true}: funcDecode,
	{Path: "encoding/json", Name: "Unmarshal"}:                                                            funcDecode,
	{Path: "encoding/json", Receiver: "Decoder", Name: "Decode", Ptr: true}:                               funcDecode,
	{Path: "github.com/ghodss/yaml", Name: "Unmarshal"}:                                                   funcDecode,
	{Path: "gopkg.in/yaml.v2", Name: "Unmarshal"}:                                                         funcDecode,
	{Path: "gopkg.in/yaml.v2", Name: "UnmarshalStrict"}:                                                   funcDecode,
	{Path: "gopkg.in/yaml.v2", Receiver: "Decoder", Name: "Decode", Ptr: true}:                            funcDecode,
	{Path: "gopkg.in/yaml.v3", Name: "Unmarshal"}:                                                         funcDecode,
	{Path: "gopkg.in/yaml.v3", Receiver: "Decoder", Name: "Decode", Ptr: true}:                            funcDecode,
	{Path: "sigs.k8s.io/yaml", Name: "Unmarshal"}:                                                         funcDecode,
	{Path: "sigs.k8s.io/yaml", Name: "UnmarshalStrict"}:                                                   funcDecode,
}
