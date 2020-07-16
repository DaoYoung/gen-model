package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/spf13/viper"
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
var dbGen *gorm.DB

func initSchemaDb() {
	if dbSchema == nil {
		dbSchema = connectDb("information_schema")
	}
}
func initGenDb() {
	initSchemaDb()
	dbSchema.Exec("create database IF NOT EXISTS gen_model")
	dbGen = connectDb("gen_model")
	dbGen.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8 comment 'struct mappers'").AutoMigrate(&structMapper{})
	dbGen.Model(&structMapper{}).AddIndex("idx_db_name", "db_name")
	dbGen.Model(&structMapper{}).AddIndex("idx_table_name", "table_name")
}
func connectDb(dbName string) (dbPool *gorm.DB) {
	var err error
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		dbName,
	)
	if dbPool, err = gorm.Open("mysql", dsn); err != nil {
		panic(err)
	}
	if viper.GetBool("debug") {
		dbPool.LogMode(true)
	}
	// defer dbPool.Close()
	return
}
