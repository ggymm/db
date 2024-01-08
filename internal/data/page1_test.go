package data

import (
	"testing"

	"db/internal/data/page"
)

func Test_Page1(t *testing.T) {
	p := page.NewPage(1, InitPage1(), nil)

	// 模拟正常启动
	SetVcOpen(p)

	// 模拟正常关闭
	SetVcClose(p)

	// 校验
	res := CheckVc(p)
	if res {
		t.Log("success")
	} else {
		t.Fatal("failed")
	}
}

func Test_Page1_Error(t *testing.T) {
	p := page.NewPage(1, InitPage1(), nil)

	// 模拟正常启动
	SetVcOpen(p)

	// 异常关闭，没有执行关闭流程
	// SetVcClose(p)

	// 校验
	res := CheckVc(p)
	if !res {
		t.Log("success")
	} else {
		t.Fatal("failed")
	}
}
