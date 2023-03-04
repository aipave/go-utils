package gmysql

import (
    "time"

    "gorm.io/gorm/logger"
)

type SQLConfig struct {
    DSN             string          `json:"DSN"             yaml:"DSN"`
    LogLevel        logger.LogLevel `json:"LogLevel"        yaml:"LogLevel"`
    MaxIdleCons     int             `json:"MaxIdleCons"     yaml:"MaxIdleCons"`
    MaxOpenCons     int             `json:"MaxOpenCons"     yaml:"MaxOpenCons"`
    ConnMaxLifeTime time.Duration   `json:"ConnMaxLifeTime" yaml:"ConnMaxLifeTime"`
}

var defaultConfig = SQLConfig{
    DSN:             "",
    LogLevel:        logger.Info,
    MaxIdleCons:     10,
    MaxOpenCons:     20,
    ConnMaxLifeTime: 30 * time.Minute, // 主动超时时间，设置30分钟减少TIME_WAIT数量
}

type Option func(c *SQLConfig)

// WithLogLevel
func WithLogLevel(level logger.LogLevel) Option {
    return func(c *SQLConfig) {
        c.LogLevel = level
    }
}

// WithDSN: data source name, dsn) is an identifier used to identify the connection information of the database system,
// such as database type, server address, port number, user name and password, etc.
/* When using this function, we can set up the database connection as follows:
config := NewSQLConfig(
    WithDSN("user:password@tcp(localhost:3306)/mydb"),
    WithMaxOpenConns(10),
    WithMaxIdleConns(5),
)
*/
func WithDSN(dsn string) Option {
    return func(c *SQLConfig) {
        c.DSN = dsn
    }
}

// WithMaxIdleCons idle connections count
func WithMaxIdleCons(n int) Option {
    return func(c *SQLConfig) {
        c.MaxIdleCons = n
    }
}

// WithMaxOpenCons open connections count
func WithMaxOpenCons(n int) Option {
    return func(c *SQLConfig) {
        c.MaxOpenCons = n
    }
}

// WithConnMaxLifeTime max expiration time of a connection
func WithConnMaxLifeTime(n time.Duration) Option {
    return func(c *SQLConfig) {
        c.ConnMaxLifeTime = n
    }
}
