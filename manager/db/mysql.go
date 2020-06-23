package db

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "github.com/spf13/viper"
    "log"
)

var Db *gorm.DB

func InitDb() error {
    var err error
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
        viper.GetString("mysql.username"),
        viper.GetString("mysql.password"),
        viper.GetString("mysql.host"),
        viper.GetInt("mysql.port"),
        "information_schema",
    )
    if Db, err = gorm.Open("mysql", dsn); err != nil {
        log.Println("dns",dsn)
        panic(err)
    }
    Db.LogMode(true)
    return nil
}