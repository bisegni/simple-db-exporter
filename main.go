package main

import (
	"github.com/bisegni/simple-db-exporter/cmd"
)

// Version ...
var Version string

// Buildtime ...
var Buildtime string

func main() {
	cmd.Version = Version
	cmd.Buildtime = Buildtime
	cmd.Execute()
}
