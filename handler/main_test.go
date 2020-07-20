package handler

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"log"
	"path/filepath"
	"testing"
)

var modelDir = "../model"

type gromMockTest struct {
	suite.Suite
	Db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *gromMockTest) SetupSuite() {
	db, mock, err := sqlmock.New()
	s.mock = mock
	require.NoError(s.T(), err)
	s.Db, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)
	s.Db.LogMode(true)
	viper.Set("is_test", true)
	viper.Set("forceCover", true)
}
func (s *gromMockTest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func TestInit(t *testing.T) {
	suite.Run(t, new(gromMockTest))
}

func (s *gromMockTest) TestMatchTables() {
	var (
		dbName    = "edu"
		tableName = "stu%"
	)
	rows := sqlmock.NewRows([]string{"TABLE_SCHEMA", "TABLE_NAME"}).
		AddRow("edu", "student").
		AddRow("edu", "teacher")
	s.mock.ExpectQuery("SELECT * (.+)").WithArgs(dbName, tableName).WillReturnRows(rows)
	pluckRows := sqlmock.NewRows([]string{"TABLE_NAME"}).
		AddRow("student")
	s.mock.ExpectQuery("SELECT TABLE_NAME (.+)").WithArgs(dbName, tableName).WillReturnRows(pluckRows)
	dbSchema = s.Db
	tbs := matchTables("edu", "stu*")
	if tbs == nil {
		s.T().Error("not match")
	} else {
		assert.Equal(s.T(), "student", tbs[0], "select TABLE_NAME should be equal")
	}
}

func (s *gromMockTest) TestGetOneTableColumns() {
	var (
		dbName    = "edu"
		tableName = "student"
	)
	rows := sqlmock.NewRows([]string{"ColumnName", "ColumnKey", "DataType", "IsNullable", "ColumnComment"}).
		AddRow("id", "", "int", "NO", "")

	s.mock.ExpectQuery("SELECT (.+)").WithArgs(dbName, tableName).WillReturnRows(rows)
	dbSchema = s.Db
	tbs := getOneTableColumns("edu", "student")
	if tbs == nil {
		s.T().Error("not match")
	}
	require.NoError(s.T(), nil)
}

func (s *gromMockTest) TestFindStructMapper() {
	var (
		dbName     = "edu"
		tableName  = "student"
		structName = "studentVO"
	)
	rows := sqlmock.NewRows([]string{"id", "db_name", "table_name", "struct_name", "model_field_name", "model_field_type", "deleted_at"}).
		AddRow("1", "edu", "student", "studentVO", "Id", "int", nil).
		AddRow("2", "edu", "student", "studentVO", "RealName", "string", nil).
		AddRow("3", "edu", "student", "studentVO", "job", "string", nil).
		AddRow("4", "edu", "student", "studentVO", "age", "int", nil).
		AddRow("5", "edu", "student", "studentVO", "sex", "int", nil).
		AddRow("6", "edu", "student", "studentVO", "birthday", "*time.time", nil)
	dbGen = s.Db
	s.mock.ExpectQuery("SELECT (.+)").WithArgs(dbName, tableName, structName).WillReturnRows(rows)
	mappers, err := findStructMapper("edu", "student", "studentVO")
	if mappers == nil {
		s.T().Error(err.Error())
	} else {
		assert.Equal(s.T(), len(*mappers), 6, "mappers should be 6")
	}
	require.NoError(s.T(), err)
}

// yaml.go tests
func (s *gromMockTest) TestGenConfigYaml() {
	student := mockTable()
	cmdRequest := mockCmdRequestByTable(student)
	GenConfigYaml(cmdRequest)
}

func (s *gromMockTest) TestReadYamlMap() {
	readYamlMap(".gen-model.yaml", "../model")
}

// cmd tests
func (s *gromMockTest) TestSetDataByViper() {
	g := &CmdRequest{}
	expected := "school"
	viper.Set("gen.searchTableName", expected)
	g.SetDataByViper()
	if g.Gen.SearchTableName != expected {
		s.T().Errorf("val is %s; expected %s", g.Gen.SearchTableName, expected)
	}
}

func (s *gromMockTest) TestGenerateStructFromSelfTable() {
	student := mockTable()
	cmdRequest := mockCmdRequestByTable(student)
	log.Println("getSearchTableName", cmdRequest.Gen.getSearchTableName())
	log.Println("getSearchStructName", cmdRequest.Gen.getSearchStructName())
	log.Println("getTables", cmdRequest.getTables())
	log.Println("getOutDir", cmdRequest.getOutDir())
	structName := camelString(student.TableName)
	modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
	columnProcessor := getProcessorSelfTable(student)
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}

func (s *gromMockTest) TestGenerateStructFromYaml() {
	student := mockTable()
	cmdRequest := mockCmdRequestByTable(student)
	cmdRequest.Gen.ModelSuffix = "VO"
	cmdRequest.Gen.Persist = ""
	cmdRequest.Gen.Source = sourceLocal
	log.Printf("%+v", cmdRequest)
	viper.Set("mock_map", mockFieldMap())
	mp, _ := filepath.Abs(modelDir)
	cmdRequest.Wg.Add(1)
	mkStructFromYaml(cmdRequest, "StudentVO", "model", mp)
}

func (s *gromMockTest) TestGenerateStructFromGenTable() {
	student := mockTable()
	cmdRequest := mockCmdRequestByTable(student)
	cmdRequest.Gen.ModelSuffix = "BO"
	cmdRequest.Gen.Persist = ""
	structName := camelString(student.TableName + cmdRequest.Gen.ModelSuffix)
	modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
	columnProcessor := getProcessorGenTable(student.TableName, mockGenMapper())
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}

// help.go tests
func (s *gromMockTest) TestContainString() {
	t := []string{"aaa", "bbb"}
	ta := containString(t, "aaa")
	assert.Equal(s.T(), ta, true, "should be true")
	tb := containString(t, "ccc")
	assert.Equal(s.T(), tb, false, "should be false")
}

func mockFieldMap() *fieldMap {
	var fields []fieldNameAndType
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
func mockCmdRequestByTable(dt *dealTable) *CmdRequest {
	mp, _ := filepath.Abs(modelDir)
	mkdir(mp)
	genConf := genConfig{SearchTableName: dt.TableName, OutDir: modelDir, Source: sourceSelfTable, Persist: sourceLocal}
	return &CmdRequest{Gen: genConf}
}
func mockGenMapper() *[]structMapper {
	var mapSlice []structMapper
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
