package middlewares_test

import (
	"github.com/sergeychur/avito_auto/internal/middlewares"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CorsGetToBody struct {
	uri         string
	method      string
	checkInside bool
	hostInHosts bool
}

func TestCreateCorsMiddleware(t *testing.T) {
	cases := []CorsGetToBody{
		{
			uri:         "http://localhost/uri",
			method:      http.MethodGet,
			checkInside: true,
			hostInHosts: true,
		},
		{
			uri:         "http://localhost/1",
			method:      http.MethodOptions,
			checkInside: false,
			hostInHosts: true,
		},
	}
	hosts := []string{"http://localhost", "https://newwordtrainer.ru"}
	function := middlewares.CreateCorsMiddleware(hosts)

	for _, testCase := range cases {
		here := false
		r := httptest.NewRequest(testCase.method, testCase.uri, nil)
		w := httptest.NewRecorder()
		h := http.HandlerFunc(
			func(w http.ResponseWriter, r *http.Request) {
				here = true
				w.WriteHeader(http.StatusOK)
			})
		function(h).ServeHTTP(w, r)
		if testCase.checkInside {
			if !here {
				t.Error("Failed cors check")
			}
		} else if testCase.hostInHosts && w.Result().StatusCode == http.StatusMethodNotAllowed {
			t.Errorf("Options don't work")
		}
	}
}