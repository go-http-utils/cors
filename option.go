package cors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/go-http-utils/headers"
)

var (
	defualtExposeHeaders = []string{
		headers.CacheControl,
		headers.ContentLanguage,
		headers.ContentType,
		headers.Expires,
		headers.LastModified,
		headers.Pragma,
	}

	defualtMethods = []string{
		http.MethodGet,
		http.MethodHead,
		http.MethodPost,
	}
)

// Option returns a configuration func to configurate the
// CORS middleware.
type Option func(*options) error

type options struct {
	allowOrigins         []string
	allowOriginValidator func(*http.Request) bool
	allowHeaders         []string
	exposeHeaders        []string
	maxAge               int
	credentials          bool
	methods              []string
}

// SetExposeHeaders configures the Access-Control-Expose-Headers CORS header.
func SetExposeHeaders(exposes []string) Option {
	return func(o *options) error {
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
			o.allowHeaders = append(o.allowHeaders, headers.Normalize(header))
		}

		return nil
	}
}

// SetMethods configures the Access-Control-Allow-Methods CORS header.
func SetMethods(methods []string) Option {
	return func(o *options) error {
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

// SetAllowOrigins configures the Access-Control-Allow-Origin CORS header.
func SetAllowOrigins(orgins []string) Option {
	return func(o *options) error {
		o.allowOrigins = orgins

		return nil
	}
}

// SetAllowOriginValidator configures the Access-Control-Allow-Origin CORS
// header by run the validator function `func(*http.Request) bool`.
// If an orgin can pass the validator, then it won't need be in orgins that
// setted by SetAllowOrigins.
func SetAllowOriginValidator(validator func(*http.Request) bool) Option {
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
