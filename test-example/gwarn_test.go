package test_example

import (
	"fmt"
	"sync"
	"testing"
	"time"

	grace "github.com/aipave/go-utils/gexit"
	"github.com/aipave/go-utils/ginfos"
	"github.com/aipave/go-utils/gwarn"
	"github.com/sirupsen/logrus"
)

type config struct {
	logLevel     *logrus.Level
	reportCaller *bool
	formatter    logrus.Formatter
	hooks        []logrus.Hook
	maxAge       int64
	filename     string
	alertUrl     string
}

type LogOption func(*config)

// WithAlerUrl
func WithAlertUrl(url string) LogOption {
	return func(cfg *config) {
		cfg.alertUrl = url
		logrus.Infof("alertUrl is %s", url)
	}
}

var pac Panic

type Panic struct {
	alertUrl string
}

func (p *Panic) watch() {
	logrus.Infof("watch alertUrl is %s", p.alertUrl)
}

// Init 初始化log rus
func Init(opts ...LogOption) {
	var cfg config
	for _, fn := range opts {
		fn(&cfg)
	}

	logrus.Infof("alertUrl is %s", cfg.alertUrl)
	var once sync.Once
	once.Do(func() {
		var p *Panic = &pac
		p.alertUrl = cfg.alertUrl
		p.watch()
	})
}

func TestInit(t *testing.T) {
	Init(WithAlertUrl("https://open.feishu.cn/open-apis/bot/v2/hook/xxxxx-xxxx-xxxx-xxxx-xxxxxxxxxx"))
}

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
	grace.Wait()
}
