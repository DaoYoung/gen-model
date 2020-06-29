package handler

import (
    "github.com/jinzhu/gorm"
    "github.com/spf13/viper"
    "fmt"
    "strings"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

// Constants for return types of golang
const (
    golangByteArray  = "[]byte"
    gureguNullInt    = "null.Int"
    sqlNullInt       = "sql.NullInt64"
    golangInt        = "int"
    golangInt64      = "int64"
    gureguNullFloat  = "null.Float"
    sqlNullFloat     = "sql.NullFloat64"
    golangFloat32    = "float32"
    golangFloat64    = "float64"
    gureguNullString = "null.String"
    sqlNullString    = "sql.NullString"
    gureguNullTime   = "null.Time"
    golangTime       = "time.Time"
)
const (
    importSql     = "database/sql"
    importNull    = "gopkg.in/guregu/null.v3"
    importTime    = "time"
    importNothing = ""
)

func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) (goType string) {
    switch mysqlType {
    case "tinyint", "int", "smallint", "mediumint":
        if nullable {
            if gureguTypes {
                goType = gureguNullInt
            } else {
                goType = sqlNullInt
            }
            break
        }
        goType = golangInt
        break
    case "bigint":
        if nullable {
            if gureguTypes {
                goType = gureguNullInt
            } else {
                goType = sqlNullInt
            }
            break
        }
        goType = golangInt64
        break
    case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
        if nullable {
            if gureguTypes {
                goType = gureguNullString
            }
            goType = sqlNullString
            break
        }
        goType = "string"
        break
    case "date", "datetime", "time", "timestamp":
        if nullable && gureguTypes {
            goType = gureguNullTime
        } else {
            goType = golangTime
        }
        break
    case "decimal", "double":
        if nullable {
            if gureguTypes {
                goType = gureguNullFloat
            } else {
                goType = sqlNullFloat
            }
            break
        }
        goType = golangFloat64
        break
    case "float":
        if nullable {
            if gureguTypes {
                goType = gureguNullFloat
            } else {
                goType = sqlNullFloat
            }
            break
        }
        goType = golangFloat32
        break
    case "binary", "blob", "longblob", "mediumblob", "varbinary":
        goType = golangByteArray
        break
    default:
        goType = ""
    }
    return
}
func getImportPackage(golangType string) string {
    im := importNothing
    if golangType == "" {
        return im
    }
    switch {
    case len(golangType) > 2 && golangType[0:3] == "sql":
        im = importSql
    case len(golangType) > 3 && golangType[0:4] == "null":
        im = importNull
    case len(golangType) > 3 && golangType[0:4] == "time":
        im = importTime
    }
    return im
}

var dbSchema *gorm.DB

func initDb() error {
    var err error
    dsn := fmt.Sprintf(
        "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
        viper.GetString("mysql.username"),
        viper.GetString("mysql.password"),
        viper.GetString("mysql.host"),
        viper.GetInt("mysql.port"),
        "information_schema",
    )
    if dbSchema, err = gorm.Open("mysql", dsn); err != nil {
        panic(err)
    }
    if viper.GetBool("debug") {
        dbSchema.LogMode(true)
    }
    return nil
}

func matchTables(dbName, tableName string) []string {
    var names []string
    columns := &([]SchemaTable{})
    pattern := strings.Replace(tableName, "*", "%", 2)
    if err := dbSchema.Where("TABLE_SCHEMA = ?", dbName).Where("TABLE_NAME like ?", pattern).Find(columns).Pluck("TABLE_NAME", &names).Error; err != nil {
        panic(err)
    }
    return names
}

func getOneTableColumns(dbName, tableName string) *[]SchemaColumn {
    columns := &([]SchemaColumn{})
    if err := dbSchema.Where("TABLE_SCHEMA = ?", dbName).Where("TABLE_NAME = ?", tableName).Find(columns).Error; err != nil {
        panic(err)
    }
    return columns
}
