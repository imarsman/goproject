package main

import (
	"fmt"
	"time"

	"gitlab.xml.team/xmlt/goproject/cmd/goproject/internal/common"
	"gitlab.xml.team/xmlt/goproject/internal/dateutils"

	"github.com/jwalton/gchalk"
)

func main() {
	fmt.Println("hello world")
	fmt.Println("App name", gchalk.Green(common.AppName()))
	fmt.Println("Date for time.Now()", dateutils.DateForTime(time.Now()))
}
