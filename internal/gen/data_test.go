package gen

import (
	"fmt"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/common"
	"testing"
)

func TestData(t *testing.T) {
	common.ProjectName = "codegen"
	stInfo, err := goparse.ASTParse("./biz_test.go")
	if err != nil {
		t.Error(err)
	}
	f := Data(stInfo)
	fmt.Printf("%#v", f)
}
