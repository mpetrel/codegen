package gen

import (
	"fmt"
	"github.com/dave/jennifer/jen"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/common"
	"github.com/mpetrel/codegen/internal/pkg/str"
)

func Service(stInfo *goparse.StructInfo) *jen.File {

	f := jen.NewFile("service")

	// 名称
	serviceName := fmt.Sprintf("%sService", stInfo.Name)
	useCaseName := fmt.Sprintf("%sUseCase", stInfo.Name)

	importBiz := fmt.Sprintf("%s/internal/biz", common.ProjectName)

	// 指定导入包名
	f.ImportNames(map[string]string{
		"github.com/sirupsen/logrus": "logrus",
		"github.com/gin-gonic/gin":   "gin",
		importBiz:                    "biz",
	})

	// 生成ServiceStruct
	f.Type().Id(serviceName).Struct(
		jen.Id("uc").Op("*").Qual(importBiz, useCaseName),
		jen.Id("log").Id("*logrus.Entry"),
	)

	// New方法
	f.Func().Id("New"+serviceName).Params(
		jen.Id("uc").Op("*").Qual(importBiz, useCaseName),
		jen.Id("logger").Op("*").Qual("github.com/sirupsen/logrus", "Entry"),
	).Id("*" + serviceName).Block(
		jen.Return(
			jen.Op("&").Id(serviceName).Values(
				jen.Dict{
					jen.Id("uc"): jen.Id("uc"),
					jen.Id("log"): jen.Id("logger").
						Dot("WithFields").
						Call(jen.Id("logrus.Fields").Values(
							jen.Dict{jen.Lit("module"): jen.Lit("service/" + str.LowerFirst(stInfo.Name))},
						)),
				},
			),
		),
	)

	// 同意方法参数
	ginCtx := jen.Id("c").Op("*").Qual("github.com/gin-gonic/gin", "Context")

	// get 方法
	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("Get").Params(ginCtx).Block(
		jen.List(jen.Id("id"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(
			jen.Id("c").Dot("Param").Call(jen.Lit("id")),
			jen.Lit(10),
			jen.Lit(64),
		),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(
				jen.Qual(importBiz, "ErrInvalidParam"),
			),
			jen.Return(),
		),
		jen.List(jen.Id("item"), jen.Err()).Op(":=").Id("s").Dot("uc").Dot("Get").Call(
			jen.Id("c").Dot("Request").Dot("Context").Call(),
			jen.Id("id"),
		),
		jen.Id("jsonResponse").Call(
			jen.Id("c"),
			jen.Id("item"),
			jen.Err(),
		),
	)

	// Create
	createVarName := str.LowerFirst(stInfo.Name)
	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("Create").Params(ginCtx).Block(
		jen.Var().Id(createVarName).Qual(importBiz, stInfo.Name),
		jen.If(
			jen.Err().Op(":=").Id("c").Dot("BindJSON").Call(jen.Id("&"+createVarName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(jen.Err()),
			jen.Return(),
		),
		jen.Id("jsonResponse").Call(
			jen.Id("c"),
			jen.Id("OK"),
			jen.Id("s").Dot("uc").Dot("Create").Call(
				jen.Id("c").Dot("Request").Dot("Context").Call(),
				jen.Id("&"+createVarName),
			),
		),
	)

	// Update
	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("Update").Params(ginCtx).Block(
		jen.Var().Id(createVarName).Qual(importBiz, stInfo.Name),
		jen.If(
			jen.Err().Op(":=").Id("c").Dot("BindJSON").Call(jen.Id("&"+createVarName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(jen.Err()),
			jen.Return(),
		),
		jen.Id("jsonResponse").Call(
			jen.Id("c"),
			jen.Id("OK"),
			jen.Id("s").Dot("uc").Dot("Update").Call(
				jen.Id("c").Dot("Request").Dot("Context").Call(),
				jen.Id("&"+createVarName),
			),
		),
	)

	// List
	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("List").Params(ginCtx).Block(
		jen.List(
			jen.Id("page"),
			jen.Id("size"),
			jen.Err(),
		).Op(":=").Id("GetPageSize").Call(jen.Id("c")),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(jen.Err()),
			jen.Return(),
		),
		// 绑定查询参数到struct
		jen.Var().Id(createVarName).Qual(importBiz, stInfo.Name),
		jen.If(
			jen.Err().Op(":=").Id("c").Dot("ShouldBindUri").Call(jen.Op("&").Id(createVarName)),
			jen.Err().Op("!=").Nil(),
		).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(jen.Err()),
			jen.Return(),
		),

		jen.List(jen.Id("ret"), jen.Err()).Op(":=").Id("s").Dot("uc").Dot("List").Call(
			jen.Id("c").Dot("Request").Dot("Context").Call(),
			jen.Op("&").Id(createVarName),
			jen.Id("page"),
			jen.Id("size"),
		),
		jen.Id("jsonResponse").Call(
			jen.Id("c"),
			jen.Id("ret"),
			jen.Err(),
		),
	)

	// Delete

	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("Delete").Params(ginCtx).Block(
		jen.List(jen.Id("id"), jen.Err()).Op(":=").Qual("strconv", "ParseUint").Call(
			jen.Id("c").Dot("Param").Call(jen.Lit("id")),
			jen.Lit(10),
			jen.Lit(64),
		),
		jen.If(jen.Err().Op("!=").Nil()).Block(
			jen.Id("_").Op("=").Id("c").Dot("Error").Call(
				jen.Err(),
			),
			jen.Return(),
		),
		jen.Id("jsonResponse").Call(
			jen.Id("c"),
			jen.Id("OK"),
			jen.Id("s").Dot("uc").Dot("Delete").Call(
				jen.Id("c").Dot("Request").Dot("Context").Call(),
				jen.Id("id"),
			),
		),
	)

	// add route
	f.Func().Params(jen.Id("s").Id("*"+serviceName)).
		Id("AddRoute").Params(jen.Id("rg").Qual("github.com/gin-gonic/gin", "RouterGroup")).Block(
		jen.Id(createVarName).Op(":=").Id("rg").Dot("Group").Call(
			jen.Lit("/"+createVarName),
		),
		jen.Id(createVarName).Dot("GET").Call(jen.Lit("/:id"), jen.Id("s").Dot("Get")),
		jen.Id(createVarName).Dot("GET").Call(jen.Lit(""), jen.Id("s").Dot("List")),
		jen.Id(createVarName).Dot("POST").Call(jen.Lit(""), jen.Id("s").Dot("Create")),
		jen.Id(createVarName).Dot("PUT").Call(jen.Lit(""), jen.Id("s").Dot("Update")),
		jen.Id(createVarName).Dot("DELETE").Call(jen.Lit("/:id"), jen.Id("s").Dot("Delete")),
	)

	return f
}
