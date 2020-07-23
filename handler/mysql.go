package handler

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" // Register mysql
	"github.com/spf13/viper"
)

// Constants for return types of golang
const (
	gureguNullInt    = "null.Int"
	gureguNullFloat  = "null.Float"
	gureguNullString = "null.String"
	golangByteArray  = "[]byte"
	golangString     = "string"
	golangInt        = "int"
	golangInt64      = "int64"
	golangFloat64    = "float64"
	golangNullTime   = "*time.Time"
	golangTime       = "time.Time"
)

const (
	importNull    = "gopkg.in/guregu/null.v4"
	importTime    = "time"
	importNothing = ""
)

func mysqlTypeToGoType(mysqlType string, nullable bool) string {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			return gureguNullInt
		}
		return golangInt
	case "bigint":
		if nullable {
			return gureguNullInt
		}
		return golangInt64
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			return gureguNullString
		}
		return golangString
	case "date", "datetime", "time", "timestamp":
		if nullable {
			return golangNullTime
		}
		return golangTime
	case "decimal", "double", "float":
		if nullable {
			return gureguNullFloat
		}
		return golangFloat64
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray
	}
	return ""
}
func getImportPackage(golangType string) string {
	im := importNothing
	if golangType == "" {
		return im
	}
	switch {
	case len(golangType) > 3 && golangType[0:4] == "null":
		im = importNull
	case len(golangType) > 3 && (golangType[0:4] == "time" || golangType[1:5] == "time"):
		im = importTime
	}
	return im
}

var dbSchema *gorm.DB
var dbGen *gorm.DB

func initSchemaDb() {
	if dbSchema == nil {
		var err error
		if dbSchema, err = connectDb("information_schema"); err != nil {
			printErrorAndExit(err)
		}
	}
}
func initGenDb() {
	initSchemaDb()
	if dbSchema != nil {
		dbSchema.Exec("create database IF NOT EXISTS gen_model")
	}
	if dbGen == nil {
		var err error
		dbGen, err = connectDb("gen_model")
		if err == nil {
			dbGen.Set("gorm:table_options", "ENGINE=InnoDB DEFAULT CHARSET=utf8 comment 'struct mappers'").AutoMigrate(&structMapper{})
			dbGen.Model(&structMapper{}).AddIndex("idx_db_name", "db_name")
			dbGen.Model(&structMapper{}).AddIndex("idx_table_name", "table_name")
		} else {
			printErrorAndExit(err)
		}
	}
}
func connectDb(dbName string) (dbPool *gorm.DB, err error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True",
		viper.GetString("mysql.username"),
		viper.GetString("mysql.password"),
		viper.GetString("mysql.host"),
		viper.GetInt("mysql.port"),
		dbName,
	)
	if dbPool, err = gorm.Open("mysql", dsn); err != nil {
		return nil, err
	}
	if viper.GetBool("debug") {
		dbPool.LogMode(true)
	}
	return
}
