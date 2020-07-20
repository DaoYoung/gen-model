package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

// YamlFile is config file
var YamlFile = ".gen-model"

// YamlMap is local mapper file suffix
var YamlMap = "FieldMapper"

// YamlExt is local file ext
var YamlExt = ".yaml"

type fieldMap struct {
	TableName string
	Fields    []fieldNameAndType // map has sort problem, fix by slice
}
type fieldNameAndType map[string]string

func (f fieldNameAndType) getValues() (fieldName, fieldType string) {
	for fieldName, fieldType = range f {
		break
	}
	return
}

// GenConfigYaml can generate config in app path
func GenConfigYaml(cmdRequest *CmdRequest) {
	projectRoot, _ := os.Getwd()
	fileName := filepath.Join(projectRoot, YamlFile+YamlExt)
	if isExist(fileName) && !viper.GetBool("forceCover") {
		printMessageAndExit("you have config file: " + fileName + ", \nset flag --forceCover=true if you want cover")
	}
	content := ""
	content += "mysql:\n"
	content += "  host: " + cmdRequest.Db.Host + "\n"
	content += "  database: " + cmdRequest.Db.Database + "\n"
	content += "  port: " + strconv.Itoa(cmdRequest.Db.Port) + "\n"
	content += "  username: " + cmdRequest.Db.Username + "\n"
	content += "  password: " + cmdRequest.Db.Password + "\n"
	content += "gen:\n"
	content += "  searchTableName: " + cmdRequest.Gen.SearchTableName + " # support patten with '*'\n"
	content += "  outDir: " + cmdRequest.Gen.OutDir + " # file path\n"
	content += "  jsonUcFirst: " + strconv.FormatBool(cmdRequest.Gen.JSONUcFirst) + " # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'\n"
	content += "  modelSuffix: " + cmdRequest.Gen.ModelSuffix + " # model name suffix\n"
	content += "  source: " + cmdRequest.Gen.Source + " # self-table: struct create by connect mysql tables; local-mapper: struct create by local mappers; db-mapper: struct create by table \"gen_model_mapper\"\n"
	content += "  persist: " + cmdRequest.Gen.Persist + " # persist struct mappers at local-mapper or db-mapper\n"

	fmt.Print("\ncreate yaml " + fileName)
	err := writeFile(fileName, content)
	if err != nil {
		printMessageAndExit(" failed " + err.Error())
	}
	fmt.Print(" success")
	exitWithCode(0)
}

func genMapYaml(filename string, columnProcessor *columnProcessor) error {
	fm := &fieldMap{TableName: columnProcessor.TableName, Fields: columnProcessor.Attrs}
	d, err := yaml.Marshal(&fm)
	if err != nil {
		return err
	}
	return writeFile(filename, fmt.Sprintf("%s", string(d)))
}
func readYamlMap(fileName, modelPath string) *fieldMap {
	data, err := ioutil.ReadFile(filepath.Join(modelPath, fileName+YamlExt))
	if err != nil {
		printErrorAndExit(err)
	}
	fieldMap := &fieldMap{}
	err = yaml.Unmarshal(data, fieldMap)
	if err != nil {
		printErrorAndExit(err)
	}
	return fieldMap
}
