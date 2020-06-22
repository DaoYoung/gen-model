package db
import (
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "fmt"
    "github.com/spf13/viper"
    "os"
    "log"
    logrus "github.com/sirupsen/logrus"

)

var Db *gorm.DB
// GormLogger struct
type GormLogger struct{}

// Print - Log Formatter
func (*GormLogger) Print(v ...interface{}) {
    switch v[0] {
    case "sql":
        logrus.WithFields(
            logrus.Fields{
                "module":        "gorm",
                "type":          "sql",
                "rows_returned": v[5],
                "src":           v[1],
                "values":        v[4],
                "duration":      v[2],
            },
        ).Info(v[3])
    case "log":
        logrus.WithFields(logrus.Fields{"module": "gorm", "type": "log"}).Print(v[2])
    }
}
func InitDb() error {
    var err error
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
        viper.GetString("mysql.user"),
        viper.GetString("mysql.password"),
        viper.GetString("mysql.host"),
        viper.GetInt("mysql.port"),
        "information_schema",
    )
    if Db, err = gorm.Open("mysql", dsn); err != nil {
        log.Println("dns",dsn)
        panic(err)
    }
    Db.SetLogger(&GormLogger{})
    // todo log every sql
    Db.LogMode(true)
    Formatter := new(logrus.TextFormatter)
    logrus.SetFormatter(Formatter)
    dir, _ := os.Getwd()
    w, err := os.OpenFile(dir+"/logs/sql.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
    if err != nil {
        panic(err)
    }
    logrus.SetOutput(w)
    return nil
}