package handler

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
