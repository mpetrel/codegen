package gen

import (
	"fmt"
	"github.com/mpetrel/codegen/internal/goparse"
	"github.com/mpetrel/codegen/internal/pkg/common"
	"testing"
	"time"
)

type RegUser struct {
	Id                uint64     `json:"id"`
	OpenId            string     `json:"openId"`
	Coins             int64      `json:"coins"`
	Name              string     `json:"name"`
	Sex               int        `json:"sex"`
	Avatar            string     `json:"avatar"`
	Level             int        `json:"level"`
	Membership        int        `json:"membership"`
	MembershipStartAt *time.Time `json:"membershipStartAt"`
	MembershipEndAt   *time.Time `json:"membershipEndAt"`
	Status            int        `json:"status"`
	Platform          int        `json:"platform"`
	CreatedAt         time.Time  `json:"createdAt"`
}

func TestData(t *testing.T) {
	common.ProjectName = "codegen"
	stInfo, err := goparse.ASTParse("./data_test.go")
	if err != nil {
		t.Error(err)
	}
	f := Data(stInfo)
	fmt.Printf("%#v", f)
}
