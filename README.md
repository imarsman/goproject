# Sample Golang project

This project using the [taskdev](https://taskfile.dev) tool to make a copy of a
starter go project that mostly adheres to the recommended Golang project layout
outlined [here](https://github.com/golang-standards/project-layout). The taskdev
website has instructions for installing the task tool, including using Go if you
have it installed.

The sample project has a date/time package and a package that can be used to do
useful things like detect if code is running in a test and to get paths to
various parts of a project's layout.

The sample project also comes with its own build file that can be used to both
build and package a project's build. More can be done to accomplish things like
putting things like config directories and files and log directories and files,
etc. into the packaged build.