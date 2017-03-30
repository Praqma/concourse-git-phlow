package repo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFormatURL(t *testing.T) {

	Convey("Running tests on FormatUrl", t, func() {

		actual := FormatURL("github.com", "user", "token")
		expected := "https://user:token@github.com"

		So(actual, ShouldEqual, expected)
	})
}
