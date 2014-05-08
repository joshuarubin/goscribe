package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandlers(t *testing.T) {
	Convey("Index should return with status OK", t, func() {
		req, _ := http.NewRequest("GET", "/", nil)
		res := httptest.NewRecorder()
		m.ServeHTTP(res, req)
		So(res.Code, ShouldEqual, http.StatusOK)
	})
}
