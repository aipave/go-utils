package test_example

import (
	"fmt"
	"testing"
	"time"

	"github.com/aipave/go-utils/gerr"
	"github.com/aipave/go-utils/ginfos"
)

func TestGerr(t *testing.T) {
	var ErrCodeNotImpl int32 = 501
	ErrCodeNotImplMsg := "The server does not support the requested feature and cannot complete the request"
	err := gerr.New(ErrCodeNotImpl, ErrCodeNotImplMsg)
	t.Log(err.Error())
	t.Log(gerr.Append(err, gerr.New(ErrCodeNotImpl, ErrCodeNotImplMsg)))

}

func TestGerrCustom(t *testing.T) {
	//mode.SetMode("test")

	gerr.Init("test", "", "topic")

	testErrSys(t)
	t.Log("--------------------------------------------------------------------------------------")
	testErrInvalidReq(t)
	t.Log("--------------------------------------------------------------------------------------")
	testCaseErrSys(t)
	t.Log("--------------------------------------------------------------------------------------")
	testCaseErrGameOn(t)
	t.Log("--------------------------------------------------------------------------------------")
	testCustomizeSysErr(t)
	t.Log("--------------------------------------------------------------------------------------")
	testCustomizeNorErr(t)

	gerr.WarnErrInfo("hhhhhhhhhhhhhhhhhhhhhhh", ginfos.FuncName()) ///< ??? only receive msg when have this line
}

func testErrSys(t *testing.T) {
	err := gerr.ErrSys
	c, d := gerr.HandleError("test ErrSys1", "zh-tw", time.Now(), err)
	fmt.Printf("test ErrSys1: code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test ErrSys2", "", time.Now(), err)
	fmt.Printf("test ErrSys2: code:{%d} desc:{%s}\n", c, d)
}

func testErrInvalidReq(t *testing.T) {
	err := gerr.ErrPermissionDenied
	c, d := gerr.HandleError("test ErrInvalidReq1", "zh-tw", time.Now(), err)
	fmt.Printf("test ErrInvalidReq1: code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test ErrInvalidReq2", "", time.Now(), err)
	fmt.Printf("test ErrInvalidReq2: code:{%d} desc:{%s}\n", c, d)
}

func testCaseErrSys(t *testing.T) {
	err := gerr.CaseErr(gerr.ErrSys, "test case err base on sys err")
	c, d := gerr.HandleError("test CaseErr ErrSys1", "zh-tw", time.Now(), err)
	fmt.Printf("test CaseErr ErrSys1: code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test CaseErr ErrSys2", "", time.Now(), err)
	fmt.Printf("test CaseErr ErrSys2: code:{%d} desc:{%s}\n", c, d)
}

func testCaseErrGameOn(t *testing.T) {
	err := gerr.CaseErr(gerr.ErrVersion, "test case err base on game on err")
	c, d := gerr.HandleError("test CaseErr", "zh-tw", time.Now(), err)
	fmt.Printf("test CaseErr : code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test CaseErr ", "", time.Now(), err)
	fmt.Printf("test CaseErr : code:{%d} desc:{%s}\n", c, d)
}

func testCustomizeSysErr(t *testing.T) {
	err := gerr.CustomizeSysErr(gerr.RetCode_kInvalidReq, "CustomizeSysErr", "CustomizeSysErr customizeInfo")
	c, d := gerr.HandleError("test CustomizeSysErr", "zh-tw", time.Now(), err)
	fmt.Printf("test CustomizeSysErr: code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test CustomizeSysErr2", "", time.Now(), err)
	fmt.Printf("test CustomizeSysErr2: code:{%d} desc:{%s}\n", c, d)
}

func testCustomizeNorErr(t *testing.T) {
	err := gerr.CustomizeNorErr(gerr.ErrSys.Code(), gerr.ERR_LEVEL_SYS, "CustomizeNorErr", "", "CustomizeNorErr customizeInfo")
	c, d := gerr.HandleError("test CustomizeNorErr1", "zh-tw", time.Now(), err)
	fmt.Printf("test CustomizeNorErr1: code:{%d} desc:{%s}\n", c, d)

	c, d = gerr.HandleError("test CustomizeNorErr2", "", time.Now(), err)
	fmt.Printf("test CustomizeNorErr2: code:{%d} desc:{%s}\n", c, d)

	err = gerr.CustomizeNorErr(100, 100, "CustomizeNorErr", "", "CustomizeNorErr customizeInfo")
	c, d = gerr.HandleError("test CustomizeNorErr3", "zh-tw", time.Now(), err)
	fmt.Printf("test CustomizeNorErr3: code:{%d} desc:{%s}\n", c, d)
}
