package main

import (
	"fmt"
	"os"
	"github.com/groenborg/pip/githandler"
)

func main() {


	_, err := githandler.Push()
	if err != nil {

	}
}
