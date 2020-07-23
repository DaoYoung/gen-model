package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

var (
	modelDir                 = "../model"
	dbName                   = "college"
	tableName                = "student"
	searchTableName          = "stu*"
	searchTableNameCondition = "stu%"
)

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
	s.Db.LogMode(false)
	viper.Set("is_test", true)
	viper.Set("forceCover", true)
	mockRequest().getOutDir()
}
func (s *gromMockTest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func TestInit(t *testing.T) {
	suite.Run(t, new(gromMockTest))
}

func (s *gromMockTest) TestSourceSelfTableAndGenLocalMapper() {
	// mock table data
	row := sqlmock.NewRows([]string{"TABLE_SCHEMA", "TABLE_NAME"}).
		AddRow(dbName, tableName)
	s.mock.ExpectQuery("SELECT * (.+)").WithArgs(dbName, searchTableNameCondition).WillReturnRows(row)
	table := sqlmock.NewRows([]string{"TABLE_NAME"}).
		AddRow(tableName)
	s.mock.ExpectQuery("SELECT TABLE_NAME FROM (.+)TABLES(.+)").WithArgs(dbName, searchTableNameCondition).WillReturnRows(table)
	fields := sqlmock.NewRows([]string{"COLUMN_NAME", "COLUMN_KEY", "DATA_TYPE", "IS_NULLABLE", "COLUMN_COMMENT"}).
		AddRow("id", "", "bigint", "NO", "").
		AddRow("age", "", "mediumint", "NO", "").
		AddRow("sex", "", "tinyint", "NO", "1:boy 2:girl").
		AddRow("job", "", "varchar", "YES", "job description").
		AddRow("birthday", "", "date", "YES", "").
		AddRow("avatar", "", "blob", "YES", "").
		AddRow("admission_score", "", "decimal", "NO", "").
		AddRow("real_name", "", "varchar", "NO", "job description")
	s.mock.ExpectQuery("SELECT * (.+)COLUMNS(.+)").WithArgs(dbName, tableName).WillReturnRows(fields)
	dbSchema = s.Db // setup db connect
	cr := mockRequest()
	cr.CreateModelStruct()
	cr.Gen.Source = sourceLocal
	cr.Gen.Persist = sourceGenTable
	dbGen = s.Db
	cr.CreateModelStruct()
}

func (s *gromMockTest) TestSourceGenTable() {
	cr := mockRequest()
	cr.Gen.Source = sourceGenTable
	cr.Gen.ModelSuffix = "VO"
	cr.Gen.SearchTableName = tableName
	rows := sqlmock.NewRows([]string{"id", "db_name", "table_name", "struct_name", "model_field_name", "model_field_type", "deleted_at"}).
		AddRow("1", dbName, tableName, "studentVO", "Id", "int", nil).
		AddRow("2", dbName, tableName, "studentVO", "RealName", "string", nil).
		AddRow("3", dbName, tableName, "studentVO", "Job", "string", nil).
		AddRow("4", dbName, tableName, "studentVO", "Age", "int", nil).
		AddRow("5", dbName, tableName, "studentVO", "Sex", "int", nil).
		AddRow("6", dbName, tableName, "studentVO", "Birthday", "*time.Time", nil)
	dbGen = s.Db
	dbSchema = s.Db
	s.mock.ExpectQuery("SELECT * (.+)struct_mappers(.+)").WithArgs(dbName, tableName, "StudentVO").WillReturnRows(rows)
	cr.CreateModelStruct()
}

func mockRequest() *CmdRequest {
	dbConf := dbConfig{Host: "test", Database: dbName, Username: "root", Password: "test", Port: 3306}
	genConf := genConfig{SearchTableName: searchTableName, OutDir: modelDir, Source: sourceSelfTable, Persist: sourceLocal}
	return &CmdRequest{Gen: genConf, Db: dbConf}
}
func (s *gromMockTest) TestGenConfigYaml() {
	cr := mockRequest()
	cr.getOutDir()
	GenConfigYaml(cr)
	viper.Set("gen.source", "wrong_source")
	cr.SetDataByViper()
	cr.CreateModelStruct()
}
