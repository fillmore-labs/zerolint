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

//nolint:ireturn
package zerolint

import (
	"log"
	"log/slog"
	"regexp"

	"fillmore-labs.com/zerolint/pkg/internal/set"
	"fillmore-labs.com/zerolint/pkg/zerolint/level"
)

// options defines configurable parameters for the linter.
type options struct {
	level           level.LintLevel
	excludes        set.Set[string]
	generated       bool
	regex           *regexp.Regexp
	logger          *log.Logger
	zeroTrace       bool
	withFlags       bool
	excludeComments bool
}

// defaultOptions returns a [options] struct initialized with default values.
func defaultOptions() *options {
	return &options{ // Default options
		level:           level.Basic,
		logger:          log.Default(),
		excludeComments: true,
	}
}

// makeOptions returns a [options] struct with overriding [Option]s applied.
func makeOptions(opts Options) *options {
	o := defaultOptions()
	opts.apply(o)

	return o
}

// Option configures specific behavior of the zerolint [analysis.Analyzer].
type Option interface {
	LogValue() slog.Value
	key() string
	apply(opts *options)
}

// Options is a list of [Option] values that also satisfies the [Option] interface.
type Options []Option

// LogValue implements the [slog.LogValuer] interface.
func (o Options) LogValue() slog.Value {
	as := make([]slog.Attr, 0, len(o))
	for _, opt := range o {
		as = append(as, slog.Attr{Key: opt.key(), Value: opt.LogValue()})
	}

	return slog.GroupValue(as...)
}

func (o Options) apply(opts *options) {
	for _, opt := range o {
		opt.apply(opts)
	}
}

func (o Options) key() string {
	return "options"
}

// WithLevel is an [Option] to configure full linting.
func WithLevel(level level.LintLevel) Option {
	return levelOption{level: level}
}

type levelOption struct {
	level level.LintLevel
}

// LogValue implements the [slog.LogValuer] interface.
func (o levelOption) LogValue() slog.Value {
	return slog.StringValue(o.level.String())
}

func (o levelOption) key() string {
	return "level"
}

func (o levelOption) apply(opts *options) {
	opts.level = o.level
}

// WithExcludes is an [Option] to configure the excluded types.
func WithExcludes(excludes []string) Option {
	return excludesOption{excludes: excludes}
}

type excludesOption struct {
	excludes []string
}

// LogValue implements the [slog.LogValuer] interface.
func (o excludesOption) LogValue() slog.Value {
	return slog.AnyValue(o.excludes)
}

func (o excludesOption) key() string {
	return "excludes"
}

func (o excludesOption) apply(opts *options) {
	if opts.excludes == nil {
		opts.excludes = set.New[string]()
	}

	for _, exclude := range o.excludes {
		opts.excludes.Add(exclude)
	}
}

// WithZeroTrace is an [Option] to configure tracing of zero-sized types.
func WithZeroTrace(zeroTrace bool) Option {
	return zeroTraceOption{zeroTrace: zeroTrace}
}

type zeroTraceOption struct {
	zeroTrace bool
}

// LogValue implements the [slog.LogValuer] interface.
func (o zeroTraceOption) LogValue() slog.Value {
	return slog.BoolValue(o.zeroTrace)
}

func (o zeroTraceOption) key() string {
	return "zeroTrace"
}

func (o zeroTraceOption) apply(opts *options) {
	opts.zeroTrace = o.zeroTrace
}

// WithGenerated is an [Option] to configure linting of generated files.
func WithGenerated(generated bool) Option {
	return generatedOption{generated: generated}
}

type generatedOption struct {
	generated bool
}

// LogValue implements the [slog.LogValuer] interface.
func (o generatedOption) LogValue() slog.Value {
	return slog.BoolValue(o.generated)
}

func (o generatedOption) key() string {
	return "generated"
}

func (o generatedOption) apply(opts *options) {
	opts.generated = o.generated
}

// WithRegex is an [Option] to configure detecting only matching types.
func WithRegex(re *regexp.Regexp) Option {
	return reOption{re: re}
}

type reOption struct {
	re *regexp.Regexp
}

// LogValue implements the [slog.LogValuer] interface.
func (o reOption) LogValue() slog.Value {
	re := ""
	if o.re != nil {
		re = o.re.String()
	}

	return slog.StringValue(re)
}

func (o reOption) key() string {
	return "regex"
}

func (o reOption) apply(opts *options) {
	opts.regex = o.re
}

// WithLogger is an [Option] to configure the used logger.
func WithLogger(logger *log.Logger) Option {
	return loggerOption{logger: logger}
}

// LogValue implements the [slog.LogValuer] interface.
func (o loggerOption) LogValue() slog.Value {
	prefix := "<nil>"
	if o.logger != nil {
		prefix = o.logger.Prefix()
	}

	return slog.StringValue(prefix)
}

func (o loggerOption) key() string {
	return "logger"
}

type loggerOption struct {
	logger *log.Logger
}

func (o loggerOption) apply(opts *options) {
	opts.logger = o.logger
}

// WithExcludeComments is an [Option] to configure parsing of `zerolint:exclude` comments.
func WithExcludeComments(excludeComments bool) Option {
	return excludeCommentsOption{excludeComments: excludeComments}
}

type excludeCommentsOption struct {
	excludeComments bool
}

// LogValue implements the [slog.LogValuer] interface.
func (o excludeCommentsOption) LogValue() slog.Value {
	return slog.BoolValue(o.excludeComments)
}

func (o excludeCommentsOption) key() string {
	return "exclude-comments"
}

func (o excludeCommentsOption) apply(opts *options) {
	opts.excludeComments = o.excludeComments
}

// WithFlags is an [Option] to configure parsing of command-line flags.
// When enabled, command-line flags (e.g., -level, -excluded) will be parsed
// and will override any corresponding options set programmatically via other `With...` functions.
func WithFlags(flags bool) Option {
	return flagsOption{flags: flags}
}

type flagsOption struct {
	flags bool
}

// LogValue implements the [slog.LogValuer] interface.
func (o flagsOption) LogValue() slog.Value {
	return slog.BoolValue(o.flags)
}

func (o flagsOption) key() string {
	return "flags"
}

func (o flagsOption) apply(opts *options) {
	opts.withFlags = o.flags
}
