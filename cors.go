package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-http-utils/headers"
)

// Version is this package's version
const Version = "0.0.1"

// Handle wraps the http.Handler with CORS support.
func Handle(h http.Handler, opts ...Option) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		option := &options{
			exposeHeaders: defualtExposeHeaders,
			methods:       defualtMethods,
			allowOrigins:  []string{"*"},
		}

		for _, opt := range opts {
			opt(option)
		}

		passed := false
		origin := req.Header.Get(headers.Origin)

		if option.allowOriginValidator != nil {
			passed = option.allowOriginValidator(req)
		}

		if !passed && option.allowOrigins != nil {
			if allowAllOrigins(option.allowOrigins) {
				passed = true
			} else {
				passed = has(option.allowOrigins, origin)
			}
		}

		if !passed {
			h.ServeHTTP(res, req)
			return
		}

		resHeader := res.Header()

		if len(option.exposeHeaders) > 0 {
			resHeader.Set(headers.AccessControlExposeHeaders,
				strings.Join(option.exposeHeaders, ","))
		}

		if option.maxAge > 0 {
			resHeader.Set(headers.AccessControlMaxAge, strconv.Itoa(option.maxAge))
		}

		if option.credentials == true {
			resHeader.Set(headers.AccessControlAllowCredentials, "true")
		}

		if len(option.methods) > 0 {
			resHeader.Set(headers.AccessControlAllowMethods,
				strings.Join(option.methods, ","))
		}

		// TODO: Allow headers

		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusNoContent)
		} else {
			h.ServeHTTP(res, req)
		}
	})
}

func allowAllOrigins(allowOrigins []string) bool {
	return len(allowOrigins) == 1 && allowOrigins[0] == "*"
}
