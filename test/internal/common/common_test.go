package common

import (
	"testing"
)

func TestApplicationRelativePath(t *testing.T) {
	t.Logf("Application relative path %s", AppBaseCodePath())
}

func TestProjectDir(t *testing.T) {
	t.Logf("Project dir %s", ProjectDir())
}

func TestRootDir(t *testing.T) {
	t.Logf("Root dir %s", RootDir())
}

func TestTestingBaseDir(t *testing.T) {
	t.Logf("Test base dir %s", TestBaseDir())
}

func TestExeParentDir(t *testing.T) {
	t.Logf("EXE parent dir %s", ExeParentDir())
}

func TestRunningTest(t *testing.T) {
	t.Logf("Running test %v", RunningTest())
}

func TestAppname(t *testing.T) {
	t.Logf("App name %v", AppName())
}
