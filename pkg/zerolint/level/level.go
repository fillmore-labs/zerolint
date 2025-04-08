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

package level

import (
	"errors"
	"strconv"
	"strings"
)

// LintLevel represents the thoroughness of detection.
type LintLevel int //nolint:recvcheck

const (
	// Default is the standard linting level that applies basic checks to the code.
	// It primarily flags comparisons of pointers to zero-size types
	// and embedding pointers in structures.
	// It is the default level used when no specific linting level is specified.
	Default LintLevel = iota

	// Extended finds more usages of pointers to zero-sized types.
	// This includes variables, struct fields and method receivers that are
	// pointers to zero-size types, casts, calls to `new`` and explicitly passing or returning nil.
	Extended

	// Full finds most usages of pointers to zero-sized types.
	// This level probably needs an 'excluded' file.
	// It includes everything from Extended, plus explicitly taken addresses of zero-sized variables
	// and function parameters and return types that are pointers to zero-size types.
	// This level is very strict, but might be useful with `-fix`.
	Full
)

// ErrUnknownLintLevel is an error returned when an unrecognized or invalid
// lint level is encountered.
var ErrUnknownLintLevel = errors.New("unknown lint level")

const (
	levelDefault  = "default"
	levelExtended = "extended"
	levelFull     = "full"
)

// UnmarshalText implements encoding.TextUnmarshaler.
func (l *LintLevel) UnmarshalText(text []byte) error {
	switch strings.ToLower(string(text)) {
	case levelDefault, "0":
		*l = Default

	case levelExtended, "1":
		*l = Extended

	case levelFull, "2":
		*l = Full

	default:
		return ErrUnknownLintLevel
	}

	return nil
}

// MarshalText implements encoding.TextMarshaler.
func (l LintLevel) MarshalText() ([]byte, error) {
	switch l {
	case Default:
		return []byte(levelDefault), nil

	case Extended:
		return []byte(levelExtended), nil

	case Full:
		return []byte(levelFull), nil

	default:
		return nil, ErrUnknownLintLevel
	}
}

// AtLeast determines if the current lint level is at least the provided lint level.
func (l LintLevel) AtLeast(m LintLevel) bool {
	return l >= m
}

// Below determines if the current lint level is lower than the provided lint level.
func (l LintLevel) Below(m LintLevel) bool {
	return l < m
}

// String returns the string representation of the LintLevel value. For unrecognized values,
// it returns their numeric value.
func (l LintLevel) String() string {
	switch l {
	case Default:
		return levelDefault

	case Extended:
		return levelExtended

	case Full:
		return levelFull

	default:
		return strconv.Itoa(int(l))
	}
}
