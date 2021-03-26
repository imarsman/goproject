# Sample Golang project

This project using the [taskdev](https://taskfile.dev) tool to make a copy of a
starter go project that mostly adheres to the recommended Golang project layout
outlined [here](https://github.com/golang-standards/project-layout). The taskdev
website has instructions for installing the task tool, including using Go to do
the build and install if you have Go installed. Taskdev is a lot like Gitlab's
gitlab-ci.yml file in terms of being able to execute code, set variables, and
use environment variables. Its syntax for variable epansion is slightly
different from that of Gitlab's gitlab-ci.yml because Taskdev uses Go's build in
templating language with a few helper functions and variables added.

For example, here is usage of a built in date function to get a timestamp in a
particular format:

```BUILD_TS: '{{dateInZone "2006-01-02T15:04:05Z" (now) "UTC"}}'```

The sample project has a date/time package very useful for timestamp parsing and
formatting, ISO-8601 date format manipulation, and ISO-8601 period calculations.
There is also a package that can be used to do useful things like detect if code
is running in a test and to get paths to various parts of a project's layout.
Most of the directory oriented common functions are useful in the context of
testing. You can use them to do things like easily locate the path to a
directory with test files and in production use to get the parent directory of a
bin directory housing a build binary.

The sample project also comes with its own taskdev file that can be used to both
build and package a project's build. More can be done to accomplish things like
putting things like config directories and files and log directories and files,
etc. into the packaged build.