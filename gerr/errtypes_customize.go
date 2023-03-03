package gerr

import (
    "fmt"
    "runtime/debug"
)

var (
    Success             *ExpandError
    ErrSys              *ExpandError
    ErrPermissionDenied *ExpandError
    ErrInvalidReq       *ExpandError
    ErrVersion          *ExpandError
)

func init() {
    Success = &ExpandError{
        code: int32(RetCode_kSuccess),
    }
    ErrVersion = &ExpandError{
        code: int32(RetCode_kLowVersionLimit),
    }
    ErrSys = &ExpandError{
        code: int32(RetCode_kSysInternal),
        msg:  "System error, please try again later",
    }
    ErrInvalidReq = &ExpandError{
        code: int32(RetCode_kInvalidReq),
    }
    ErrPermissionDenied = &ExpandError{
        code: int32(RetCode_kPermissionDenied),
    }
}

func Init(mode, warnUrl, topic string) {
    initErrHandler(mode, warnUrl, topic)
}

// ExpandError
type ExpandError struct {
    level      int32
    code       int32
    msg        string
    sceneInfos string
    stackInfos string
}

func (e *ExpandError) Error() string {
    return e.msg
}

func (e *ExpandError) Code() int32 {
    return e.code
}

func (e *ExpandError) Level() int32 {
    return e.level
}

// LangError
func (e *ExpandError) LangError(lang string) string {
    //todo return with lang

    return e.Error()
}

// MoreError Error details description, mainly used to display error details when warning
func (e *ExpandError) MoreError() string {
    if e.sceneInfos != "" {
        if e.stackInfos != "" {
            return fmt.Sprintf("msg:%s \nsceneInfos:%s \nstack:%s", e.Error(), e.sceneInfos, e.stackInfos)
        }
        return fmt.Sprintf("msg:%s \nsceneInfos:%s", e.Error(), e.sceneInfos)
    } else {
        if e.stackInfos != "" {
            return fmt.Sprintf("msg:%s \nstack:%s", e.Error(), e.stackInfos)
        }
    }
    return e.Error()
}

// CaseErr Generate an error object with some custom information for specific scenarios based on underlying errors
func CaseErr(e *ExpandError, diyInfo string) *ExpandError {
    e_ := &ExpandError{
        code:       e.code,
        msg:        e.Error(),
        sceneInfos: diyInfo,
    }
    // If it is a system error that generates stack information, other errors
    // are normal business errors and do not generate stack information
    if e.code == ErrSys.code || e.level == ERR_LEVEL_SYS {
        e_.stackInfos = string(debug.Stack())
    }
    return e_
}

// CustomizeSysErr Customize the system-level error object, mainly to undertake other service errors, and you can add description information
func CustomizeSysErr(code int32, msg, diyInfo string) *ExpandError {
    return &ExpandError{
        code:       ErrSys.code,
        msg:        fmt.Sprintf("%v[%d]", msg, code),
        sceneInfos: diyInfo,
    }
}

// CustomizeNorErr
// If you don't want the error object to be initialized at the beginning, you can use this function to generate it at any time
func CustomizeNorErr(code, level int32, msg, langKey, diyInfo string) *ExpandError {
    return &ExpandError{
        level:      level,
        code:       code,
        msg:        msg,
        sceneInfos: diyInfo,
    }
}
