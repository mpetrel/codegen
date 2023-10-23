package goparse

import (
	"fmt"
	"github.com/mpetrel/codegen/internal/pkg/str"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"os"
)

type StructInfo struct {
	Name    string
	Pointer string
	Fields  []FieldsItem
}

type FieldsItem struct {
	Name    string
	Type    string
	DBTag   string
	JsonTag string
}

func ASTParse(filePath string) (*StructInfo, error) {
	// 读取文件为字符串
	codeFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	codeSrc := ""
	codeByte, err := io.ReadAll(codeFile)
	if err != nil {
		return nil, err
	}
	codeSrc = string(codeByte)
	fSet := token.NewFileSet()
	file, err := parser.ParseFile(fSet, "demo", codeSrc, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	structInfo := &StructInfo{}

	for _, node := range file.Decls {
		switch node.(type) {
		case *ast.GenDecl:
			genDecl := node.(*ast.GenDecl)
			for _, spec := range genDecl.Specs {
				switch spec.(type) {
				case *ast.TypeSpec:
					typeSpec := spec.(*ast.TypeSpec)
					fmt.Printf("Struct name: %s\n", typeSpec.Name.Name)
					structInfo.Name = typeSpec.Name.Name
					structInfo.Pointer = fmt.Sprintf("*%s", structInfo.Name)
				}
			}
		}
	}

	ast.Inspect(file, func(node ast.Node) bool {
		s, ok := node.(*ast.StructType)
		if !ok {
			return true
		}

		fields := make([]FieldsItem, 0)
		for _, field := range s.Fields.List {
			if field.Names[0].Name == "Id" {
				continue
			}
			fmt.Printf("Field: %s, Type: %s\n", field.Names[0].Name, types.ExprString(field.Type))
			fields = append(fields, FieldsItem{
				Name:    field.Names[0].Name,
				Type:    types.ExprString(field.Type),
				JsonTag: str.LowerFirst(field.Names[0].Name),
				DBTag:   getDBTag(types.ExprString(field.Type)),
			})
		}
		structInfo.Fields = fields
		return false
	})
	return structInfo, nil
}

func getDBTag(typeName string) string {
	switch typeName {
	case "int", "int32":
		return "type:int"
	case "uint64", "int64":
		return "type:bigint"
	case "int8":
		return "type:tinyint"
	case "bool":
		return "type:tinyint"
	case "time.Time":
		return "type:datetime"
	default:
		return "type:varchar(256)"
	}
}
