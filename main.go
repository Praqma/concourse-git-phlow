package main

import (
	"github.com/go-errors/errors"
	"fmt"
)
var Crashed = errors.Errorf("oh dear")

func Crash() error {
	return errors.New(Crashed)
}
func main() {


	err := Crash()
	if err != nil {
		if errors.Is(err, Crashed) {
			fmt.Println(err.(*errors.Error).ErrorStack())
		} else {
			panic(err)
		}
	}

}
