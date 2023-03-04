package test_example_test

import (
    "testing"

    "github.com/alyu01/go-utils/dbms/gmysql"
    "github.com/alyu01/go-utils/ginfos"
    "github.com/sirupsen/logrus"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

// 数据库client
var (
    readClient *gorm.DB // readonly
)

func InitMysql() {
    // user:password@tcp(localhost:3306)/mydb

    var mysqlDSN string = "root:123@tcp(localhost:3306)/mydatabase"
    readClient = gmysql.NewMySQLClient(gmysql.WithDSN(mysqlDSN), gmysql.WithLogLevel(logger.Error))
}

var MySQL _MySQLMgr

type _MySQLMgr struct {
}

// Using orm methods can simplify access and manipulation of databases while improving code readability and maintainability
// If the query condition is a null value when using the orm method, if the is null or is not null condition
// is not used for processing, it will cause uncertainty or randomness in the query results.

// In gorm, if the query condition is an empty string or zero value,
// the where method will directly ignore the condition and return all the data in the table.
// If you use more than one condition in a query, ignoring the condition may cause non-determinism
// or randomness in the query results.
//
// so the right version is
// if name = "" {
//      readClient.Debug.Where("name IS NULL")....
// } else {
//      readClient.Debug.Where("name = ?", name)....
// }

type User struct {
    Id   int64  `gorm:"column:status"`
    Name string `gorm:"column:name"`
}

func (m *_MySQLMgr) FindIdFromMytable(name string) (user User, err error) {
    user.Name = name
    rows, err := readClient.Debug().Table("mydatabase.mytable").Select("id").
        Where("name = ?", name).Rows()
    if err != nil {
        logrus.Errorf("err")
    }
    for rows.Next() {
        err := readClient.ScanRows(rows, &user.Id)
        if err != nil {
            logrus.Errorf("err")
        }

    }
    return
}

// Using orm methods can simplify access and manipulation of databases while improving code readability and maintainability
// If the query condition is a null value when using the orm method, if the is null or is not null condition
// is not used for processing, it will cause uncertainty or randomness in the query results.
func (m *_MySQLMgr) FindNameFromMytable(id int64) (user User, err error) {
    rows, err := readClient.Debug().Table("mydatabase.mytable").Select("name").
        Where("id = ?", id).Rows()
    if err != nil {
        logrus.Errorf("err")
    }
    for rows.Next() {
        err := readClient.ScanRows(rows, &user)
        if err != nil {
            logrus.Errorf("err")
        }

    }
    return
}

// Using orm methods can simplify access and manipulation of databases while improving code readability and maintainability
// If the query condition is a null value when using the orm method, if the is null or is not null condition
// is not used for processing, it will cause uncertainty or randomness in the query results.
func (m *_MySQLMgr) FindNameFromMytable2(id int64) (user User) {
    readClient.Debug().Table("mydatabase.mytable").Select("name").
        Where("id = ?", id).Find(&user.Name)
    return
}

// When using native sql statements,
// if the input parameters are not handled correctly, it may cause sql injection problems
func (m *_MySQLMgr) UpdateFromMytable(name string) (err error) {
    sqlStr := "update mydatabase.mytable set id = 1234 where name = ?"
    var res = &gorm.DB{}
    res = readClient.Debug().Exec(sqlStr, name)
    if res.Error != nil {
        logrus.Errorf("fail: err")
        return res.Error
    }
    logrus.Infof("success %v", res.RowsAffected)
    return
}

func TestInitMysql(t *testing.T) {
    InitMysql()

    id, err := MySQL.FindIdFromMytable("lyu")
    t.Log(id, err)
    logrus.Errorf("ID:%v, %v", id, err)

    t.Logf("---------%v-------", ginfos.FuncName())
    name, err := MySQL.FindNameFromMytable(1234)
    t.Log(name, err)
    logrus.Errorf("name:%v, %v", name, err)

    t.Logf("---------%v-------", ginfos.FuncName())
    name2 := MySQL.FindNameFromMytable2(0)
    t.Log(name2, err)
    logrus.Errorf("name:%v, %v", name2, err)

    t.Logf("---------%v-------", ginfos.FuncName())
    err = MySQL.UpdateFromMytable("lyu")
    t.Log(err)
}
