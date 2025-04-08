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

package zerolint

import (
	"log"
	"regexp"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/internal/visitor"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// options defines configurable parameters for the linter.
type options struct {
	visitor.Options
	logger    *log.Logger
	zeroTrace bool
	flags     bool
}

// Option configures specific behavior of the zerolint [analysis.Analyzer].
type Option interface {
	apply(opts *options)
}

// Options is a list of [Option] values that also satisfies the [Option] interface.
type Options []Option

func (o Options) apply(opts *options) {
	for _, opt := range o {
		opt.apply(opts)
	}
}

// WithLevel is an [Option] to configure full linting.
func WithLevel(level level.LintLevel) Option { //nolint:ireturn
	return levelOption{level: level}
}

type levelOption struct {
	level level.LintLevel
}

func (o levelOption) apply(opts *options) {
	opts.Level = o.level
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

// WithZeroTrace is an [Option] to configure tracing of zero-sized types.
func WithZeroTrace(zeroTrace bool) Option { //nolint:ireturn
	return zeroTraceOption{zeroTrace: zeroTrace}
}

type zeroTraceOption struct {
	zeroTrace bool
}

func (o zeroTraceOption) apply(opts *options) {
	opts.zeroTrace = o.zeroTrace
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

// WithRegex is an [Option] to configure detecting only matching types.
func WithRegex(re *regexp.Regexp) Option { //nolint:ireturn
	return reOption{re: re}
}

type reOption struct {
	re *regexp.Regexp
}

func (o reOption) apply(opts *options) {
	opts.Regex = o.re
}

// WithLogger is an [Option] to configure the used logger.
func WithLogger(logger *log.Logger) Option { //nolint:ireturn
	return loggerOption{logger: logger}
}

type loggerOption struct {
	logger *log.Logger
}

func (o loggerOption) apply(opts *options) {
	opts.logger = o.logger
}

// WithFlags is an [Option] to configure parsing of command-line flags.
// When enabled, command-line flags (e.g., -level, -excluded) will be parsed
// and will override any corresponding options set programmatically via other `With...` functions.
// If the "-V" flag is specified on the command line, it prints the program version and exits.
func WithFlags(flags bool) Option { //nolint:ireturn
	return flagsOption{flags: flags}
}

type flagsOption struct {
	flags bool
}

func (o flagsOption) apply(opts *options) {
	opts.flags = o.flags
}
