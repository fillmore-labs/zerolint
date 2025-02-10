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

package analyzer

import (
	"log"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/visitor"
	"golang.org/x/tools/go/analysis"
)

// NewRun returns a configurable function for the Run field of [Analyzer].
func NewRun(opts ...Option) func(*analysis.Pass) (any, error) {
	option := options{
		Logger: log.Default(),
	}
	for _, opt := range opts {
		opt.apply(&option)
	}

	v := visitor.New(option)

	return func(pass *analysis.Pass) (any, error) {
		res, err := v.Run(pass)

		if option.ZeroTrace && v.HasDetected() {
			option.Logger.Printf("found zero-sized types in %q:\n", v.Pass.Pkg.Path())
			for name := range v.AllDetected() {
				option.Logger.Printf("- %s\n", name)
			}
		}

		return res, err
	}
}

// options defines configurable parameters for the linter.
type options = visitor.Options

// Option defines configurations for [NewRun].
type Option interface {
	apply(opts *options)
}

// WithLogger is an [Option] to configure the used logger.
func WithLogger(logger *log.Logger) Option { //nolint:ireturn
	return loggerOption{logger: logger}
}

type loggerOption struct {
	logger *log.Logger
}

func (o loggerOption) apply(opts *options) {
	opts.Logger = o.logger
}

// WithExcludes is an [Option] to configure the excluded types.
func WithExcludes(excludes []string) Option { //nolint:ireturn
	return excludesOption{excludes: excludes}
}

type excludesOption struct {
	excludes []string
}

func (o excludesOption) apply(opts *options) {
	if opts.Excludes == nil {
		opts.Excludes = set.New[string]()
	}
	for _, exclude := range o.excludes {
		opts.Excludes.Insert(exclude)
	}
}

// WithZeroTrace is an [Option] to configure tracing of zero sized types.
func WithZeroTrace(zeroTrace bool) Option { //nolint:ireturn
	return zeroTraceOption{zeroTrace: zeroTrace}
}

type zeroTraceOption struct {
	zeroTrace bool
}

func (o zeroTraceOption) apply(opts *options) {
	opts.ZeroTrace = o.zeroTrace
}

// WithFull is an [Option] to configure full linting.
func WithFull(full bool) Option { //nolint:ireturn
	return fullOption{full: full}
}

type fullOption struct {
	full bool
}

func (o fullOption) apply(opts *options) {
	opts.Full = o.full
}

// WithGenerated is an [Option] to configure linting of generated files.
func WithGenerated(generated bool) Option { //nolint:ireturn
	return generatedOption{generated: generated}
}

type generatedOption struct {
	generated bool
}

func (o generatedOption) apply(opts *options) {
	opts.Generated = o.generated
}
