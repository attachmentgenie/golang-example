package main

import (
	"github.com/attachmentgenie/golang-example/cmd"
	promversion "github.com/prometheus/common/version"
)

var (
	commit  = "none"
	date    = "unknown"
	service = "example"
	version = "dev"
)

func main() {
	promversion.Revision = commit
	promversion.BuildDate = date
	promversion.Version = version
	cmd.Service = service

	cmd.Execute()
}
