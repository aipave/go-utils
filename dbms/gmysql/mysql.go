package gmysql

import (
    "database/sql"
    "time"

    "github.com/sirupsen/logrus"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// NewMySQLClient initialclient
func NewMySQLClient(opts ...Option) *gorm.DB {
    var c = defaultConfig
    for _, fn := range opts {
        fn(&c)
    }

    db, err := gorm.Open(mysql.Open(c.DSN), &gorm.Config{
        Logger: logger.New(logrus.StandardLogger(), logger.Config{
            SlowThreshold:             100 * time.Millisecond, // 100 ms
            Colorful:                  true,
            IgnoreRecordNotFoundError: false,
            LogLevel:                  c.LogLevel,
        }),
    })

    if err != nil {
        logrus.Fatalf("open mysql fail:%s, dsn:%v", err.Error(), c.DSN)
    }

    var sqlDB *sql.DB
    sqlDB, err = db.DB()
    if err != nil {
        logrus.Fatalf("get db fail:%v, dsn: %v", err.Error(), c.DSN)
    }

    sqlDB.SetMaxIdleConns(c.MaxIdleCons)
    sqlDB.SetMaxOpenConns(c.MaxOpenCons)
    sqlDB.SetConnMaxLifetime(c.ConnMaxLifeTime)

    logrus.Infof("initialized mysql: %v", c.DSN)
    return db
}

// CloseRows .
func CloseRows(rows *sql.Rows) {
    if rows != nil {
        err := rows.Close()
        if err != nil {
            logrus.Errorf("close rows:%v err:%v", rows, err)
        }
    }
}
