package str

import (
	"fmt"
	"testing"
)

func TestLowerFirst(t *testing.T) {
	rawStr := "UserRepo"
	expectStr := "userRepo"
	newStr := LowerFirst(rawStr)
	if newStr != expectStr {
		t.Errorf("expect %s, got %s", expectStr, newStr)
	}
	fmt.Println(newStr)
}
