package main

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
)

type TestStrategy struct {
	coFail     bool
	mergeFail  int
	rebaseFail bool
}

func (t *TestStrategy) Checkout(br string) (error) {
	if t.coFail {
		return errors.New("")
	}
	return nil
}
func (t *TestStrategy) MergeFF(br string) (error) {
	if t.mergeFail == 1 {
		return nil
	}
	t.mergeFail++
	return errors.New("")
}

func (t *TestStrategy) RebaseOnto(br string) (error) {
	if t.rebaseFail {
		return errors.New("")
	}
	return nil
}

func TestApplyAndRunStrategy(t *testing.T) {

	Convey("Running test on IntegrationEngine", t, func() {

		Convey("success scenario no errors thrown", func() {
			t := TestStrategy{false, 0, false}
			err := ApplyAndRunStrategy("master", "ready", &t)
			So(err, ShouldBeNil)
		})

		Convey("ff-merge error, rebase/merge should succeed", func() {
			t := TestStrategy{false, 1, false}
			err := ApplyAndRunStrategy("master", "ready", &t)
			So(err, ShouldBeNil)
		})

		Convey("ff-merge error, rebase error", func() {
			t := TestStrategy{false, 0, true}
			err := ApplyAndRunStrategy("master", "ready", &t)
			So(err, ShouldNotBeNil)
		})

		Convey("checkout fails", func() {
			t := TestStrategy{true, 0, false}
			err := ApplyAndRunStrategy("master", "ready", &t)
			So(err, ShouldNotBeNil)
		})
	})
}
