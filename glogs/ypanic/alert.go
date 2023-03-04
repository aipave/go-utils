package ypanic

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"

    "github.com/alyu01/go-utils/ginfos"
    "github.com/sirupsen/logrus"
)

// https://open.feishu.cn/document/ukTMukTMukTM/uEjNwUjLxYDM14SM2ATN

type card struct {
    MsgType string `json:"msg_type"`
    Card    struct {
        Header struct {
            Title struct {
                Tag     string `json:"tag"`
                Content string `json:"content"`
            } `json:"title"`
            Template string `json:"template"`
        } `json:"header"`
        Elements []Element `json:"elements"`
    } `json:"card"`
}

type Element struct {
    Tag     string `json:"tag"`
    Content string `json:"content"`
}

func buildAlert(stack string) (c card) {
    c.MsgType = "interactive" // card type

    // header
    c.Card.Header.Title.Tag = "plain_text"
    c.Card.Header.Title.Content = fmt.Sprintf("%v panic alert", ginfos.Runtime.Exec())
    c.Card.Header.Template = "red"

    // body
    c.Card.Elements = append(c.Card.Elements, Element{
        Tag:     "markdown",
        Content: fmt.Sprintf("IP: %v\n%v", ginfos.Runtime.IP(), stack),
    })

    // dividing line
    c.Card.Elements = append(c.Card.Elements, Element{
        Tag: "hr",
    })

    // @all
    c.Card.Elements = append(c.Card.Elements, Element{
        Tag:     "markdown",
        Content: "<at id=all></at>",
    })

    return
}

func triggerAlert(card interface{}) {
    var currentIP = ginfos.Runtime.IP()
    for _, devIP := range []string{"10.10.0.111"} {
        if currentIP == devIP {
            return
        }

        if strings.HasPrefix(currentIP, "172.") || strings.HasPrefix(currentIP, "10.10.100.") {
            return
        }
    }

    content, _ := json.Marshal(card)
    resp, err := http.Post("https://open.feishu.cn/open-apis/bot/v2/hook/8e7d3ce7-3371-44fb-910f-00254177f0d1", "application/json", bytes.NewReader(content))
    if err != nil {
        logrus.Errorf("alert err:%v\n", err)
        return
    }
    defer resp.Body.Close()

    var data []byte
    data, err = io.ReadAll(resp.Body) ///< ioutil.ReadAll deprecated
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