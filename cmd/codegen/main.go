package main

import (
	"flag"
	"fmt"
	"github.com/mpetrel/codegen/internal/gen"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/common"
	"github.com/mpetrel/codegen/internal/pkg/str"
	"github.com/mpetrel/codegen/internal/pkg/tool"
	"os"
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

	structFile = "C:\\Users\\Larry\\GolandProjects\\sd-platform\\sd-platform-service\\internal\\biz\\order.go"
	common.ProjectName = "sd-platform-service"

	stInfo, err := goparse.ASTParse(structFile)
	if err != nil {
		panic(err)
	}

	fileName := str.LowerFirst(stInfo.Name)

	// 检查文件夹是否存在
	if err = mkdirIfNotExist("output/biz"); err != nil {
		panic(err)
	}
	if err = mkdirIfNotExist("output/data"); err != nil {
		panic(err)
	}
	if err = mkdirIfNotExist("output/service"); err != nil {
		panic(err)
	}

	bizF := gen.Biz(stInfo)
	err = bizF.Save(fmt.Sprintf("output/biz/%s.go", fileName))
	if err != nil {
		fmt.Printf("%v", err)
	}
	dataF := gen.Data(stInfo)
	err = dataF.Save(fmt.Sprintf("output/data/%s.go", fileName))
	if err != nil {
		fmt.Println(err)
	}
	serviceF := gen.Service(stInfo)
	_ = serviceF.Save(fmt.Sprintf("output/service/%s.go", fileName))
}

func mkdirIfNotExist(dirName string) error {
	exist, err := tool.PathExist(dirName)
	if err != nil {
		return err
	}
	if exist {
		return nil
	}
	// 创建文件夹
	return os.MkdirAll(dirName, os.ModePerm)
}
