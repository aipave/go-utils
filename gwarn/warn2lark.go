package gwarn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	ginfos "github.com/aipave/go-utils/ginfos"
	"github.com/sirupsen/logrus"
)

type Alarm struct {
	WebHookAddr string
	Title       string
	ShowIp      bool
	NoticeAll   bool
	CardColor   string
	FontColor   string
}

/*
 * 生成卡片
 */
func (m *Alarm) generateCardMsg(content string, level string) (c msgCard) {
	c.MsgType = "interactive" // card type

	// header
	c.Card.Header.Title.Tag = "plain_text"
	c.Card.Header.Title.Content = m.Title
	if level == "info" {
		c.Card.Header.Template = "green"
	} else {
		c.Card.Header.Template = "red"
	}
	if m.CardColor != "default" {
		c.Card.Header.Template = m.CardColor
	}

	body := fmt.Sprintf("%v", content)
	if m.ShowIp {
		body = fmt.Sprintf("%v\n(IP:%v)", content, ginfos.Runtime.IP())
	}

	// body
	if m.FontColor != "dafault" {
		c.Card.Elements = append(c.Card.Elements, element{
			Tag:     "markdown",
			Content: fmt.Sprintf("<font color='%v'>%v</font>", m.FontColor, body),
		})
	} else {
		c.Card.Elements = append(c.Card.Elements, element{
			Tag:     "markdown",
			Content: body,
		})
	}

	// dividing line
	c.Card.Elements = append(c.Card.Elements, element{
		Tag: "hr",
	})

	// @all
	if m.NoticeAll {
		c.Card.Elements = append(c.Card.Elements, element{
			Tag:     "markdown",
			Content: "<at id=all></at>",
		})
	}

	return
}

/*
 * alarm to lark
 */
func (m *Alarm) alert(contents string, level string) {
	card := m.generateCardMsg(contents, level)
	content, _ := json.Marshal(card)
	logrus.Infof("%v", string(content))
	resp, err := http.Post(m.WebHookAddr, "application/json", bytes.NewReader(content))
	if err != nil {
		logrus.Errorf("alert err:%v\n", err)
		return
	}
	defer resp.Body.Close()

	var data []byte
	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("read body err:%v\n", err)
		return
	}

	var codec struct {
		Code int64  `json:"code"`
		Msg  string `json:"msg"`
	}

	err = json.Unmarshal(data, &codec)
	if err != nil || codec.Code != 0 {
		logrus.Errorf("alert failed with err:%v resp:%+v\n", err, string(data))
	} else {
		logrus.Infof("alerted success resp:%v", string(data))
	}
}

/*
 * parse content
 */
func (m *Alarm) parse(contents interface{}) (msg string, err error) {
	switch contents.(type) {
	case []byte:
		msg = string(contents.([]byte))
	case string:
		msg = contents.(string)
		break
	case []string:
		for _, content := range contents.([]string) {
			msg += content + "\n"
		}
		break
	case [][]string:
		for _, content := range contents.([][]string) {
			msg += strings.Join(content, ",") + "\n"
		}
		break
	default:
		err = fmt.Errorf("not support type")
		break
	}
	return msg, err
}

/*
 * get alarm
 * WebHookAddr: hook addr
 * Title: subject
 */
func AleterGetter(WebHookAddr string, Title string, fns ...func(*Alarm)) *Alarm {
	a := &Alarm{
		WebHookAddr: WebHookAddr,
		Title:       Title,
		ShowIp:      true,
		NoticeAll:   true,
		CardColor:   "default",
		FontColor:   "default",
	}
	for _, fn := range fns {
		fn(a)
	}
	return a
}

// SetShowIp will not valid when use go test -v xxx_test.go, why?
func SetShowIp(showIp bool) func(t *Alarm) {
	return func(t *Alarm) {
		t.ShowIp = showIp
	}
}

func SetNoticeAll(noticeAll bool) func(t *Alarm) {
	return func(t *Alarm) {
		t.NoticeAll = noticeAll
	}
}

func SetCardColor(color string) func(t *Alarm) {
	return func(t *Alarm) {
		t.CardColor = color
	}
}

func SetFontColor(color FontColor) func(t *Alarm) {
	return func(t *Alarm) {
		t.FontColor = string(color)
	}
}

func (m *Alarm) Warning(contents interface{}) error {
	msg, err := m.parse(contents)
	if err != nil {
		return err
	}
	m.alert(msg, "alert")
	return nil
}

func (m *Alarm) Notice(contents interface{}) error {
	msg, err := m.parse(contents)
	if err != nil {
		return err
	}
	m.alert(msg, "info")
	return nil
}
