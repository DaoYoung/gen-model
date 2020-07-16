package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

// CmdRequest is request arguments manager
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
	SearchTableName string

	OutDir string

	// true: uppercase first letter in json tag
	// (default) false: lowercase first letter in json tag
	JSONUcFirst bool

	// model name suffix
	ModelSuffix string

	// self-table: golang struct create by your tables
	// local-mapper: use local files to create
	// db-mapper: create struct with db table: gen_model.struct_mappers
	Source string

	// local-mapper: save mappers at local files
	// db-mapper: create db table: gen_model.struct_mappers, and save mappers in it
	Persist string
}

func (gc *genConfig) getSearchTableName() string {
	return gc.SearchTableName + gc.ModelSuffix
}
func (gc *genConfig) getSearchStructName() string {
	return camelString(gc.SearchTableName) + gc.ModelSuffix
}
func (gc *genConfig) isBuildLocalMapper() bool {
	return gc.Persist == sourceLocal && gc.Source != sourceLocal
}
func (gc *genConfig) isBuildGenTable() bool {
	return gc.Persist == sourceGenTable && gc.Source != sourceGenTable
}

const (
	sourceSelfTable = "self-table"
	sourceLocal     = "local-mapper"
	sourceGenTable  = "db-mapper"
)

func (g *CmdRequest) getTables() []string {
	if strings.Contains(g.Gen.SearchTableName, "*") {
		return matchTables(g.Db.Database, g.Gen.SearchTableName)
	}
	return []string{g.Gen.SearchTableName}
}

func (g *CmdRequest) getOutDir() string {
	if g.Gen.OutDir == "" {
		g.Gen.OutDir = "model"
	}
	p, _ := filepath.Abs(g.Gen.OutDir)
	outDir := filepath.Dir(p)
	mkdir(outDir)
	return p
}

func (g *CmdRequest) getAbsPathAndPackageName() (absPath, packageName string) {
	if g.Gen.OutDir == "" {
		g.Gen.OutDir = "model"
	}
	var err error
	var appPath string
	if absPath, err = filepath.Abs(g.Gen.OutDir); err != nil {
		printErrorAndExit(err)
	}
	if !isExist(absPath) {
		printMessageAndExit("OutDir not exist: " + absPath)
	}
	if appPath, err = os.Getwd(); err != nil {
		printErrorAndExit(err)
	}
	if absPath == appPath {
		packageName = "main"
	} else {
		_, packageName = filepath.Split(absPath)
	}
	return absPath, packageName
}

// SetDataByViper bind viper value
func (g *CmdRequest) SetDataByViper() {
	g.Gen.SearchTableName = viper.GetString("gen.searchTableName")
	g.Gen.OutDir = viper.GetString("gen.outDir")
	g.Gen.JSONUcFirst = viper.GetBool("gen.jsonUcFirst")
	g.Gen.ModelSuffix = viper.GetString("gen.modelSuffix")
	g.Gen.Source = viper.GetString("gen.source")
	g.Gen.Persist = viper.GetString("gen.persist")
	g.Db.Host = viper.GetString("mysql.host")
	g.Db.Database = viper.GetString("mysql.database")
	g.Db.Port = viper.GetInt("mysql.port")
	g.Db.Username = viper.GetString("mysql.username")
	g.Db.Password = viper.GetString("mysql.password")
}

func (g *CmdRequest) selfTable2Struct() {
	initSchemaDb()
	fmt.Println("search table " + g.Gen.SearchTableName + " in db: " + g.Db.Database)
	tables := g.getTables()
	for _, tn := range tables {
		g.Wg.Add(1)
		go mkStructFromSelfTable(tn, g)
	}
	g.Wg.Wait()
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

func (g *CmdRequest) genTable2Struct() {
	initSchemaDb()
	initGenDb()
	fmt.Println("search table " + g.Gen.getSearchTableName() + " in gen table: get_model.struct_mappers")
	tables := g.getTables()
	for _, tn := range tables {
		g.Wg.Add(1)
		go mkStructFromGenTable(tn, g)
	}
	g.Wg.Wait()
	if len(tables) == 0 {
		fmt.Println("\n\n  nothing found out :( ")
	}
	os.Exit(0)
}

// CreateModelStruct hanlder
func (g *CmdRequest) CreateModelStruct() {
	defer func() {
		if r := recover(); r != nil {
			printErrorMsg(r)
		}
	}()
	switch g.Gen.Source {
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
