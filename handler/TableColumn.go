package handler

type TableColumn struct {
    // TableSchema string
    // TableName string
    ColumnName string
    ColumnKey  string
    DataType   string
    IsNullable string
}

func (tc *TableColumn) TableName() string {
    return "COLUMNS"
}
