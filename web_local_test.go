// +build !wercker

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/kr/pretty"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBaseUrl(t *testing.T) {
	Convey("Base URL should be set", t, func() {
		So(os.Getenv("BASE_URL"), ShouldNotBeEmpty)
	})
}

func genPostReq(path string, params *url.Values) (*http.Request, error) {
	req, err := http.NewRequest("POST", path, strings.NewReader(params.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return req, nil
}

func TestGetTranscription(t *testing.T) {
	//Convey("getTranscription should work", t, func() {
	//})
	req, _ := genPostReq("/v1/transcribe", &url.Values{
		"callback_url": {baseURL + "/v1/transcribe/process"},
		"audio_url":    {baseURL + "/audio/testing123.mp3"},
	})
	res := httptest.NewRecorder()
	m.ServeHTTP(res, req)
	body, _ := ioutil.ReadAll(res.Body)
	var data interface{}
	json.Unmarshal(body, &data)
	pretty.Println(data)
}
