# yoho-logs

## 说明：
1. 定义log formatter
1. 重定向log日志目录
1. log rotate设置

## Usage
```go
func yourLogic() {
    defer yrecover.Recover()

    // some logic here
}
```

## Recover
使用recover捕获的panic会输出到panic.log然后触发alert告警
```go
func yourLogic() {
    defer yrecover.Recover()

    // some logic here
}
```

## 未捕获的panic
程序重新被拉起来后，会根据如下分割线分割成多个日志块，并获取最后一个块
如果匹配到panic日志，会触发alert告警
```go
progress started at: ---------%v-----------
```