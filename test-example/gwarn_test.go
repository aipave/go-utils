package test_example

import (
    "fmt"
    "testing"
    "time"

    "github.com/alyu01/go-utils/ginfos"
    "github.com/alyu01/go-utils/gwarn"
)

func TestGwarn(t *testing.T) {
    url := "https://open.feishu.cn/open-apis/bot/v2/hook/2f1dc72c-8d2d-4641-bd95-31bbd6fcd2c7"

    titleN := "Notice Topic"
    // theSetShowIpOnTheRealMachineIsValid
    gwarn.AleterGetter(url, titleN, gwarn.SetShowIp(true), gwarn.SetCardColor("grey"),
        gwarn.SetFontColor(gwarn.FontColorRed)).
        Notice(fmt.Sprintf("area: ID\ntime: %v \nQ: hello world", time.Now().Unix()))

    titleW := "Warning Topic"
    gwarn.AleterGetter(url, titleW, gwarn.SetNoticeAll(false), gwarn.SetShowIp(true), gwarn.SetCardColor("blue"),
        gwarn.SetFontColor(gwarn.FontColorGrey)).
        Warning(fmt.Sprintf("area: ID\ntime: %v \nQ: hello world\nIP:%v", time.Now().Unix(), ginfos.Runtime.IP()))
}
