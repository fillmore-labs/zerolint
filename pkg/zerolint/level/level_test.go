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

package level_test

import (
	"errors"
	"testing"

	. "fillmore-labs.com/zerolint/pkg/zerolint/level"
)

func TestLintLevel_UnmarshalText(t *testing.T) { //nolint:funlen
	t.Parallel()

	tests := []struct {
		name      string
		text      []byte
		wantLevel LintLevel
		wantErr   error
	}{
		{
			name:      "default lowercase",
			text:      []byte("default"),
			wantLevel: Default,
		},
		{
			name:      "default uppercase",
			text:      []byte("DEFAULT"),
			wantLevel: Default,
		},
		{
			name:      "default numeric",
			text:      []byte("0"),
			wantLevel: Default,
		},
		{
			name:      "extended lowercase",
			text:      []byte("extended"),
			wantLevel: Extended,
		},
		{
			name:      "extended uppercase",
			text:      []byte("EXTENDED"),
			wantLevel: Extended,
		},
		{
			name:      "extended numeric",
			text:      []byte("1"),
			wantLevel: Extended,
		},
		{
			name:      "full lowercase",
			text:      []byte("full"),
			wantLevel: Full,
		},
		{
			name:      "full uppercase",
			text:      []byte("FULL"),
			wantLevel: Full,
		},
		{
			name:      "full numeric",
			text:      []byte("2"),
			wantLevel: Full,
		},
		{
			name:    "unknown string",
			text:    []byte("unknown"),
			wantErr: ErrUnknownLintLevel,
		},
		{
			name:    "unknown numeric",
			text:    []byte("3"),
			wantErr: ErrUnknownLintLevel,
		},
		{
			name:    "empty string",
			text:    []byte(""),
			wantErr: ErrUnknownLintLevel,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var l LintLevel

			err := l.UnmarshalText(tt.text)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UnmarshalText() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil && l != tt.wantLevel {
				t.Errorf("UnmarshalText() level = %v, wantLevel %v", l, tt.wantLevel)
			}
		})
	}
}

func TestLintLevel_MarshalText(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		level   LintLevel
		want    []byte
		wantErr error
	}{
		{"Default", Default, []byte("default"), nil},
		{"Extended", Extended, []byte("extended"), nil},
		{"Full", Full, []byte("full"), nil},
		{"Unknown", LintLevel(99), nil, ErrUnknownLintLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := tt.level.MarshalText()
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("MarshalText() error = %v, wantErr %v", err, tt.wantErr)

				return
			}

			if string(got) != string(tt.want) {
				t.Errorf("MarshalText() got = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestLintLevel_AtLeast(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    LintLevel
		m    LintLevel
		want bool
	}{
		{"Default at least Default", Default, Default, true},
		{"Default at least Extended", Default, Extended, false},
		{"Extended at least Default", Extended, Default, true},
		{"Full at least Extended", Full, Extended, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.l.AtLeast(tt.m); got != tt.want {
				t.Errorf("%s AtLeast %s = %t, want %t", tt.l, tt.m, got, tt.want)
			}
		})
	}
}

func TestLintLevel_Below(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		l    LintLevel
		m    LintLevel
		want bool
	}{
		{"Default below Default", Default, Default, false},
		{"Default below Extended", Default, Extended, true},
		{"Extended below Default", Extended, Default, false},
		{"Full below Extended", Full, Extended, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.l.Below(tt.m); got != tt.want {
				t.Errorf("%s Below %s = %t, want %t", tt.l, tt.m, got, tt.want)
			}
		})
	}
}

func TestLintLevel_String(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		level LintLevel
		want  string
	}{
		{"Default", Default, "default"},
		{"Extended", Extended, "extended"},
		{"Full", Full, "full"},
		{"Unknown", LintLevel(42), "42"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.level.String(); got != tt.want {
				t.Errorf("LintLevel.String() = %s, want %s", got, tt.want)
			}
		})
	}
}
