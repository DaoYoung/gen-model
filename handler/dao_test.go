package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type daoMockTest struct {
	suite.Suite
	Db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *daoMockTest) SetupSuite() {
	db, mock, err := sqlmock.New()
	s.mock = mock
	require.NoError(s.T(), err)
	s.Db, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)
	s.Db.LogMode(false)
	viper.Set("is_test", true)
	viper.Set("forceCover", true)
}
func (s *daoMockTest) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func TestDao(t *testing.T) {
	suite.Run(t, new(daoMockTest))
}
func (s *daoMockTest) TestMatchTables() {
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

func (s *daoMockTest) TestGetOneTableColumns() {
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

func (s *daoMockTest) TestFindStructMapper() {
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
func (s *daoMockTest) TestMysqlTypeToGoType() {
	assert.Equal(s.T(), mysqlTypeToGoType("smallint", true), gureguNullInt, "error smallint")
	assert.Equal(s.T(), mysqlTypeToGoType("bigint", true), gureguNullInt, "error bigint")
	assert.Equal(s.T(), mysqlTypeToGoType("double", true), gureguNullFloat, "error double")
	assert.Equal(s.T(), mysqlTypeToGoType("timestamp", false), golangTime, "error timestamp")
	assert.Equal(s.T(), mysqlTypeToGoType("jsonp", false), "", "error enum")
	assert.Equal(s.T(), getImportPackage(""), importNothing, "error importNothing")
}
func (s *daoMockTest) TearDownSuite() {
	viper.Set("debug", true)
	dbSchema = nil
	dbGen = nil
	initGenDb()
}
