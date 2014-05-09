// +build !wercker

package main

import (
	"os"
	"testing"
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
	getTranscription(baseURL + "/audio/testing123.mp3")
}
