# logs

## Description
1. Define log formatter
2. Redirect log directory
3. Set up log rotation.

## Usage
```go
func yourLogic() {
    defer grecover.Recover()
    // some logic here
}
```

## Recover
The panic captured by recover will be output to the panic log and then trigger an alert
```go
func yourLogic() {
    defer grecover.Recover()

    // some logic here
}
```

## Uncaught Panic
After the program is restarted, it will be split into multiple 
log blocks according to the following separator, 
and the last block will be obtained. 
If a panic log is matched, an alert will be triggered.

```go
progress started at: ---------%v-----------
```