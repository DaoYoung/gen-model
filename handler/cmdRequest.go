package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

type CmdRequest struct {
	Db  dbConfig
	Gen genConfig
	Wg  sync.WaitGroup
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
	HasGormTag           bool   // gorm tag, `gorm:"column:name"`
	HasJsonTag           bool   // json tag, `json:"age"`
	HasGureguNullPackage bool   // have package: "gopkg.in/guregu/null.v3"
	ModelSuffix          string // model name suffix
	SourceType           string // self-table: struct create by connect mysql tables local-mapper: struct create by local mappers gen-table: struct create by table "gen_model_mapper"
	PersistType          string // persist struct mappers at local-mapper or gen-table
}

func (gc *genConfig) getSearchTableName() string {
	return gc.SearchTableName + gc.ModelSuffix
}
func (gc *genConfig) getSearchStructName() string {
	return camelString(gc.SearchTableName) + gc.ModelSuffix
}
func (gc *genConfig) isBuildLocalMapper() bool {
	return gc.PersistType == sourceLocal && gc.SourceType != sourceLocal
}
func (gc *genConfig) isBuildGenTable() bool {
	return gc.PersistType == sourceGenTable && gc.SourceType != sourceGenTable
}

const (
	sourceSelfTable = "self-table"
	sourceLocal     = "local-mapper"
	sourceGenTable  = "gen-table"
)

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
	if absPath, err = filepath.Abs(g.Gen.OutPutPath); err != nil {
		printErrorAndExit(err)
	}
	if !isExist(absPath) {
		printMessageAndExit("OutPutPath not exist: " + absPath)
	}
	if appPath, err = os.Getwd(); err != nil {
		printErrorAndExit(err)
	}
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	log.Println(path)
	log.Println(absPath, appPath)
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
	g.Gen.ModelSuffix = viper.GetString("gen.modelSuffix")
	g.Gen.SourceType = viper.GetString("gen.sourceType")
	g.Gen.PersistType = viper.GetString("gen.persistType")
	g.Db.Host = viper.GetString("mysql.host")
	g.Db.Database = viper.GetString("mysql.database")
	g.Db.Port = viper.GetInt("mysql.port")
	g.Db.Username = viper.GetString("mysql.username")
	g.Db.Password = viper.GetString("mysql.password")
}
func (cmdRequest *CmdRequest) selfTable2Struct() {
	initSchemaDb()
	fmt.Println("search table " + cmdRequest.Gen.SearchTableName + " in db: " + cmdRequest.Db.Database)
	tables := cmdRequest.getTables()
	for _, tn := range tables {
		cmdRequest.Wg.Add(1)
		go mkStructFromSelfTable(tn, cmdRequest)
	}
	cmdRequest.Wg.Wait()
	if len(tables) == 0 {
		fmt.Println("\n\n  nothing found out :( ")
	}
	os.Exit(0)
}
func (g *CmdRequest) localMap2Struct() {
	modelPath, packageName := g.getAbsPathAndPackageName()
	fmt.Println("pattern: " + g.Gen.getSearchStructName() + " search mapper yaml at " + modelPath)
	files, _ := ioutil.ReadDir(modelPath)
	count := 0
	for _, f := range files {
		fn := f.Name()
		suffix := path.Ext(fn)
		if suffix == YamlExt {
			fileName := strings.TrimSuffix(fn, suffix)
			if isFileNameMatch(g.Gen.SearchTableName, g.Gen.ModelSuffix, fileName) {
				count++
				g.Wg.Add(1)
				go mkStructFromYaml(g, fileName, packageName, modelPath)
			}
		}
	}
	g.Wg.Wait()
	if count == 0 {
		fmt.Println("\n\n  none yaml found out :( ")
	}
	os.Exit(0)
}
func (cmdRequest *CmdRequest) genTable2Struct() {
	initSchemaDb()
	initGenDb()
	fmt.Println("search table " + cmdRequest.Gen.getSearchTableName() + " in gen table: get_model.struct_mappers")
	tables := cmdRequest.getTables()
	for _, tn := range tables {
		cmdRequest.Wg.Add(1)
		go mkStructFromGenTable(tn, cmdRequest)
	}
	cmdRequest.Wg.Wait()
	if len(tables) == 0 {
		fmt.Println("\n\n  nothing found out :( ")
	}
	os.Exit(0)
}
func (g *CmdRequest) CreateModelStruct() {
	switch g.Gen.SourceType {
	case sourceSelfTable:
		g.selfTable2Struct()
		break
	case sourceLocal:
		g.localMap2Struct()
		break
	case sourceGenTable:
		g.genTable2Struct()
		break
	default:
		printMessageAndExit("wrong sourceType, set value with \"" + sourceSelfTable + "\" or \"" + sourceLocal + "\" or \"" + sourceGenTable + "\"")
	}
}

func isFileNameMatch(pattern, suffix, fileName string) bool {
	fileName = strings.TrimSuffix(fileName, YamlMap)
	pattern = camelString(pattern) + suffix
	if strings.Contains(pattern, "*") {
		isMatch, _ := regexp.MatchString(strings.Replace(pattern, "*", "(.*)", -1), fileName)
		return isMatch
	}
	return fileName == pattern
}
