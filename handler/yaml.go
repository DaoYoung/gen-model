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

var YamlFile = ".gen-model"
var YamlMap = "FieldMapper"
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
	content += "  outPutPath: " + cmdRequest.Gen.OutPutPath + " # file path\n"
	content += "  isLowerCamelCaseJson: " + strconv.FormatBool(cmdRequest.Gen.IsLowerCamelCaseJson) + " # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'\n"
	content += "  hasGormTag: " + strconv.FormatBool(true) + " # gorm tag, `gorm:\"column:name\"`\n"
	content += "  hasJsonTag: " + strconv.FormatBool(true) + " # json tag, `json:\"age\"`\n"
	content += "  hasGureguNullPackage: " + strconv.FormatBool(cmdRequest.Gen.HasGureguNullPackage) + " # have package: \"gopkg.in/guregu/null.v3\"\n"
	content += "  modelSuffix: " + cmdRequest.Gen.ModelSuffix + " # model name suffix\n"
	content += "  sourceType: " + cmdRequest.Gen.SourceType + " # self-table: struct create by connect mysql tables; local-mapper: struct create by local mappers; gen-table: struct create by table \"gen_model_mapper\"\n"
	content += "  persistType: " + cmdRequest.Gen.PersistType + " # persist struct mappers at local-mapper or gen-table\n"

	fmt.Print("\ncreate yaml " + fileName)
	// fmt.Println(content)
	err := writeFile(fileName, content)
	if err != nil {
		printMessageAndExit(" failed " + err.Error())
	}
	fmt.Print(" success")
	os.Exit(0)
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
