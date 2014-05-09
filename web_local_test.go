// +build !wercker

package main

import (
	"encoding/json"
	"fmt"
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

func TestGetTranscription(t *testing.T) {
	//Convey("getTranscription should work", t, func() {
	//})
	callbackURL := baseURL + "/v1/transcribe/process"
	audioURL := baseURL + "/audio/testing123.mp3"
	req, _ := http.NewRequest("POST", "/v1/transcribe", strings.NewReader(fmt.Sprintf(
		"callback_url=%s&audio_url=%s",
		url.QueryEscape(callbackURL),
		url.QueryEscape(audioURL),
	)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	res := httptest.NewRecorder()
	m.ServeHTTP(res, req)
	body, _ := ioutil.ReadAll(res.Body)
	var data interface{}
	json.Unmarshal(body, &data)
	pretty.Println(data)
}
