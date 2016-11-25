package cors

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-http-utils/headers"
	"github.com/stretchr/testify/suite"
)

type OptionSuite struct {
	suite.Suite

	o *options
}

func (s *OptionSuite) SetupTest() {
	s.o = &options{allowOrigin: true}
}

func (s *OptionSuite) TestExpose() {
	err := SetExposeHeaders([]string{
		headers.Accept,
		headers.Accept,
		headers.Authorization,
	})(s.o)

	s.Nil(err)
	s.Equal([]string{headers.Accept, headers.Authorization}, s.o.exposeHeaders)
}

func (s *OptionSuite) TestAllowHeaders() {
	err := SetAllowHeaders([]string{
		headers.Accept,
		headers.Accept,
		headers.Authorization,
	})(s.o)

	s.Nil(err)
	s.Equal([]string{headers.Accept, headers.Authorization}, s.o.allowHeaders)
}

func (s *OptionSuite) TestMethods() {
	err := SetMethods([]string{
		http.MethodGet,
		http.MethodHead,
		"Get",
	})(s.o)

	s.Nil(err)
	s.Equal([]string{http.MethodGet, http.MethodHead}, s.o.methods)
}

func (s *OptionSuite) TestMaxAge() {
	err := SetMaxAge(-1)(s.o)

	s.Error(err, "cors: maxAge should > 0")

	err = SetMaxAge(600)(s.o)

	s.Nil(err)
	s.Equal(600, s.o.maxAge)

	err = SetMaxAge(601)(s.o)

	s.Nil(err)
	s.Equal(600, s.o.maxAge)
}

func (s *OptionSuite) TestCredentials() {
	err := SetCredentials(true)(s.o)

	s.Nil(err)
	s.Equal(true, s.o.credentials)
}

func (s *OptionSuite) TestAllowOrigin() {
	err := SetAllowOrigin(true)(s.o)

	s.Nil(err)
	s.Equal(true, s.o.allowOrigin)
}

func (s *OptionSuite) TestAllowOriginsValidator() {
	err := SetAllowOriginValidator(func(_ *http.Request) string {
		return "test"
	})(s.o)

	s.Nil(err)
	s.Equal("test",
		s.o.allowOriginValidator(httptest.NewRequest(http.MethodGet, "/", nil)))
}

func TestOption(t *testing.T) {
	suite.Run(t, new(OptionSuite))
}
