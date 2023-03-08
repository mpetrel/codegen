package gen

import (
	"fmt"
	"github.com/mpetrel/codegen/internal/goparse"
	"testing"
)

func TestData(t *testing.T) {
	stInfo, err := goparse.ASTParse("./biz_test.go")
	if err != nil {
		t.Error(err)
	}
	f := Data(stInfo)
	fmt.Printf("%#v", f)
}
