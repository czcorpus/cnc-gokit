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

package cors

import "github.com/gin-gonic/gin"

func getRequestOrigin(ctx *gin.Context) string {
	currOrigin, ok := ctx.Request.Header["Origin"]
	if ok {
		return currOrigin[0]
	}
	return ""
}

// CORSMiddleware generates access-control-allow-origin/credentials
// headers so API services running on different domains can be
// accessed from a main (web) application without additional
// treatment (e.g. proxying)
func CORSMiddleware(allowedOrigins []string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var allowedOrigin string
		currOrigin := getRequestOrigin(ctx)
		for _, origin := range allowedOrigins {
			if currOrigin == origin {
				allowedOrigin = origin
				break
			}
		}
		if allowedOrigin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			ctx.Writer.Header().Set(
				"Access-Control-Allow-Headers",
				"Content-Type, Content-Length, Accept-Encoding, Authorization, Accept, Origin, Cache-Control, X-Requested-With",
			)
			ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		}
		if ctx.Request.Method == "OPTIONS" {
			ctx.AbortWithStatus(204)
			return
		}
		ctx.Next()
	}
}
