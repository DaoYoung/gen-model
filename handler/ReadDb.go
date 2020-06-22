package handler

import "github.com/DaoYoung/gen-model/manager/db"

func table2struct()  {
    
}

func getOneTableColumns(dbName ,tableName string) *[]TableColumn {
    columns := &([]TableColumn{})
    if err := db.Db.Where("TABLE_SCHEMA = ?", dbName).Where("TABLE_NAME = ?", tableName).Find(columns).Error; err != nil {
        panic(err)
    }
    return columns
}