package handler

import (
    "github.com/spf13/viper"
    "strings"
    "path/filepath"
)

type CmdRequest struct {
    Db  dbConfig
    Gen genConfig
}

type dbConfig struct {
    Host     string
    Database string
    Username string
    Password string
    Port     int
}

type genConfig struct {
    SearchTableName      string
    OutPutPath           string
    IsLowerCamelCaseJson bool
    HasGormTag           bool // gorm tag, `gorm:"column:name"`
    HasJsonTag           bool // json tag, `json:"age"`
    HasGureguNullPackage bool // have package: "gopkg.in/guregu/null.v3"
}

func (g *CmdRequest) getTables() []string {
    if strings.Contains(g.Gen.SearchTableName, "*") {
        return matchTables(g.Db.Database, g.Gen.SearchTableName)
    }
    return []string{g.Gen.SearchTableName}
}

func (g *CmdRequest) getOutPutPath() string {
    if g.Gen.OutPutPath == "" {
        g.Gen.OutPutPath = "model"
    }
    p,_ := filepath.Abs(g.Gen.OutPutPath)
    outDir := filepath.Dir(p)
    mkdir(outDir)
    return p
}

func (g *CmdRequest) SetDataByViper() {
    g.Gen.SearchTableName = viper.GetString("searchTableName")
    g.Gen.OutPutPath = viper.GetString("outPutPath")
    g.Gen.IsLowerCamelCaseJson = viper.GetBool("isLowerCamelCaseJson")
    g.Db.Host = viper.GetString("mysql.host")
    g.Db.Database = viper.GetString("mysql.database")
    g.Db.Port = viper.GetInt("mysql.port")
    g.Db.Username = viper.GetString("mysql.username")
    g.Db.Password = viper.GetString("mysql.password")
}
