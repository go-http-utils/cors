package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-http-utils/headers"
)

// Version is this package's version
const Version = "0.1.0"

// Handler wraps the http.Handler with CORS support.
func Handler(h http.Handler, opts ...Option) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
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

		origin := ""

		if option.allowOrigin {
			origin = req.Header.Get(headers.Origin)

			if origin == "" {
				origin = "*"
			}
		} else if option.allowOriginValidator != nil {
			origin = option.allowOriginValidator(req)
		}

		if origin == "" {
			return
		}

		resHeader := res.Header()

		resHeader.Set(headers.AccessControlAllowOrigin, origin)

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

		var allowHeaders string

		if len(option.allowHeaders) > 0 {
			allowHeaders = strings.Join(option.allowHeaders, ",")
		} else {
			allowHeaders = req.Header.Get(headers.AccessControlRequestHeaders)
		}

		resHeader.Set(headers.AccessControlAllowHeaders, allowHeaders)

		if req.Method == http.MethodOptions {
			res.WriteHeader(http.StatusNoContent)
			return
		}

		h.ServeHTTP(res, req)
	})
}
