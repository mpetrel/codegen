package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/str"
)

func Biz(stInfo *goparse.StructInfo) *jen.File {
	f := jen.NewFile("biz")

	// 指定导入包名
	f.ImportNames(map[string]string{
		"github.com/sirupsen/logrus": "logrus",
	})
	// 生成结构体
	var structFields []jen.Code
	for _, field := range stInfo.Fields {
		structFields = append(structFields, jen.Id(field.Name).Id(field.Type).Tag(map[string]string{"json": field.JsonTag}))
	}
	f.Type().Id(stInfo.Name).Struct(
		structFields...,
	)

	// 定义各类型/变量名
	repoName := fmt.Sprintf("%sRepo", stInfo.Name)
	useCaseName := fmt.Sprintf("%sUseCase", stInfo.Name)

	ctx := jen.Id("ctx").Qual("context", "Context")
	// 生成data接口
	f.Type().Id(repoName).Interface(
		jen.Id("Create").Params(
			ctx,
			jen.Id("item").Id(stInfo.Pointer),
		).Error(),
		jen.Id("Delete").Params(
			ctx,
			jen.Id("id").Uint64(),
		).Error(),
		jen.Id("Update").Params(
			ctx,
			jen.Id("item").Id(stInfo.Pointer),
		).Error(),
		jen.Id("Get").Params(
			ctx,
			jen.Id("id").Uint64(),
		).Id("(").List(jen.Id(stInfo.Pointer), jen.Error()).Id(")"),
		jen.Id("List").Params(
			ctx,
			jen.Id("item").Id(stInfo.Pointer),
			jen.Id("page"),
			jen.Id("size").Int(),
		).Id("(").List(jen.Id("*PageResult").Index(jen.Id(stInfo.Pointer)), jen.Error()).Id(")"),
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
						jen.Dict{jen.Lit("module"): jen.Lit("biz/" + str.LowerFirst(stInfo.Name))},
					)),
			}),
		),
	)

	// 生成 use-case 方法
	f.Func().Params(
		jen.Id("uc").Id("*"+useCaseName),
	).Id("Create").Params(
		ctx,
		jen.Id("item").Id(stInfo.Pointer),
	).Error().Block(
		jen.Return(
			jen.Id("uc.repo.Create(ctx, item)"),
		),
	)

	f.Func().Params(
		jen.Id("uc").Id("*"+useCaseName),
	).Id("Delete").Params(
		ctx,
		jen.Id("id").Uint64(),
	).Error().Block(
		jen.Return(
			jen.Id("uc.repo.Delete(ctx, id)"),
		),
	)

	f.Func().Params(
		jen.Id("uc").Id("*"+useCaseName),
	).Id("Update").Params(
		ctx,
		jen.Id("item").Id(stInfo.Pointer),
	).Error().Block(
		jen.Return(
			jen.Id("uc.repo.Update(ctx, item)"),
		),
	)

	f.Func().Params(
		jen.Id("uc").Id("*"+useCaseName),
	).Id("Get").Params(
		ctx,
		jen.Id("id").Uint64(),
	).Id("(").List(jen.Id(stInfo.Pointer), jen.Error()).Id(")").
		Block(
			jen.Return(
				jen.Id("uc.repo.Get(ctx, id)"),
			),
		)

	f.Func().Params(
		jen.Id("uc").Id("*"+useCaseName),
	).Id("List").Params(
		ctx,
		jen.Id("item").Id(stInfo.Pointer),
		jen.Id("page"),
		jen.Id("size").Int(),
	).Id("(").List(jen.Id("*PageResult").Index(jen.Id(stInfo.Pointer)), jen.Error()).Id(")").
		Block(
			jen.Return(
				jen.Id("uc.repo.List(ctx, item, page, size)"),
			),
		)

	return f
}
