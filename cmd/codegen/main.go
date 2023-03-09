package main

import (
	"flag"
	"github.com/mpetrel/codegen/internal/pkg/common"
)

var (
	structFile  string
	projectName string
)

func init() {
	flag.StringVar(&structFile, "sf", "", "struct file path, eg: -sf user.go")
	flag.StringVar(&projectName, "pn", "", "project name, eg: -pn shop-service")
}

func main() {
	flag.Parse()
	common.ProjectName = projectName
}
