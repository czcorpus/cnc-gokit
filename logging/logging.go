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
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// GinMiddleware is a zerolog logging middleware for Gin.
// It is inspired by the original logging routine from the
// Gin project.
func GinMiddleware() gin.HandlerFunc {

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
		logEvent.
			Float64("latency", t0.Sub(start).Seconds()).
			Str("clientIP", ctx.ClientIP()).
			Str("method", ctx.Request.Method).
			Int("status", ctx.Writer.Status()).
			Str("errorMessage", ctx.Errors.ByType(gin.ErrorTypePrivate).String()).
			Int("bodySize", ctx.Writer.Size()).
			Str("path", path).
			Send()
	}
}
