package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/str"
)

func Data(stInfo *goparse.StructInfo) *jen.File {
	f := jen.NewFile("data")

	// 生成 data struct

	// 生成repo结构体
	repoName := fmt.Sprintf("%sRepo", str.LowerFirst(stInfo.Name))
	f.Type().Id(repoName).Struct(
		jen.Id("data").Id("*Data"),
		jen.Id("log").Id("*logrus.Entry"),
	)

	return f
}
