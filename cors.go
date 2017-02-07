package cors

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-http-utils/headers"
)

// Version is this package's version
const Version = "1.1.0"

// Handler wraps the http.Handler with CORS support.
func Handler(h http.Handler, opts ...Option) http.Handler {
	option := &options{
		allowOrigin: true,
		methods: []string{
			http.MethodGet,
			http.MethodHead,
			http.MethodPut,
			http.MethodPost,
			http.MethodDelete,
		},
	}

	for _, opt := range opts {
		opt(option)
	}

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		resHeader := res.Header()
		resHeader.Add(headers.Vary, headers.Origin)
		origin := req.Header.Get(headers.Origin)

		// Not a CORS request.
		if origin == "" {
			h.ServeHTTP(res, req)

			return
		}

		allowOrigin := ""

		if option.allowOriginValidator != nil {
			allowOrigin = option.allowOriginValidator(req)
		} else {
			allowOrigin = req.Header.Get(headers.Origin)
		}

		if allowOrigin == "" {
			res.WriteHeader(http.StatusForbidden)
			res.Write([]byte(fmt.Sprintf("Invalid origin %v", origin)))
			return
		}

		if allowOrigin != "*" && option.credentials {
			// When responding to a credentialed request, server must specify a
			// domain, and cannot use wild carding.
			// See *important note* in https://developer.mozilla.org/en-US/docs/Web/HTTP/Access_control_CORS#Requests_with_credentials .
			resHeader.Set(headers.AccessControlAllowCredentials, "true")
		}

		resHeader.Set(headers.AccessControlAllowOrigin, allowOrigin)

		// Preflighted requests
		if req.Method == http.MethodOptions {
			resHeader.Add(headers.Vary, headers.AccessControlAllowMethods)
			resHeader.Add(headers.Vary, headers.AccessControlAllowHeaders)
			requestMethod := req.Header.Get(headers.AccessControlRequestMethod)

			if requestMethod == "" {
				resHeader.Del(headers.AccessControlAllowOrigin)
				resHeader.Del(headers.AccessControlAllowCredentials)

				res.WriteHeader(http.StatusForbidden)
				res.Write([]byte("Invalid preflighted request, missing Access-Control-Request-Method header"))

				return
			}

			if len(option.methods) > 0 {
				resHeader.Set(headers.AccessControlAllowMethods,
					strings.Join(option.methods, ","))
			}

			var allowHeaders string

			if len(option.allowHeaders) > 0 {
				allowHeaders = strings.Join(option.allowHeaders, ",")
			} else {
				allowHeaders = req.Header.Get(headers.AccessControlRequestHeaders)
			}

			resHeader.Set(headers.AccessControlAllowHeaders, allowHeaders)

			if option.maxAge > 0 {
				resHeader.Set(headers.AccessControlMaxAge, strconv.Itoa(option.maxAge))
			}

			res.WriteHeader(http.StatusNoContent)

			return
		}

		if len(option.exposeHeaders) > 0 {
			resHeader.Set(headers.AccessControlExposeHeaders,
				strings.Join(option.exposeHeaders, ","))
		}

		h.ServeHTTP(res, req)
	})
}
