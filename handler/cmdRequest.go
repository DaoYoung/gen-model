package handler

import (
    "github.com/spf13/viper"
    "strings"
    "path/filepath"
    "log"
    "os"
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
    p, _ := filepath.Abs(g.Gen.OutPutPath)
    outDir := filepath.Dir(p)
    mkdir(outDir)
    return p
}

func (g *CmdRequest) getAbsPathAndPackageName() (absPath, packageName string) {
    if g.Gen.OutPutPath == "" {
        g.Gen.OutPutPath = "model"
    }
    var err error
    var appPath string
    if absPath, err = filepath.Abs(g.Gen.OutPutPath);err != nil {
        log.Println(err)
        os.Exit(1)
    }
    if !isExist(absPath) {
        log.Println("OutPutPath not exist: " + absPath)
        os.Exit(1)
    }
    if appPath, err = os.Getwd(); err != nil {
        log.Println(err)
        os.Exit(1)
    }
    if absPath == appPath {
        packageName = "main"
    } else {
        _, packageName = filepath.Split(absPath)
    }
    return absPath, packageName
}

func (g *CmdRequest) SetDataByViper() {
    g.Gen.SearchTableName = viper.GetString("gen.searchTableName")
    g.Gen.OutPutPath = viper.GetString("gen.outPutPath")
    g.Gen.IsLowerCamelCaseJson = viper.GetBool("gen.isLowerCamelCaseJson")
    g.Db.Host = viper.GetString("mysql.host")
    g.Db.Database = viper.GetString("mysql.database")
    g.Db.Port = viper.GetInt("mysql.port")
    g.Db.Username = viper.GetString("mysql.username")
    g.Db.Password = viper.GetString("mysql.password")
}
