package test_example

import (
	"fmt"
	"testing"
	"time"

	"github.com/aipave/go-utils/ginfos"
	"github.com/aipave/go-utils/gwarn"
)

func TestGwarn(t *testing.T) {
	url := "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx"

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
