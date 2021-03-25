package main

import (
	"fmt"

	"gitlab.xml.team/xmlt/goproject/cmd/goproject/internal/common"

	"github.com/jwalton/gchalk"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("App name", gchalk.Green(common.AppName()))
}
