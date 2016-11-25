package cors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-http-utils/headers"
)

// Option returns a configuration func to configurate the
// CORS middleware.
type Option func(*options) error

type options struct {
	allowOrigin          bool
	allowOriginValidator func(*http.Request) string
	allowHeaders         []string
	exposeHeaders        []string
	maxAge               int
	credentials          bool
	methods              []string
}

// SetExposeHeaders configures the Access-Control-Expose-Headers CORS header.
func SetExposeHeaders(exposes []string) Option {
	return func(o *options) error {
		o.exposeHeaders = []string{}

		for _, header := range exposes {
			normalized := headers.Normalize(header)

			if has(o.exposeHeaders, normalized) {
				continue
			}

			o.exposeHeaders = append(o.exposeHeaders, normalized)
		}

		return nil
	}
}

// SetAllowHeaders configures the Access-Control-Allow-Headers CORS header.
func SetAllowHeaders(allows []string) Option {
	return func(o *options) error {
		o.allowHeaders = []string{}

		for _, header := range allows {
			normalized := headers.Normalize(header)

			if has(o.allowHeaders, normalized) {
				continue
			}

			o.allowHeaders = append(o.allowHeaders, normalized)
		}

		return nil
	}
}

// SetMethods configures the Access-Control-Allow-Methods CORS header.
func SetMethods(methods []string) Option {
	return func(o *options) error {
		o.methods = []string{}

		for _, method := range methods {
			normalized := strings.ToUpper(strings.TrimSpace(method))

			if has(o.methods, normalized) {
				continue
			}

			o.methods = append(o.methods, normalized)
		}

		return nil
	}
}

// SetMaxAge configures the Access-Control-Max-Age CORS header (in seconds).
// If the maxAge > 600 (10 minutes), then it will be reset to 600 (10 minutes).
func SetMaxAge(maxAge int) Option {
	return func(o *options) error {
		if maxAge < 0 {
			return fmt.Errorf("cors: maxAge should > 0")
		}

		if maxAge > 600 {
			maxAge = 600
		}

		o.maxAge = maxAge

		return nil
	}
}

// SetCredentials configures the Access-Control-Allow-Credentials CORS header.
func SetCredentials(credentials bool) Option {
	return func(o *options) error {
		o.credentials = credentials

		return nil
	}
}

// SetAllowOrigin configures the Access-Control-Allow-Origin CORS header.
// Set to true to reflect the request origin. Set to false to disable
// CORS.
func SetAllowOrigin(allow bool) Option {
	return func(o *options) error {
		o.allowOrigin = allow

		return nil
	}
}

// SetAllowOriginValidator configures the Access-Control-Allow-Origin CORS
// header by run the validator function `func(*http.Request) string`.
// The validator function accpets an `*http.Request` argument and return
// the Access-Control-Allow-Origin value.
func SetAllowOriginValidator(validator func(*http.Request) string) Option {
	return func(o *options) error {
		o.allowOriginValidator = validator

		return nil
	}
}

func has(hs []string, h string) bool {
	for _, e := range hs {
		if e == h {
			return true
		}
	}

	return false
}
