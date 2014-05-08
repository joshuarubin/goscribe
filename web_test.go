package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnvironment(t *testing.T) {
	Convey("TelAPI Account SID should be set", t, func() {
		So(os.Getenv("TELAPI_ACCOUNT_SID"), ShouldNotBeEmpty)
	})
}

func TestHandlers(t *testing.T) {
	Convey("Index should return with status OK", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		So(res.Code, ShouldEqual, http.StatusOK)
	})
}
