package goparse

import (
	"fmt"
	"testing"
)

type User struct {
	Id   uint64
	Name string
	Sex  int
	Age  int
	Role []uint64
}

func TestASTParse(t *testing.T) {
	result, err := ASTParse("./goparse_test.go")
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("%#v", result)
}
