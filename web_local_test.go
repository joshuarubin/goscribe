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
