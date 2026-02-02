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
	"errors"
	"testing"

	pkgerrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	experrors "golang.org/x/exp/errors"
	"golang.org/x/xerrors"
	gotest "gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"
)

type cmpErr struct{}

func (*cmpErr) Error() string { return "" } // want " \\(zl:err\\)$"

func wrap(err, target error) (error, error) {
	return err, target
}

func wrapt(t assert.TestingT, err, target error) (assert.TestingT, error, error) {
	return t, err, target
}

func TestCmp(t *testing.T) {
	errors.Is(&cmpErr{}, &cmpErr{})    // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	experrors.Is(&cmpErr{}, &cmpErr{}) // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	xerrors.Is(&cmpErr{}, &cmpErr{})   // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	pkgerrors.Is(&cmpErr{}, &cmpErr{}) // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"

	errors.As(nil, &cmpErr{})
	experrors.As(nil, &cmpErr{})
	xerrors.As(nil, &cmpErr{})
	pkgerrors.As(nil, &cmpErr{})

	assert.ErrorIs(t, &cmpErr{}, &cmpErr{})          // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	assert.ErrorIsf(t, &cmpErr{}, &cmpErr{}, "")     // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	assert.NotErrorIs(t, &cmpErr{}, &cmpErr{})       // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	assert.NotErrorIsf(t, &cmpErr{}, &cmpErr{}, "")  // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	require.ErrorIs(t, &cmpErr{}, &cmpErr{})         // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	require.ErrorIsf(t, &cmpErr{}, &cmpErr{}, "")    // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	require.NotErrorIs(t, &cmpErr{}, &cmpErr{})      // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	require.NotErrorIsf(t, &cmpErr{}, &cmpErr{}, "") // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"

	assert.ErrorAs(t, nil, &cmpErr{})
	assert.ErrorAsf(t, nil, &cmpErr{}, "")
	assert.NotErrorAs(t, nil, &cmpErr{})
	assert.NotErrorAsf(t, nil, &cmpErr{}, "")
	require.ErrorAs(t, nil, &cmpErr{})
	require.ErrorAsf(t, nil, &cmpErr{}, "")
	require.NotErrorAs(t, nil, &cmpErr{})
	require.NotErrorAsf(t, nil, &cmpErr{}, "")

	var s suite.Suite
	r := s.Require()

	s.ErrorIs(&cmpErr{}, &cmpErr{})         // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	s.ErrorIsf(&cmpErr{}, &cmpErr{}, "")    // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	s.NotErrorIs(&cmpErr{}, &cmpErr{})      // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	s.NotErrorIsf(&cmpErr{}, &cmpErr{}, "") // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	r.ErrorIs(&cmpErr{}, &cmpErr{})         // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	r.ErrorIsf(&cmpErr{}, &cmpErr{}, "")    // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	r.NotErrorIs(&cmpErr{}, &cmpErr{})      // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	r.NotErrorIsf(&cmpErr{}, &cmpErr{}, "") // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"

	s.ErrorAs(nil, &cmpErr{})
	s.ErrorAsf(nil, &cmpErr{}, "")
	s.NotErrorAs(nil, &cmpErr{})
	s.NotErrorAsf(nil, &cmpErr{}, "")
	r.ErrorAs(nil, &cmpErr{})
	r.ErrorAsf(nil, &cmpErr{}, "")
	r.NotErrorAs(nil, &cmpErr{})
	r.NotErrorAsf(nil, &cmpErr{}, "")

	gotest.Equal(t, &cmpErr{}, &cmpErr{})   // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"
	gotest.ErrorIs(t, &cmpErr{}, &cmpErr{}) // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"

	cmp.Equal(&cmpErr{}, &cmpErr{}) // want " \\(zl:cmp\\)$" "\\(zl:add\\)$" "\\(zl:add\\)$"

	assert.ErrorIs(wrapt(t, &cmpErr{}, &cmpErr{})) // want "\\(zl:add\\)$" "\\(zl:add\\)$"
	s.ErrorIs(wrap(&cmpErr{}, &cmpErr{}))          // want "\\(zl:add\\)$" "\\(zl:add\\)$"
}
