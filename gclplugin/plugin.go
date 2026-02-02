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

package gclplugin

import (
	"regexp"

	"github.com/golangci/plugin-module-register/register"
	"golang.org/x/tools/go/analysis"

	"fillmore-labs.com/zerolint/pkg/zerolint"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

func init() {
	register.Plugin(zerolint.Name, New)
}

// Settings are the linters settings.
type Settings struct {
	Excluded  []string         `json:"excluded,omitempty"`
	Level     *level.LintLevel `json:"level,omitempty"`
	Match     *regexp.Regexp   `json:"match,omitempty"`
	Generated *bool            `json:"generated,omitempty"`
}

// New creates a new [Plugin] instance with the given [Settings].
func New(settings any) (register.LinterPlugin, error) {
	s, err := register.DecodeSettings[Settings](settings)
	if err != nil {
		return nil, err
	}

	return Plugin{settings: s}, nil
}

// Plugin is a zerolint linter as a [register.LinterPlugin].
type Plugin struct {
	settings Settings
}

// BuildAnalyzers returns the [analysis.Analyzer]s for a zerolint run.
func (p Plugin) BuildAnalyzers() ([]*analysis.Analyzer, error) {
	var opts []zerolint.Option

	if p.settings.Level != nil {
		opts = append(opts, zerolint.WithLevel(*p.settings.Level))
	}

	if len(p.settings.Excluded) > 0 {
		opts = append(opts, zerolint.WithExcludes(p.settings.Excluded))
	}

	if p.settings.Match != nil {
		opts = append(opts, zerolint.WithRegex(p.settings.Match))
	}

	if p.settings.Generated != nil {
		opts = append(opts, zerolint.WithGenerated(*p.settings.Generated))
	}

	z := zerolint.New(opts...)

	return []*analysis.Analyzer{z}, nil
}

// GetLoadMode returns the golangci load mode.
func (Plugin) GetLoadMode() string {
	return register.LoadModeTypesInfo
}
