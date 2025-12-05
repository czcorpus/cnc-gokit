// Copyright 2023 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2023 Institute of the Czech National Corpus,
//                Faculty of Arts, Charles University
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package logging

import (
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	logEventPrefix         = "logEvent_"
	dfltLoggingMaxFileSize = 500
	dfltLoggingMaxFiles    = 3
	dfltLoggingMaxAgeDays  = 28
)

var (
	levelMapping = map[LogLevel]zerolog.Level{
		"debug":   zerolog.DebugLevel,
		"info":    zerolog.InfoLevel,
		"warning": zerolog.WarnLevel,
		"warn":    zerolog.WarnLevel,
		"error":   zerolog.ErrorLevel,
	}
)

type LogLevel string

func (ll LogLevel) IsDebugMode() bool {
	return ll == "debug"
}

func (ll LogLevel) IsValid() bool {
	_, ok := levelMapping[ll]
	return ok
}

type LoggingConf struct {

	// Path specifies logging file path. If empty, then stderr
	// is used.
	Path string `json:"path"`

	// Level specifies level of logging (debug, info, warning (warn), error)
	Level LogLevel `json:"level"`

	MaxFileSize int `json:"maxFileSize"`
	MaxFiles    int `json:"maxFiles"`
	MaxAgeDays  int `json:"maxAgeDays"`
}

func (conf *LoggingConf) validate() error {
	if conf.MaxFileSize == 0 {
		conf.MaxFileSize = dfltLoggingMaxFileSize
		log.Warn().Msgf("missing logging.maxFileSize, setting %d", dfltLoggingMaxFileSize)
	}
	if conf.MaxFiles == 0 {
		conf.MaxFiles = dfltLoggingMaxFiles
		log.Warn().Msgf("missing logging.maxFiles, setting %d", dfltLoggingMaxFiles)
	}
	if conf.MaxAgeDays == 0 {
		conf.MaxAgeDays = dfltLoggingMaxAgeDays
		log.Warn().Msgf("missing logging.maxAgeDays, setting %d", dfltLoggingMaxAgeDays)
	}
	return nil
}

// SetupLogging is a common setup for different
// CNC HTTP services.
func SetupLogging(conf LoggingConf) {
	if err := conf.validate(); err != nil {
		log.Fatal().Err(err).Msgf("invalid config")
	}
	lev, ok := levelMapping[conf.Level]
	if !ok {
		log.Fatal().Msgf("Invalid logging level: %s", conf.Level)
	}
	zerolog.SetGlobalLevel(lev)
	if conf.Path != "" {
		log.Logger = log.Output(&lumberjack.Logger{
			Filename:   conf.Path,
			MaxSize:    conf.MaxFileSize,
			MaxBackups: conf.MaxFiles,
			MaxAge:     conf.MaxAgeDays,
			Compress:   false,
		})

	} else {
		log.Logger = log.Output(
			zerolog.ConsoleWriter{
				Out:        os.Stderr,
				TimeFormat: time.RFC3339,
			},
		)
	}
}

// -------

type middlewareConf struct {
	monitoringIPs      []string
	monitoringUASubstr string
}

func reqMatchesMonitoring(ctx *gin.Context, conf middlewareConf) bool {
	var ipMatch bool
	for _, ip := range conf.monitoringIPs {
		if ip == ctx.ClientIP() {
			ipMatch = true
			break
		}
	}
	uaMatch := strings.Contains(
		strings.ToLower(ctx.Request.UserAgent()),
		strings.ToLower(conf.monitoringUASubstr),
	)
	return ipMatch && uaMatch
}

func GinMiddlewareWithMonitoringIPs(ips []string) func(conf *middlewareConf) {
	return func(conf *middlewareConf) {
		conf.monitoringIPs = ips
	}
}

func GinMiddlewareWithMonitoringUASubstr(substr string) func(conf *middlewareConf) {
	return func(conf *middlewareConf) {
		conf.monitoringUASubstr = substr
	}
}

// GinMiddleware is a zerolog logging middleware for Gin.
// It is inspired by the original logging routine from the
// Gin project.
func GinMiddleware(opts ...func(conf *middlewareConf)) gin.HandlerFunc {

	var conf middlewareConf

	for _, opt := range opts {
		opt(&conf)
	}

	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		if ctx.Request.URL.RawQuery != "" {
			path = path + "?" + ctx.Request.URL.RawQuery
		}

		ctx.Next()

		var logEvent *zerolog.Event
		if ctx.Writer.Status() >= 500 {
			logEvent = log.Error()

		} else {
			logEvent = log.Info()
		}
		t0 := time.Now()
		errs := ctx.Errors.ByType(gin.ErrorTypePrivate)
		if len(errs) > 0 {
			logEvent = logEvent.Str("errorMessage", errs.String())
		}
		logEvent = logEvent.
			Float64("latency", t0.Sub(start).Seconds()).
			Str("clientIP", ctx.ClientIP()).
			Str("method", ctx.Request.Method).
			Int("status", ctx.Writer.Status()).
			Int("bodySize", ctx.Writer.Size()).
			Str("userAgent", ctx.Request.UserAgent()).
			Str("path", path)

		if reqMatchesMonitoring(ctx, conf) {
			logEvent = logEvent.Bool("isMonitoring", true)
		}

		for k, v := range ctx.Keys {
			kStr, ok := k.(string)
			if !ok {
				continue
			}
			if strings.HasPrefix(kStr, logEventPrefix) {
				logEvent = logEvent.Any(kStr[len(logEventPrefix):], v)
			}
		}
		logEvent.Send()
	}
}

func AddCustomEntry(ctx *gin.Context, key string, value any) {
	ctx.Set(logEventPrefix+key, value)
}

// AddLogEvent
// Deprecated: please use `AddCustomEntry` instead
func AddLogEvent(ctx *gin.Context, key string, value any) {
	AddCustomEntry(ctx, key, value)
}
