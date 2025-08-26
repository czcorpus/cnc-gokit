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

// These unit tests were generated with assistance from Claude Code.

import (
	"bytes"
	"io"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel_IsDebugMode(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected bool
	}{
		{"debug", true},
		{"info", false},
		{"warning", false},
		{"warn", false},
		{"error", false},
		{"invalid", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			result := tt.level.IsDebugMode()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLogLevel_IsValid(t *testing.T) {
	tests := []struct {
		level    LogLevel
		expected bool
	}{
		{"debug", true},
		{"info", true},
		{"warning", true},
		{"warn", true},
		{"error", true},
		{"invalid", false},
		{"DEBUG", false}, // case sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(string(tt.level), func(t *testing.T) {
			result := tt.level.IsValid()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestLoggingConf_validate(t *testing.T) {
	tests := []struct {
		name     string
		conf     LoggingConf
		expected LoggingConf
	}{
		{
			name: "all defaults applied",
			conf: LoggingConf{},
			expected: LoggingConf{
				MaxFileSize: dfltLoggingMaxFileSize,
				MaxFiles:    dfltLoggingMaxFiles,
				MaxAgeDays:  dfltLoggingMaxAgeDays,
			},
		},
		{
			name: "partial defaults applied",
			conf: LoggingConf{
				MaxFileSize: 100,
				MaxFiles:    0,
				MaxAgeDays:  0,
			},
			expected: LoggingConf{
				MaxFileSize: 100,
				MaxFiles:    dfltLoggingMaxFiles,
				MaxAgeDays:  dfltLoggingMaxAgeDays,
			},
		},
		{
			name: "no defaults needed",
			conf: LoggingConf{
				MaxFileSize: 200,
				MaxFiles:    5,
				MaxAgeDays:  30,
			},
			expected: LoggingConf{
				MaxFileSize: 200,
				MaxFiles:    5,
				MaxAgeDays:  30,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output to avoid noise in test output
			var buf bytes.Buffer
			originalLogger := log.Logger
			log.Logger = zerolog.New(&buf)
			defer func() { log.Logger = originalLogger }()

			err := tt.conf.validate()
			assert.NoError(t, err)
			assert.Equal(t, tt.expected.MaxFileSize, tt.conf.MaxFileSize)
			assert.Equal(t, tt.expected.MaxFiles, tt.conf.MaxFiles)
			assert.Equal(t, tt.expected.MaxAgeDays, tt.conf.MaxAgeDays)
		})
	}
}

func TestSetupLogging(t *testing.T) {
	// Save original logger and global level
	originalLogger := log.Logger
	originalLevel := zerolog.GlobalLevel()
	defer func() {
		log.Logger = originalLogger
		zerolog.SetGlobalLevel(originalLevel)
	}()

	t.Run("console logging setup", func(t *testing.T) {
		var buf bytes.Buffer
		log.Logger = zerolog.New(&buf)

		conf := LoggingConf{
			Level:       "info",
			MaxFileSize: 100,
			MaxFiles:    2,
			MaxAgeDays:  14,
		}

		SetupLogging(conf)
		assert.Equal(t, zerolog.InfoLevel, zerolog.GlobalLevel())
	})

	t.Run("file logging setup", func(t *testing.T) {
		var buf bytes.Buffer
		log.Logger = zerolog.New(&buf)

		tempDir := t.TempDir()
		logFile := filepath.Join(tempDir, "test.log")

		conf := LoggingConf{
			Path:        logFile,
			Level:       "debug",
			MaxFileSize: 200,
			MaxFiles:    3,
			MaxAgeDays:  21,
		}

		SetupLogging(conf)
		assert.Equal(t, zerolog.DebugLevel, zerolog.GlobalLevel())
	})
}

func TestReqMatchesMonitoring(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name          string
		clientIP      string
		userAgent     string
		monitoringIPs []string
		uaSubstr      string
		expected      bool
	}{
		{
			name:          "both IP and UA match",
			clientIP:      "192.168.1.100",
			userAgent:     "monitoring-bot/1.0",
			monitoringIPs: []string{"192.168.1.100", "10.0.0.1"},
			uaSubstr:      "monitoring",
			expected:      true,
		},
		{
			name:          "IP matches but UA doesn't",
			clientIP:      "192.168.1.100",
			userAgent:     "regular-browser/1.0",
			monitoringIPs: []string{"192.168.1.100"},
			uaSubstr:      "monitoring",
			expected:      false,
		},
		{
			name:          "UA matches but IP doesn't",
			clientIP:      "192.168.1.101",
			userAgent:     "monitoring-bot/1.0",
			monitoringIPs: []string{"192.168.1.100"},
			uaSubstr:      "monitoring",
			expected:      false,
		},
		{
			name:          "case insensitive UA matching",
			clientIP:      "192.168.1.100",
			userAgent:     "MONITORING-BOT/1.0",
			monitoringIPs: []string{"192.168.1.100"},
			uaSubstr:      "monitoring",
			expected:      true,
		},
		{
			name:          "neither matches",
			clientIP:      "192.168.1.101",
			userAgent:     "regular-browser/1.0",
			monitoringIPs: []string{"192.168.1.100"},
			uaSubstr:      "monitoring",
			expected:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "/", nil)
			req.Header.Set("User-Agent", tt.userAgent)
			req.RemoteAddr = tt.clientIP + ":12345"

			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req

			conf := middlewareConf{
				monitoringIPs:      tt.monitoringIPs,
				monitoringUASubstr: tt.uaSubstr,
			}

			result := reqMatchesMonitoring(ctx, conf)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGinMiddlewareOptions(t *testing.T) {
	t.Run("GinMiddlewareWithMonitoringIPs", func(t *testing.T) {
		ips := []string{"192.168.1.100", "10.0.0.1"}
		option := GinMiddlewareWithMonitoringIPs(ips)

		conf := middlewareConf{}
		option(&conf)

		assert.Equal(t, ips, conf.monitoringIPs)
	})

	t.Run("GinMiddlewareWithMonitoringUASubstr", func(t *testing.T) {
		substr := "monitoring-bot"
		option := GinMiddlewareWithMonitoringUASubstr(substr)

		conf := middlewareConf{}
		option(&conf)

		assert.Equal(t, substr, conf.monitoringUASubstr)
	})
}

func TestGinMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Setup test logger
	var buf bytes.Buffer
	originalLogger := log.Logger
	log.Logger = zerolog.New(&buf)
	defer func() { log.Logger = originalLogger }()

	t.Run("successful request logging", func(t *testing.T) {
		buf.Reset()

		router := gin.New()
		router.Use(GinMiddleware())
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test?param=value", nil)
		req.Header.Set("User-Agent", "test-agent/1.0")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 200, w.Code)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "GET")
		assert.Contains(t, logOutput, "/test?param=value")
		assert.Contains(t, logOutput, "test-agent/1.0")
		assert.Contains(t, logOutput, "200")
	})

	t.Run("error request logging", func(t *testing.T) {
		buf.Reset()

		router := gin.New()
		router.Use(GinMiddleware())
		router.GET("/error", func(c *gin.Context) {
			c.JSON(500, gin.H{"error": "internal error"})
		})

		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, 500, w.Code)

		logOutput := buf.String()
		assert.Contains(t, logOutput, "GET")
		assert.Contains(t, logOutput, "/error")
		assert.Contains(t, logOutput, "500")
		// Error level logging should be used for 5xx status codes
		assert.Contains(t, logOutput, `"level":"error"`)
	})

	t.Run("monitoring detection", func(t *testing.T) {
		buf.Reset()

		router := gin.New()
		router.Use(GinMiddleware(
			GinMiddlewareWithMonitoringIPs([]string{"192.168.1.100"}),
			GinMiddlewareWithMonitoringUASubstr("monitoring"),
		))
		router.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("User-Agent", "monitoring-bot/1.0")
		req.RemoteAddr = "192.168.1.100:12345"
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		logOutput := buf.String()
		assert.Contains(t, logOutput, `"isMonitoring":true`)
	})

	t.Run("custom log entries", func(t *testing.T) {
		buf.Reset()

		router := gin.New()
		router.Use(GinMiddleware())
		router.GET("/test", func(c *gin.Context) {
			AddCustomEntry(c, "customKey", "customValue")
			AddCustomEntry(c, "requestID", "12345")
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/test", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		logOutput := buf.String()
		assert.Contains(t, logOutput, `"customKey":"customValue"`)
		assert.Contains(t, logOutput, `"requestID":"12345"`)
	})

	t.Run("gin errors handling", func(t *testing.T) {
		buf.Reset()

		router := gin.New()
		router.Use(GinMiddleware())
		router.GET("/error", func(c *gin.Context) {
			c.Error(gin.Error{Err: io.EOF, Type: gin.ErrorTypePrivate})
			c.JSON(200, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest("GET", "/error", nil)
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		logOutput := buf.String()
		assert.Contains(t, logOutput, `"errorMessage"`)
		assert.Contains(t, logOutput, "EOF")
	})
}

func TestAddCustomEntry(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	AddCustomEntry(ctx, "testKey", "testValue")
	AddCustomEntry(ctx, "numberKey", 42)

	// Check that the values are stored with the correct prefix
	value, exists := ctx.Get(logEventPrefix + "testKey")
	assert.True(t, exists)
	assert.Equal(t, "testValue", value)

	numberValue, exists := ctx.Get(logEventPrefix + "numberKey")
	assert.True(t, exists)
	assert.Equal(t, 42, numberValue)
}

func TestAddLogEvent_Deprecated(t *testing.T) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Test deprecated function still works
	AddLogEvent(ctx, "deprecatedKey", "deprecatedValue")

	value, exists := ctx.Get(logEventPrefix + "deprecatedKey")
	assert.True(t, exists)
	assert.Equal(t, "deprecatedValue", value)
}

func TestConstants(t *testing.T) {
	assert.Equal(t, "logEvent_", logEventPrefix)
	assert.Equal(t, 500, dfltLoggingMaxFileSize)
	assert.Equal(t, 3, dfltLoggingMaxFiles)
	assert.Equal(t, 28, dfltLoggingMaxAgeDays)
}

func TestLevelMapping(t *testing.T) {
	expectedMappings := map[LogLevel]zerolog.Level{
		"debug":   zerolog.DebugLevel,
		"info":    zerolog.InfoLevel,
		"warning": zerolog.WarnLevel,
		"warn":    zerolog.WarnLevel,
		"error":   zerolog.ErrorLevel,
	}

	for level, expectedZLevel := range expectedMappings {
		zLevel, exists := levelMapping[level]
		assert.True(t, exists, "Level %s should exist in mapping", level)
		assert.Equal(t, expectedZLevel, zLevel, "Level %s should map to %v", level, expectedZLevel)
	}
}
