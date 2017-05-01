package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestExistsOnOrigin(t *testing.T) {

	Convey("Should return true on a located branch", t, func() {
		Convey("Empty branch should return false", func() {

			st := BranchExistsOnOrigin("wip/ready/hej")

			So(st, ShouldBeTrue)
		})

	})
}
