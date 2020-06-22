package handler

type SchemaColumn struct {
    ColumnName string
    ColumnKey  string
    DataType   string
    IsNullable string
}

func (tc *SchemaColumn) TableName() string {
    return "COLUMNS"
}

type SchemaTable struct {
    TableSchema string
    TableNameAlias  string `gorm:"column:TABLE_NAME"`
}
func (tc *SchemaTable) TableName() string {
    return "TABLES"
}
