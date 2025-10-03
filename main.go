package main

import (
	"github.com/attachmentgenie/golang-example/cmd"
	promVersion "github.com/prometheus/common/version"
)

var (
	commit  = "none"
	date    = "unknown"
	service = "example"
	version = "dev"
)

func main() {
	promVersion.Revision = commit
	promVersion.BuildDate = date
	promVersion.Version = version
	cmd.Service = service

	cmd.Execute()
}
