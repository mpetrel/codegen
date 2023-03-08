package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mpetrel/codegen/internal/goparse"
)

func Biz(stInfo *goparse.StructInfo) *jen.File {
	f := jen.NewFile("biz")

	// 生成结构体
	var structFields []jen.Code
	for _, field := range stInfo.Fields {
		structFields = append(structFields, jen.Id(field.Name).Id(field.Type))
	}
	f.Type().Id(stInfo.Name).Struct(
		structFields...,
	)

	// 定义各类型/变量名
	repoName := fmt.Sprintf("%sRepo", stInfo.Name)
	useCaseName := fmt.Sprintf("%sUseCase", stInfo.Name)

	// 生成data接口
	f.Type().Id(repoName).Interface(
		jen.Id("Create").Params(
			jen.Id("ctx").Id("context.Context"),
			jen.Id("item").Id(stInfo.Pointer),
		),
	)

	// 生成UseCase
	f.Type().Id(useCaseName).Struct(
		jen.Id("repo").Id(repoName),
		jen.Id("tx").Id("Transaction"),
		jen.Id("log").Id("*logrus.Entry"),
	)

	f.Func().Id("New"+useCaseName).Params(
		jen.Id("repo").Id(repoName),
		jen.Id("tx").Id("Transaction"),
		jen.Id("logger").Op("*").Qual("github.com/sirupsen/logrus", "Entry"),
	).Id("*" + useCaseName).Block(
		jen.Return(
			jen.Op("&").Id(useCaseName).Values(jen.Dict{
				jen.Id("repo"): jen.Id("repo"),
				jen.Id("tx"):   jen.Id("tx"),
				jen.Id("log"): jen.Id("logger").
					Dot("WithFields").
					Call(jen.Id("logrus.Fields").Values(
						jen.Dict{jen.Lit("module"): jen.Lit("biz/" + stInfo.Name)},
					)),
			}),
		),
	)

	return f
}
