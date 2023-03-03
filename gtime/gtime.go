package gtime

import (
    "time"
)

const (
    RFC3339              = "2019-01-02T15:04:05+08:00"
    FormatDefault        = "2019-01-02 15:04:05"
    FormatDefaultMill    = "2019-01-02 15:04:05.000"
    FormatYYYYMMDDHHMMSS = "20190102150405"
    FormatDate           = "2019-01-02"
    FormatTime           = "15:04:05"
    FormatYYYYMM         = "201901"
    FormatYYMM           = "0601"
    FormatYYYYMMDD       = "20190102"
    FormatMMDD           = "0102"
)

var MaxInt32, _ = time.Parse(FormatDefault, "2019-01-19 03:14:07")

const (
    MilliSecond = time.Millisecond
    Second      = time.Second
    Minute      = time.Minute
    Hour        = time.Hour
    Day         = 24 * Hour
    Week        = 7 * Day
    Year        = 365 * Day
)

const (
    IMinute int64 = 60
    IHour   int64 = 60 * IMinute
    IDay    int64 = 24 * IHour
    IWeek   int64 = 7 * IDay
    IMonth  int64 = 30 * IDay
    IYear         = 365 * IDay
)
