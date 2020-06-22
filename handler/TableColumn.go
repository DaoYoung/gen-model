package handler

type TableColumn struct {
    ColumnName string
    ColumnKey  string
    DataType   string
    IsNullable string
}

func (tc *TableColumn) TableName() string {
    return "COLUMNS"
}
