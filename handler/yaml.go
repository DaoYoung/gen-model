package handler

import (
    "log"
    "strconv"
    "os"
    "github.com/spf13/viper"
    "path/filepath"
    "fmt"
    "gopkg.in/yaml.v2"
)

var YamlFile = ".gen-model"


type fieldMap struct {
    TableName string
    Fields map[string]string
}

func GenConfigYaml(cmdRequest *CmdRequest) {
    log.Printf("GenConfigYaml %+v", cmdRequest)
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
    content += "  modelSuffix: " + cmdRequest.Gen.ModelSuffix + " # model name suffix\n"
    content += "  sourceType: " + cmdRequest.Gen.SourceType + " # self-table: struct create by connect mysql tables local: struct create by local mappers gen-table: struct create by table \"gen_model_mapper\"\n"
    content += "  persistType: " + cmdRequest.Gen.PersistType + " # persist struct mappers at local or db\n"
    content += "  isLowerCamelCaseJson: " + strconv.FormatBool(cmdRequest.Gen.IsLowerCamelCaseJson) + " # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'\n"
    content += "  hasGormTag: " + strconv.FormatBool(true) + " # gorm tag, `gorm:\"column:name\"`\n"
    content += "  hasJsonTag: " + strconv.FormatBool(true) + " # json tag, `json:\"age\"`\n"
    content += "  hasGureguNullPackage: " + strconv.FormatBool(cmdRequest.Gen.HasGureguNullPackage) + " # have package: \"gopkg.in/guregu/null.v3\"\n"
    fileName := cmdRequest.getOutPutPath() + "/" + YamlFile + ".yaml"
    fileName =filepath.FromSlash(fileName)
    log.Println("GenConfigYaml: ", fileName)
    if isExist(fileName) && !viper.GetBool("force-cover") {
        printMessageAndExit("you have config file: " + filepath.FromSlash(fileName) + ", \nset falg --force-cover=true if you want cover")
    }
    fmt.Print("\ncreate yaml " + fileName)
    err := writeFile(fileName, content)
    if err != nil {
        printMessageAndExit(" failed/n")
    }
    fmt.Print(" success")
    os.Exit(0)
}

func genMapYaml(tableName string,filename string, columnProcessor *columnProcessor) error {
    fm := &fieldMap{TableName:tableName, Fields:columnProcessor.Attrs}
    d, err := yaml.Marshal(&fm)
    if err != nil {
        return err
    }
    return writeFile(filename, fmt.Sprintf("%s", string(d)))
}