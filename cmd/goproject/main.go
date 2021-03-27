package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	_ "embed"

	"github.com/imarsman/goproject/cmd/goproject/internal/common"
	"github.com/imarsman/goproject/internal/compress"
	"github.com/imarsman/goproject/internal/dateutils"

	"github.com/jwalton/gchalk"
)

// The Taskfile.yml generates the two files below. Golang 1.16 is required to
// use the embed capability that became available in that release.

//go:embed .appbuildts
var appBuildTS string

//go:embed .appbuildversion
var appBuildVersion string

func main() {
	fmt.Println("hello world")
	fmt.Println("App name", gchalk.Green(common.AppName()))
	fmt.Println("Date for time.Now()", dateutils.DateForTime(time.Now()))
	fmt.Println("string hello is compressed", compress.IsGzipped([]byte("hello")))
	bytes, _ := compress.GzipBytes([]byte("hello"))
	fmt.Println("compressed string hello bytes is compressed", compress.IsGzipped(bytes))

	fmt.Println()
	fmt.Println(gchalk.WithBlue().Bold("Sample Golang project"))
	fmt.Println(gchalk.Green("This is a sample Golang project. You can customize it to suit your needs."))
	fmt.Println()
	fmt.Println(gchalk.BrightYellow("Version: \t") + gchalk.Green(strings.TrimSpace(appBuildVersion)))
	fmt.Println(gchalk.BrightYellow("Build: \t\t") + gchalk.Green(strings.TrimSpace(appBuildTS)))
	fmt.Println(gchalk.BrightYellow("Platform: \t") + gchalk.Green(runtime.GOOS))
	fmt.Println(gchalk.BrightYellow("Architecture: \t") + gchalk.Green(runtime.GOARCH))
	fmt.Println()
	os.Exit(0)

}
