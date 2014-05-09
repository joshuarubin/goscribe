package telapi

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestEnvironment(t *testing.T) {
	Convey("TelAPI Account SID should be set", t, func() {
		So(AccountSID, ShouldNotBeEmpty)
	})

	Convey("TelAPI Auth Token should be set", t, func() {
		So(AuthToken, ShouldNotBeEmpty)
	})
}
