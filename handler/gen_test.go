package handler

import (
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"testing"
)

func TestGenerateStructFromSelfTable(t *testing.T) {
	student := mockTable()
	cmdRequest := mockConfig(student)
	log.Println("getSearchTableName", cmdRequest.Gen.getSearchTableName())
	log.Println("getSearchStructName", cmdRequest.Gen.getSearchStructName())
	log.Println("getTables", cmdRequest.getTables())
	log.Println("getOutPutPath", cmdRequest.getOutPutPath())
	structName := camelString(student.TableName)
	modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
	columnProcessor := getProcessorSelfTable(student, cmdRequest)
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}

func TestGenerateStructFromYaml(t *testing.T) {
	student := mockTable()
	cmdRequest := mockConfig(student)
	cmdRequest.Gen.ModelSuffix = "VO"
	cmdRequest.Gen.PersistType = ""
	cmdRequest.Gen.SourceType = sourceLocal
	log.Printf("%+v", cmdRequest)
	viper.Set("mock_map", mockFieldMap())
	mp, _ := filepath.Abs("../model")
	cmdRequest.Wg.Add(1)
	mkStructFromYaml(cmdRequest, "StudentVO", "model", mp)
}

func TestGenerateStructFromGenTable(t *testing.T) {
	student := mockTable()
	cmdRequest := mockConfig(student)
	cmdRequest.Gen.ModelSuffix = "BO"
	cmdRequest.Gen.PersistType = ""
	structName := camelString(student.TableName + cmdRequest.Gen.ModelSuffix)
	modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
	columnProcessor := getProcessorGenTable(student.TableName, mockGenMapper(), cmdRequest)
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}

func mockGenMapper() *[]structMapper {
	mapSlice := []structMapper{}
	sm := structMapper{
		DbName:         "test",
		TableName:      "student",
		StructName:     "StudentBO",
		ModelFieldName: "Id",
		ModelFieldType: "int",
	}
	mapSlice = append(mapSlice, sm)
	sm.ModelFieldName = "RealName"
	sm.ModelFieldType = "string"
	mapSlice = append(mapSlice, sm)
	sm.ModelFieldName = "Job"
	mapSlice = append(mapSlice, sm)
	return &mapSlice
}

func mockFieldMap() *fieldMap {
	fields := []fieldNameAndType{}
	nameAndType := fieldNameAndType{}
	nameAndType["Id"] = "int"
	fields = append(fields, nameAndType)
	nameAndType = fieldNameAndType{}
	nameAndType["RealName"] = "string"
	fields = append(fields, nameAndType)
	return &fieldMap{
		TableName: "student",
		Fields:    fields,
	}
}
func mockTable() *dealTable {
	dealTable := &(dealTable{})
	dealTable.TableName = "student"
	studentColumns := []SchemaColumn{
		{ColumnName: "id", ColumnKey: "PRI", DataType: "int", IsNullable: "NO"},
		{ColumnName: "real_name", DataType: "varchar", IsNullable: "NO"},
		{ColumnName: "job", DataType: "varchar", IsNullable: "NO", ColumnComment: "job description"},
		{ColumnName: "age", DataType: "mediumint", IsNullable: "NO"},
		{ColumnName: "sex", DataType: "tinyint", IsNullable: "NO", ColumnComment: "1:boy 2:girl"},
		{ColumnName: "birthday", DataType: "date", IsNullable: "YES"},
	}
	dealTable.Columns = &studentColumns
	return dealTable
}

func mockConfig(dt *dealTable) *CmdRequest {
	mp, _ := filepath.Abs("../model")
	mkdir(mp)
	genConf := genConfig{SearchTableName: dt.TableName, OutPutPath: "../model", IsLowerCamelCaseJson: true, HasGormTag: true, HasJsonTag: true, HasGureguNullPackage: true, SourceType: sourceSelfTable, PersistType: sourceLocal}
	return &CmdRequest{Gen: genConf}
}
