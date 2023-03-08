package gen

import (
	"fmt"
	"github.com/mpetrel/codegen/internal/goparse"
	"testing"
)

type User struct {
	Id   uint64
	Name string
	Sex  int
	Age  int
	Role []uint64
}

func TestBiz(t *testing.T) {
	stInfo, err := goparse.ASTParse("./gen_test.go")
	if err != nil {
		t.Error(err)
	}
	f := Biz(stInfo)
	fmt.Printf("%#v", f)
}
