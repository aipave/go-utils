package gerr

import "fmt"

// custom err type
type Err struct {
    error
    code int32
    msg  string
}

type Warn struct {
    mode string
}

func NewErrInst(code int32, msg string) error {
    return &Err{
        code: code,
        msg:  msg,
    }

}

func Msg(err error) string {
    if err == nil {
        return SuccessMsg
    }
    if e, ok := err.(*Err); ok {
        return e.Msg()
    }
    return err.Error()
}

// append err msg
func Append(err, appendErr error) error {
    if err == nil {
        return appendErr
    }
    if appendErr == nil {
        return err
    }
    if e, ok := err.(*Err); ok {
        return e.Append(appendErr)
    }
    return NewErrInst(DefaultErrCode, err.Error()+"; "+appendErr.Error())
}

func (e *Err) Error() string {
    return fmt.Sprintf("code: %v, msg: %v", e.Code(), e.Msg())
}

func (e *Err) Code() int32 {
    return e.code
}

func (e *Err) Msg() string {
    return e.msg
}

func (e *Err) Append(err error) error {
    if err == nil {
        return e
    }
    e.msg += "; " + err.Error()
    return e
}
