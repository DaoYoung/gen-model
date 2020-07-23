package handler

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"testing"
)

type gromMockExtra struct {
	suite.Suite
	Db   *gorm.DB
	mock sqlmock.Sqlmock
}

func (s *gromMockExtra) SetupSuite() {
	Welcome()
	db, mock, err := sqlmock.New()
	s.mock = mock
	require.NoError(s.T(), err)
	s.Db, err = gorm.Open("mysql", db)
	require.NoError(s.T(), err)
	s.Db.LogMode(false)
	viper.Set("is_test", true)
	mockRequest().getOutDir()
}
func (s *gromMockExtra) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}
func TestExtra(t *testing.T) {
	suite.Run(t, new(gromMockExtra))
}

func (s *gromMockExtra) aTestSourceSelfTableFail() {
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
	cr.CreateModelStruct()
	cr.Gen.ModelSuffix = "VO"
	cr.CreateModelStruct()
}
func (s *gromMockExtra) TestPathError() {
	cr := mockRequest()
	cr.Gen.OutDir = ""
	cr.getAbsPathAndPackageName()
	cr.Gen.OutDir = ":?>1"
	cr.getAbsPathAndPackageName()
	cr.Gen.OutDir = "../"
	cr.getAbsPathAndPackageName()

}
func (s *gromMockExtra) TestYamlError() {
	cr := mockRequest()
	GenConfigYaml(cr)
	readYamlMap(".fault.yaml", "test")
}
