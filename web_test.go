package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-martini/martini"
	"github.com/joshuarubin/goscribe/telapi"
	"github.com/martini-contrib/render"
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

func TestSuccess(t *testing.T) {
	m := martini.New()
	m.Use(render.Renderer())

	type Test struct {
		Handler martini.Handler
		Result  map[string]interface{}
		Status  int
	}

	tests := []Test{
		Test{
			Handler: func(r render.Render) {
				success(r)
			},
			Result: map[string]interface{}{
				"status": float64(http.StatusOK),
			},
			Status: http.StatusOK,
		},
		Test{
			Handler: func(r render.Render) {
				success(r, "just a string")
			},
			Result: map[string]interface{}{
				"status": float64(http.StatusOK),
				"data":   "just a string",
			},
			Status: http.StatusOK,
		},
		Test{
			Handler: func(r render.Render) {
				success(r, "string0", "string1")
			},
			Result: map[string]interface{}{
				"status": float64(http.StatusOK),
				"data": []interface{}{
					"string0",
					"string1",
				},
			},
			Status: http.StatusOK,
		},
		Test{
			Handler: func(r render.Render) {
				telapiError(r, fmt.Errorf("some random error"))
			},
			Result: map[string]interface{}{
				"status": float64(http.StatusInternalServerError),
				"error":  "some random error",
			},
			Status: http.StatusInternalServerError,
		},
		Test{
			Handler: func(r render.Render) {
				telapiError(r, telapi.Error{})
			},
			Result: map[string]interface{}{
				"status": float64(http.StatusInternalServerError),
				"error":  "unknown error",
			},
			Status: http.StatusInternalServerError,
		},
	}

	Convey("Server should return properly", t, func() {
		for _, test := range tests {
			m.Action(test.Handler)

			res := httptest.NewRecorder()
			m.ServeHTTP(res, (*http.Request)(nil))

			So(res.Code, ShouldEqual, test.Status)

			data, err := ioutil.ReadAll(res.Body)
			So(err, ShouldBeNil)

			var j interface{}
			err = json.Unmarshal(data, &j)
			So(err, ShouldBeNil)

			dataMap, ok := j.(map[string]interface{})
			So(ok, ShouldBeTrue)

			So(dataMap, ShouldResemble, test.Result)
		}
	})
}
