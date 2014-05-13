package telapi

import (
	"net/http"
	"net/http/httptest"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestError(t *testing.T) {
	const jsonErrorTxt0 = `{
    "status": 403,
    "message": "Invalid credentials supplied",
    "code": 10004,
    "more_info": "http://www.telapi.com/docs/api/rest/overview/errors/10004"
}`

	var jsonErr0 *JSONError
	var res0 *http.Response

	Convey("Standard error should work", t, func() {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}

		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		res, err := http.Get(ts.URL)
		So(err, ShouldBeNil)

		body, err := responseError(res)
		So(body, ShouldBeNil)
		So(err, ShouldNotBeNil)

		tmpErr, ok := err.(Error)
		So(ok, ShouldBeTrue)

		_, ok = tmpErr.Data.(*JSONError)
		So(ok, ShouldBeFalse)

		So(tmpErr.Data, ShouldBeNil)
		So(tmpErr.Status(), ShouldEqual, http.StatusInternalServerError)
		So(tmpErr.Body, ShouldBeEmpty)
	})

	Convey("JSONError should work", t, func() {
		handler := func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusForbidden)
			w.Header().Set("Content-Type", "application/json")
			_, err := w.Write([]byte(jsonErrorTxt0))
			So(err, ShouldBeNil)
		}

		ts := httptest.NewServer(http.HandlerFunc(handler))
		defer ts.Close()

		var err error
		res0, err = http.Get(ts.URL)
		So(err, ShouldBeNil)

		body, err := responseError(res0)
		So(body, ShouldBeNil)
		So(err, ShouldNotBeNil)

		tmpErr, ok := err.(Error)
		So(ok, ShouldBeTrue)

		jsonErr0, ok = tmpErr.Data.(*JSONError)
		So(ok, ShouldBeTrue)

		So(jsonErr0.Status, ShouldEqual, http.StatusForbidden)
		So(jsonErr0.Message, ShouldEqual, "Invalid credentials supplied")
		So(jsonErr0.Code, ShouldEqual, 10004)
		So(jsonErr0.MoreInfo, ShouldEqual, "http://www.telapi.com/docs/api/rest/overview/errors/10004")
	})

	Convey("JSONError should turn into telapi.Error", t, func() {
		So(Error{}.Status(), ShouldEqual, 0)
		So(jsonErr0, ShouldNotBeNil)

		taError := Error{
			Data: jsonErr0,
		}

		So(taError.Status(), ShouldEqual, http.StatusForbidden)

		tmpStatus, tmpJSON := taError.JSON()
		So(tmpStatus, ShouldEqual, http.StatusForbidden)
		So(tmpJSON, ShouldResemble, map[string]interface{}{
			"status": http.StatusForbidden,
			"data":   jsonErr0,
		})
	})
}
