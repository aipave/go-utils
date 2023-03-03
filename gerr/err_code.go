package gerr

const (
    Mode_Prod = "pro"
    Mode_Test = "dev"
)

const (
    // error code constant
    ERR_LEVEL_SYS = -1
    ERR_LEVEL_NOR = 0

    DefaultErrCode = 100

    // example
    SuccessCode = 0
    SuccessMsg  = "success"

    ErrCodeNotFound    = 404
    ErrCodeNotFoundMsg = "Not found"

    // example
    RetCode_kSuccess          = 0
    RetCode_kPermissionDenied = 2000
    RetCode_kSysInternal      = 2001
    RetCode_kInvalidReq       = 2002
    RetCode_kLowVersionLimit  = 2003

    RetCode_kSysInternalMsg = "System error, please try again later"
)
