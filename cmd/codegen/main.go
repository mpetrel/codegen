package main

import "flag"

var (
	structFile string
)

func init() {
	flag.StringVar(&structFile, "sf", "", "struct file path, eg: -sf user.go")
}

func main() {
	flag.Parse()

}
