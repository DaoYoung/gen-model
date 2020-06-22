package handler

import (
    "github.com/jinzhu/gorm"
    "github.com/spf13/viper"
    "log"
    "fmt"
)
var dbPool *gorm.DB
func connect(databaseName string) error {
    var err error
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
        viper.GetString("mysql.user"),
        viper.GetString("mysql.password"),
        viper.GetString("mysql.host"),
        viper.GetInt("mysql.port"),
        databaseName,
    )
    if dbPool, err = gorm.Open("mysql", dsn); err != nil {
        log.Println("dns",dsn)
        panic(err)
    }
    dbPool.LogMode(true)
    return nil
}