package cors_test

import (
	"net/http"

	"github.com/go-http-utils/cors"
)

func Example() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", cors.Handler(mux))
}

func ExampleOption_SetExposeHeaders() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", cors.Handler(mux,
		cors.SetExposeHeaders([]string{"ETag"})))
}

func ExampleOption_SetCredentials() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	http.ListenAndServe(":8080", cors.Handler(mux,
		cors.SetExposeHeaders([]string{"ETag"}), cors.SetCredentials(true)))
}

func ExampleOption_SetAllowOriginValidator() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		res.Write([]byte("Hello World"))
	})

	validator := func(_ *http.Request) string {
		return "*"
	}

	http.ListenAndServe(":8080", cors.Handler(mux,
		cors.SetAllowOriginValidator(validator)))
}
