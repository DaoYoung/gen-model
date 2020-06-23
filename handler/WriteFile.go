package handler

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var Yamlfile = ".gen-model"

type DealTable struct {
	TableName string
	Columns   *[]SchemaColumn
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
func jsonWrite(data []byte) {
	fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
func containString(s []string, e string) bool {
	for _, a := range s {
		if strings.EqualFold(a, e) {
			return true
		}
	}
	return false
}

type columnProcessor struct {
	AttrSegment   string
	ImportSegment string
}

func columnProcess(columns *[]SchemaColumn, hasGormAnnotation, hasJsonAnnotation, hasGureguNullPackage bool) *columnProcessor {
	columnProcessor := &(columnProcessor{})
	var importPackages []string
	for _, column := range *columns {
		var annotations []string
		fieldType, needPackage := mysqlTypeToGoType(column.DataType, column.IsNull(), hasGureguNullPackage)
		if needPackage != "" && !containString(importPackages, needPackage) {
			importPackages = append(importPackages, needPackage)
		}
		if hasGormAnnotation {
			primary := ""
			if column.ColumnKey == "PRI" {
				primary = ";primary_key"
			}
			annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", column.ColumnName, primary))
		}
		if hasJsonAnnotation {
			annotations = append(annotations, fmt.Sprintf("json:\"%s\"", camelString(column.ColumnName)))
		}
		columnProcessor.AttrSegment += fmt.Sprintf("\n    %s %s `%s`",
			camelString(column.ColumnName),
			fieldType,
			strings.Join(annotations, " "))
	}
	if len(importPackages) > 0 {
		sort.Strings(importPackages)
		columnProcessor.ImportSegment = "import (\n"
		for _, p := range importPackages {
			columnProcessor.ImportSegment += "    \"" + p + "\"\n"
		}
		columnProcessor.ImportSegment += ")"
	}
	return columnProcessor
}
func structWrite(dealTable *DealTable, genRequest *GenRequest) {
	structName := camelString(dealTable.TableName)
	absPath, err := filepath.Abs(genRequest.OutPutPath)
	if err != nil {
		log.Println("error OutPutPath: " + genRequest.OutPutPath)
		os.Exit(0)
	}
	if !isExist(absPath) {
		log.Println("OutPutPath not exist: " + absPath)
		os.Exit(0)
	}
	packageName := ""
	appPath, err := os.Getwd()
	fileName := absPath + "/" + structName + ".go"
	if absPath == appPath {
		packageName = "main"
	} else {
		_, packageName = filepath.Split(fileName)

	}
	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	fp.Truncate(0)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	str := "package " + packageName + "\n\n"

	gormAnnotation := true
	jsonAnnotation := true
	columnProcessor := columnProcess(dealTable.Columns, gormAnnotation, jsonAnnotation, true)
	str += columnProcessor.ImportSegment
	str += "\ntype " + structName + " struct {"
	str += columnProcessor.AttrSegment
	str += "\n}\n\n"
	str += "func (model *" + structName + ") TableName() string {\n    return \"" + dealTable.TableName + "\"\n}"
	strmodel := fmt.Sprintf("%s", str)
	_, err = fp.Write([]byte(strmodel))
	if err != nil {
		log.Fatal(err)
	}
}
func GenConfigYaml(genRequest *GenRequest) {
	log.Println("GenConfigYaml", genRequest)
	content :=
		`mysql:
  host: ` + genRequest.DbConfig.Host + `
  database: ` + genRequest.DbConfig.Database + `
  port: ` + strconv.Itoa(genRequest.DbConfig.Port) + `
  username: ` + genRequest.DbConfig.Username + `
  password: ` + genRequest.DbConfig.Password + `
searchTableName: ` + genRequest.SearchTableName + ` # support patten with '*'
outPutPath: ` + genRequest.OutPutPath + ` # file path
isLowerCamelCaseJson: ` + strconv.FormatBool(genRequest.IsLowerCamelCaseJson) + ` # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'
`
	fileName := genRequest.getOutPutPath() + "/" + Yamlfile + ".yaml"
	log.Println("fileName: ", fileName)
	if isExist(fileName) {
		log.Println("you must delete config file: " + fileName + ", if you want create new one")
	}
	mkdir(genRequest.getOutPutPath())
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	defer f.Close()
	if err != nil {
		log.Println(err.Error())
	} else {
		_, err = f.Write([]byte(content))
	}
}
func mkdir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
		os.Chmod(path, 0777)
	}
}
func isExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil || os.IsExist(err)
}

// Constants for return types of golang
const (
	golangByteArray  = "[]byte"
	gureguNullInt    = "null.Int"
	sqlNullInt       = "sql.NullInt64"
	golangInt        = "int"
	golangInt64      = "int64"
	gureguNullFloat  = "null.Float"
	sqlNullFloat     = "sql.NullFloat64"
	golangFloat      = "float"
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

func mysqlTypeToGoType(mysqlType string, nullable bool, gureguTypes bool) (goType, importPackage string) {
	switch mysqlType {
	case "tinyint", "int", "smallint", "mediumint":
		if nullable {
			if gureguTypes {
				return gureguNullInt, importNull
			}
			return sqlNullInt, importSql
		}
		return golangInt, importNothing
	case "bigint":
		if nullable {
			if gureguTypes {
				return gureguNullInt, importNull
			}
			return sqlNullInt, importSql
		}
		return golangInt64, importNothing
	case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext", "json":
		if nullable {
			if gureguTypes {
				return gureguNullString, importNull
			}
			return sqlNullString, importSql
		}
		return "string", importNothing
	case "date", "datetime", "time", "timestamp":
		if nullable && gureguTypes {
			return gureguNullTime, importNull
		}
		return golangTime, importTime
	case "decimal", "double":
		if nullable {
			if gureguTypes {
				return gureguNullFloat, importNull
			}
			return sqlNullFloat, importSql
		}
		return golangFloat64, importNothing
	case "float":
		if nullable {
			if gureguTypes {
				return gureguNullFloat, importNull
			}
			return sqlNullFloat, importSql
		}
		return golangFloat32, importNothing
	case "binary", "blob", "longblob", "mediumblob", "varbinary":
		return golangByteArray, importNothing
	}
	return "", importNothing
}
