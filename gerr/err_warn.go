package gerr

import (
    "fmt"
    "time"

    helper "github.com/alyu01/go-utils/ginfos"
    "github.com/alyu01/go-utils/gtime"
    yalert "github.com/alyu01/go-utils/gwarn"

    "github.com/sirupsen/logrus"
)

type Warn struct {
    alerter *yalert.Alarm
    mode    string
}

var warn *Warn

// initErrHandler
func initErrHandler(mode, warnUrl, topic string) {
    warn = &Warn{
        alerter: yalert.AleterGetter(warnUrl, topic),
        mode:    mode,
    }
}

// WarnErrInfo
func WarnErrInfo(info, funcN string) {
    if warn.mode != Mode_Prod {
        //	return
    }

    msg := fmt.Sprintf("Service:%s \nFunction:%s \nError:%s", warn.mode+":"+helper.Runtime.Exec(), funcN, info)
    if warn.alerter.WebHookAddr == "" {
        logrus.Warnf("[WarnErrInfo] alerter.WebHookAddr is nil, msg:%v", msg)
        return
    }

    _ = warn.alerter.Warning(msg)
}

// NoticeInfo
func NoticeInfo(info, fun string) {
    if warn.mode != Mode_Prod {
        //	return
    }

    msg := fmt.Sprintf("Service:%s \nFunction:%s \ninfo:%s", warn.mode+":"+helper.Runtime.Exec(), fun, info)
    if warn.alerter.WebHookAddr == "" {
        logrus.Warnf("[NoticeInfo] alerter.WebHookAddr is nil, msg:%v", msg)
        return
    }

    _ = warn.alerter.Notice(msg)
}

// HandleError  Handle errors and report prometheus label interface name lang multilingual start interface call start time
func HandleError(label string, lang string, start time.Time, err error) (int32, string) {
    code := Success.Code()
    desc := Success.Error()
    // isSysErr := false

    if err != nil {
        errExpand, ok := err.(*ExpandError)
        if !ok {
            errExpand = CustomizeSysErr(ErrSys.code, err.Error(), "")
        }

        desc = errExpand.LangError(lang)
        code = errExpand.Code()

        // System-level error DB error, redis error, rpc error, etc.,
        // you can specify which system-level errors are counted by yourself,
        // report the failure of the prometheus interface,
        // and you can also send an alarm to Feishu
        if code == ErrSys.Code() || errExpand.Level() == ERR_LEVEL_SYS {
            logrus.Errorf("[HandleError] label:%s error:%s", label, desc)
            // isSysErr = true
            go func() {
                WarnErrInfo(errExpand.MoreError(), label)
                // todo prometheus
            }()
        } else {
            // General errors are not processed. Parameters are wrong, not in the appropriate range,
            // not in opening hours, etc. You can specify which ones are considered general errors.
            logrus.Warnf("[HandleError] label:%s warn:%s", label, errExpand.MoreError())
        }
    }

    // todo prometheus
    //if !isSysErr {
    //    // go gprometheus.Inc
    //}

    cost := time.Since(start)
    logrus.Infof("%v: start:%v cost:%v", label, start.Format(gtime.FormatDefault), cost)
    // todo prometheus
    // go gprometheus.Inc
    return code, desc
}
