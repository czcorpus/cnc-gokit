// Copyright 2026 Tomas Machalek <tomas.machalek@gmail.com>
// Copyright 2026 Institute of the Czech National Corpus,
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

package mcptools

import (
	"context"
	"net/http"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/rs/zerolog/log"
)

type remoteAddrCtxKey struct{}

// clientIPFromRequest extracts the caller's IP address, preferring the
// X-Forwarded-For/X-Real-Ip headers set by a reverse proxy over the raw
// connection address.
func clientIPFromRequest(r *http.Request) string {
	if fwd := r.Header.Get("X-Forwarded-For"); fwd != "" {
		return strings.TrimSpace(strings.Split(fwd, ",")[0])
	}
	if realIP := r.Header.Get("X-Real-Ip"); realIP != "" {
		return realIP
	}
	return r.RemoteAddr
}

// WithClientIPContext returns an HTTPContextFunc suitable for
// server.WithHTTPContextFunc which makes the caller's IP address available
// to tool call hooks via RemoteAddrFromContext.
func WithClientIPContext(ctx context.Context, r *http.Request) context.Context {
	return context.WithValue(ctx, remoteAddrCtxKey{}, clientIPFromRequest(r))
}

// RemoteAddrFromContext returns the caller's IP address as stored by
// WithClientIPContext, or "" when running in stdio mode (no HTTP request).
func RemoteAddrFromContext(ctx context.Context) string {
	addr, _ := ctx.Value(remoteAddrCtxKey{}).(string)
	return addr
}

// DefaultLoggingHooks builds server hooks that log every tool call together with
// the caller's remote address (HTTP mode only) and, when available, the
// identity of the connecting MCP client (e.g. the LLM/agent calling the tool).
func DefaultLoggingHooks() *server.Hooks {
	hooks := &server.Hooks{}
	hooks.AddBeforeCallTool(func(ctx context.Context, id any, message *mcp.CallToolRequest) {
		evt := log.Info().
			Str("tool", message.Params.Name).
			Str("remoteAddr", RemoteAddrFromContext(ctx))

		if session := server.ClientSessionFromContext(ctx); session != nil {
			evt = evt.Str("sessionId", session.SessionID())
			if sessionWithInfo, ok := session.(server.SessionWithClientInfo); ok {
				clientInfo := sessionWithInfo.GetClientInfo()
				evt = evt.Str("clientName", clientInfo.Name).Str("clientVersion", clientInfo.Version)
			}
		}
		if ua := message.Header.Get("User-Agent"); ua != "" {
			evt = evt.Str("userAgent", ua)
		}
		evt.Msg("MCP tool called")
	})
	return hooks
}
