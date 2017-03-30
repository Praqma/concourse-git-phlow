package repo

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestValidateInRequest(t *testing.T) {

	Convey("Running tests on Validate In request", t, func() {

		Convey("Validate request should return nil", func() {

			So(ValidateInRequest(), ShouldBeNil)

		})
	})

}
