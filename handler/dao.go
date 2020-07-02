package handler

import (
    "strings"
    "time"
)

type SchemaColumn struct {
    ColumnName string `gorm:"column:COLUMN_NAME"`
    ColumnKey  string `gorm:"column:COLUMN_KEY"`
    DataType   string `gorm:"column:DATA_TYPE"`
    IsNullable string `gorm:"column:IS_NULLABLE"`
}

func (tc *SchemaColumn) TableName() string {
    return "COLUMNS"
}

func (tc *SchemaColumn) IsNull() bool {
    return tc.IsNullable == "YES"
}

type SchemaTable struct {
    TableSchema    string `gorm:"column:TABLE_SCHEMA"`
    TableNameAlias string `gorm:"column:TABLE_NAME"`
}

func (tc *SchemaTable) TableName() string {
    return "TABLES"
}

type dealTable struct {
    TableName string
    Columns   *[]SchemaColumn
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

type structMapper struct {
    Id    int    `gorm:"primary_key;auto_increment;"`
    DbName       string    `gorm:"type:varchar(150);not null;comment:'database name';"`
    TableName     string `gorm:"type:varchar(150);not null;comment:'table name';"`
    ModelFieldName   string `gorm:"type:varchar(150);not null;comment:'golang struct field name';"`
    ModelFieldType    string  `gorm:"type:varchar(50);not null;comment:'golang struct field type';"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt *time.Time
}

func createOrUpdateMappers(dbName string, columnProcessor *columnProcessor) {
    var existFields []string
    condition := &structMapper{}
    condition.DbName = dbName
    condition.TableName = columnProcessor.TableName
    for _,fieldNameAndType := range columnProcessor.Attrs{
        condition.Id = 0
        fn,ft := fieldNameAndType.getValues()
        condition.ModelFieldName = fn
        existFields = append(existFields, fn)
        dbGen.Where(condition).Assign(structMapper{ModelFieldType: ft}).FirstOrCreate(&structMapper{})
    }
    if len(existFields) > 0 {
        dbGen.Where("model_field_name NOT IN (?)", existFields).Delete(structMapper{})
    }
}