package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/common"
	"github.com/mpetrel/codegen/internal/pkg/str"
	"golang.org/x/exp/slices"
)

func Data(stInfo *goparse.StructInfo) *jen.File {
	f := jen.NewFile("data")

	//导入包设置
	importOrm := fmt.Sprintf("%s/internal/pkg/orm", common.ProjectName)
	importBiz := fmt.Sprintf("%s/internal/biz", common.ProjectName)
	importLogrus := "github.com/sirupsen/logrus"

	f.ImportNames(map[string]string{
		importOrm:    "orm",
		importBiz:    "biz",
		importLogrus: "logrus",
	})

	// ctx := jen.Id("ctx").Id("context.Context")
	ctx := jen.Id("ctx").Qual("context", "Context")

	excludeFields := []string{"createdAt", "updatedAt", "deletedAt"}
	// 生成 data struct
	var structFields []jen.Code
	structFields = append(structFields, jen.Qual(importOrm, "Model"))
	for _, field := range stInfo.Fields {
		if slices.Contains(excludeFields, field.Name) {
			continue
		}
		structFields = append(structFields, jen.Id(field.Name).Id(field.Type).Tag(map[string]string{"gorm": field.DBTag}))
	}
	f.Type().Id(stInfo.Name).Struct(
		structFields...,
	)

	// 公共名称
	toDataName := fmt.Sprintf("toData%s", stInfo.Name)
	toBizName := fmt.Sprintf("toBiz%s", stInfo.Name)

	// 生成repo结构体
	repoName := fmt.Sprintf("%sRepo", str.LowerFirst(stInfo.Name))
	f.Type().Id(repoName).Struct(
		jen.Id("data").Id("*Data"),
		jen.Id("log").Op("*").Qual(importLogrus, "Entry"),
	)

	// 实现repo方法
	repoMV := str.GetFirstLower(stInfo.Name)
	repoPtr := fmt.Sprintf("*%sRepo", str.LowerFirst(stInfo.Name))
	// db 调用前缀
	// dbCallPrefix := jen.Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx"))
	f.Func().Params(jen.Id(repoMV).Id(repoPtr)).
		Id("Create").
		Params(
			ctx,
			jen.Id("item").Op("*").Qual(importBiz, stInfo.Name),
		).Error().Block(
		jen.Id("dataItem").Op(":=").Id(toDataName).Call(jen.Id("item")),
		jen.If(
			jen.Err().Op(":=").Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx")).
				Dot("Create").Call(jen.Id("dataItem")).Dot("Error"),
			jen.Err().Op("!=").Nil(),
		).Block(jen.Return(jen.Err())),
		// set id
		jen.Id("item").Dot("Id").Op("=").Id("dataItem").Dot("Id"),
		jen.Return(jen.Nil()),
	)

	// 删除
	f.Func().Params(jen.Id(repoMV).Id(repoPtr)).
		Id("Delete").
		Params(ctx, jen.Id("id").Uint64()).Error().Block(
		jen.Return(
			jen.Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx")).Dot("Delete").Call(
				jen.Op("&").Id(stInfo.Name).Values(),
				jen.Id("id"),
			).Dot("Error"),
		),
	)

	// 更新
	f.Func().Params(jen.Id(repoMV).Id(repoPtr)).Id("Update").
		Params(ctx, jen.Id("item").Op("*").Qual(importBiz, stInfo.Name)).Error().Block(
		jen.Id("dataItem").Op(":=").Id(toDataName).Call(jen.Id("item")),
		jen.Return(
			jen.Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx")).Dot("Scopes").Call(jen.Id("CommonOmit").Call()).Dot("Updates").
				Call(jen.Id("dataItem")).Dot("Error"),
		),
	)

	// id查询
	f.Func().Params(jen.Id(repoMV).Id(repoPtr)).Id("Get").
		Params(
			ctx,
			jen.Id("id").Uint64()).
		Id("(").List(jen.Op("*").Qual(importBiz, stInfo.Name), jen.Error()).Id(")").
		Block(
			jen.Var().Id("item").Id(stInfo.Name),
			jen.If(
				jen.Err().Op(":=").Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx")).
					Dot("First").Call(jen.Op("&").Id("item"), jen.Id("id")).Dot("Error"),
				jen.Err().Op("!=").Nil(),
			).Block(jen.Return(jen.Nil(), jen.Err())),
			jen.Return(
				jen.Id(toBizName).Call(jen.Op("&").Id("item")),
				jen.Nil(),
			),
		)

	// list 查询
	f.Func().Params(jen.Id(repoMV).Id(repoPtr)).Id("List").
		Params(
			ctx,
			jen.Id("item").Op("*").Qual(importBiz, stInfo.Name),
			jen.Id("page"),
			jen.Id("size").Int(),
		).Id("(").
		List(jen.Op("*").Qual(importBiz, "PageResult").Index(jen.Op("*").Qual(importBiz, stInfo.Name)), jen.Error()).Id(")").
		Block(
			jen.Var().Id("items").Index().Id(stInfo.Name),
			jen.Var().Id("total").Int64(),
			jen.Err().Op(":=").Id(repoMV).Dot("data").Dot("DB").Call(jen.Id("ctx")).
				Dot("Scopes").Call(jen.Id("Paginate").Call(jen.Id("page"), jen.Id("size"))).
				Dot("Find").Call(jen.Op("&").Id("items")).
				Dot("Scopes").Call(jen.Id("ClearPaginate").Call()).
				Dot("Count").Call(jen.Op("&").Id("total")).
				Dot("Error"),
			jen.If(jen.Err().Op("!=").Nil()).Block(jen.Return(jen.Nil(), jen.Err())),

			jen.Id("ret").Op(":=").Make(
				jen.Index().Op("*").Qual(importBiz, stInfo.Name),
				jen.Len(jen.Id("items")),
			),
			jen.For(jen.Id("i").Op(":=").Range().Id("items")).Block(
				jen.Id("ret[i]").Op("=").Id(toBizName).Call(jen.Op("&").Id("items").Index(jen.Id("i"))),
			),

			jen.Return(
				jen.Op("&").Qual(importBiz, "PageResult").Index(jen.Op("*").Qual(importBiz, stInfo.Name)).
					Values(jen.Dict{
						jen.Id("Items"): jen.Id("ret"),
						jen.Id("Total"): jen.Id("total"),
					}),
				jen.Nil(),
			),
		)

	// newRepo方法
	f.Func().Id(fmt.Sprintf("New%sRepo", stInfo.Name)).Params(
		jen.Id("data").Op("*").Id("Data"),
		jen.Id("logger").Op("*").Qual("github.com/sirupsen/logrus", "Logger"),
	).Qual(importBiz, fmt.Sprintf("%sRepo", stInfo.Name)).Block(jen.Return(
		jen.Op("&").Id(repoName).Values(jen.Dict{
			jen.Id("data"): jen.Id("data"),
			jen.Id("log"): jen.Id("logger").
				Dot("WithFields").
				Call(jen.Id("logrus.Fields").Values(
					jen.Dict{jen.Lit("module"): jen.Lit("data/" + str.LowerFirst(stInfo.Name))},
				)),
		}),
	))

	// 生成结构体赋值
	bizList := make(jen.Dict)
	dataList := make(jen.Dict)
	firstLower := str.GetFirstLower(stInfo.Name)

	for _, v := range stInfo.Fields {
		bizList[jen.Id(v.Name)] = jen.Id(firstLower).Dot(v.Name)
		dataList[jen.Id(v.Name)] = jen.Id(firstLower).Dot(v.Name)
	}
	// toBiz方法
	f.Func().Id(toBizName).Params(jen.Id(firstLower).Op("*").Id(stInfo.Name)).Op("*").Qual(importBiz, stInfo.Name).
		Block(
			jen.Return(
				jen.Op("&").Qual(importBiz, stInfo.Name).Values(bizList),
			),
		)

	// toData方法
	dataList[jen.Id("Model")] = jen.Qual(importOrm, "Model").Values(
		jen.Dict{jen.Id("Id"): jen.Id(firstLower).Dot("Id")},
	)
	f.Func().Id(toDataName).Params(jen.Id(firstLower).Op("*").Qual(importBiz, stInfo.Name)).
		Op("*").Id(stInfo.Name).Block(
		jen.Return(
			jen.Op("&").Id(stInfo.Name).Values(
				dataList,
			),
		),
	)

	return f
}
