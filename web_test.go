package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestHandlers(t *testing.T) {
	Convey("Index should return with status OK", t, func() {
		request, _ := http.NewRequest("GET", "/", nil)
		response := httptest.NewRecorder()

		handleIndex(response, request)
		So(response.Code, ShouldEqual, http.StatusOK)
	})
}
