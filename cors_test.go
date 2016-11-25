package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/suite"
)

type CorsSuite struct {
	suite.Suite
}

func (s *CorsSuite) TestDefaultAllowOrigin() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("*", res.Header.Get(headers.AccessControlAllowOrigin))
}

func (s *CorsSuite) TestReflectAllowOrigin() {
	origin := "helloworld.org"

	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	s.Nil(err)

	req.Header.Set(headers.Origin, origin)

	res, err := sendRequest(req)

	s.Nil(err)
	s.Equal(origin, res.Header.Get(headers.AccessControlAllowOrigin))
}

func (s *CorsSuite) TestValidatorAllowOrigin() {
	origin := "helloworld.org"
	validator := func(req *http.Request) string {
		return origin
	}

	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetAllowOriginValidator(validator), SetAllowOrigin(false)))

	server := httptest.NewServer(mux)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	s.Nil(err)

	req.Header.Set(headers.Origin, "test.org")

	res, err := sendRequest(req)

	s.Nil(err)
	s.Equal(origin, res.Header.Get(headers.AccessControlAllowOrigin))
}

func (s *CorsSuite) TestNotAllowOrigin() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetAllowOrigin(false)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("", res.Header.Get(headers.AccessControlAllowOrigin))
	s.Equal("", res.Header.Get(headers.AccessControlAllowMethods))
}

func (s *CorsSuite) TestEmptyExpose() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("", res.Header.Get(headers.AccessControlExposeHeaders))
}

func (s *CorsSuite) TestExpose() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetExposeHeaders([]string{headers.AcceptRanges, headers.AcceptDatetime})))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal(headers.AcceptRanges+","+headers.AcceptDatetime,
		res.Header.Get(headers.AccessControlExposeHeaders))
}

func (s *CorsSuite) TestMaxAge() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetMaxAge(600)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("600", res.Header.Get(headers.AccessControlMaxAge))
}

func (s *CorsSuite) TestDefualtMethods() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("GET,HEAD,PUT,POST,DELETE",
		res.Header.Get(headers.AccessControlAllowMethods))
}

func (s *CorsSuite) TestMethods() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetMethods([]string{http.MethodHead, http.MethodTrace})))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("HEAD,TRACE", res.Header.Get(headers.AccessControlAllowMethods))
}

func (s *CorsSuite) TestCredentials() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetCredentials(true)))

	server := httptest.NewServer(mux)

	res, err := http.Get(server.URL + "/")

	s.Nil(err)
	s.Equal(http.StatusOK, res.StatusCode)
	s.Equal("true", res.Header.Get(headers.AccessControlAllowCredentials))
}

func (s *CorsSuite) TestDefualtAllowHeader() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	s.Nil(err)

	req.Header.Set(headers.AccessControlRequestHeaders, "FOO-BAR")

	res, err := sendRequest(req)

	s.Nil(err)
	s.Equal("FOO-BAR", res.Header.Get(headers.AccessControlAllowHeaders))
}

func (s *CorsSuite) TestAllowHeader() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc),
		SetAllowHeaders([]string{"FOO", "BAR"})))

	server := httptest.NewServer(mux)

	req, err := http.NewRequest(http.MethodGet, server.URL+"/", nil)

	s.Nil(err)

	req.Header.Set(headers.AccessControlRequestHeaders, "FOO-BAR")

	res, err := sendRequest(req)

	s.Nil(err)
	s.Equal("Foo,Bar", res.Header.Get(headers.AccessControlAllowHeaders))
}

func (s *CorsSuite) TestOptionRequest() {
	mux := http.NewServeMux()
	mux.Handle("/", Handler(http.HandlerFunc(helloHandlerFunc)))

	server := httptest.NewServer(mux)

	req, err := http.NewRequest(http.MethodOptions, server.URL+"/", nil)

	s.Nil(err)

	res, err := sendRequest(req)

	s.Nil(err)
	s.Equal(http.StatusNoContent, res.StatusCode)
}

func TestCors(t *testing.T) {
	suite.Run(t, new(CorsSuite))
}

func helloHandlerFunc(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(http.StatusOK)

	res.Write([]byte("Hello World"))
}

func sendRequest(req *http.Request) (*http.Response, error) {
	cli := &http.Client{}

	return cli.Do(req)
}
