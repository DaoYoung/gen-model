package handler

import (
	"github.com/DaoYoung/gen-model/manager/db"
	"path/filepath"
	"strings"
)

type GenRequest struct {
	DbConfig DbConfig
	SearchTableName            string
	OutPutPath           string
	IsLowerCamelCaseJson bool
}

func (g *GenRequest) getTables() []string {
	if strings.Contains(g.SearchTableName, "*") {
		return matchTables(g.DbConfig.Database, g.SearchTableName)
	}
	return []string{g.SearchTableName}
}

func (g *GenRequest) getOutPutPath() string{
	if g.OutPutPath == "" {
		g.OutPutPath = "model"
	}
	return filepath.Dir(g.OutPutPath)
}

func matchTables(dbName, tableName string) []string {
	var names []string
	columns := &([]SchemaTable{})
	pattern := strings.Replace(tableName, "*", "%", 2)
	if err := db.Db.Where("TABLE_SCHEMA = ?", dbName).Where("TABLE_NAME like ?", pattern).Find(columns).Pluck("TABLE_NAME", &names).Error; err != nil {
		panic(err)
	}
	return names
}

func Table2struct(genRequest *GenRequest) {
	tables := genRequest.getTables()
	dealTable := &(DealTable{})
	for _, tn := range tables {
		dealTable.TableName = tn
		dealTable.Columns = getOneTableColumns(genRequest.DbConfig.Database, tn)
		structWrite(dealTable, genRequest)
	}
}

func getOneTableColumns(dbName, tableName string) *[]SchemaColumn {
	columns := &([]SchemaColumn{})
	if err := db.Db.Where("TABLE_SCHEMA = ?", dbName).Where("TABLE_NAME = ?", tableName).Find(columns).Error; err != nil {
		panic(err)
	}
	return columns
}
