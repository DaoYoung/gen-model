package handler

import (
    "log"
    "strconv"
    "os"
)

var YamlFile = ".gen-model"

func GenConfigYaml(cmdRequest *CmdRequest) {
    log.Println("GenConfigYaml", cmdRequest)
    content := ""
    content += "mysql:\n"
    content += "  host:" + cmdRequest.Db.Host + "\n"
    content += "  database:" + cmdRequest.Db.Database + "\n"
    content += "  port:" + strconv.Itoa(cmdRequest.Db.Port) + "\n"
    content += "  username:" + cmdRequest.Db.Username + "\n"
    content += "  password:" + cmdRequest.Db.Password + "\n"
    content += "gen:\n"
    content += "  searchTableName:" + cmdRequest.Gen.SearchTableName + " # support patten with '*'\n"
    content += "  outPutPath:" + cmdRequest.Gen.OutPutPath + " # file path\n"
    content += "  isLowerCamelCaseJson:" + strconv.FormatBool(cmdRequest.Gen.IsLowerCamelCaseJson) + " # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'\n"
    content += "  hasGormTag:" + strconv.FormatBool(cmdRequest.Gen.HasGormTag) + " # gorm tag, `gorm:\"column:name\"`\n"
    content += "  hasJsonTag:" + strconv.FormatBool(cmdRequest.Gen.HasJsonTag) + " # json tag, `json:\"age\"`\n"
    content += "  hasGureguNullPackage:" + strconv.FormatBool(cmdRequest.Gen.HasGureguNullPackage) + " # have package: \"gopkg.in/guregu/null.v3\"\n"
    fileName := cmdRequest.getOutPutPath() + "/" + YamlFile + ".yaml"
    log.Println("GenConfigYaml: ", fileName)
    if isExist(fileName) {
        log.Println("you must delete config file: " + fileName + ", if you want create new one")
        os.Exit(1)
    }
    writeFile(fileName, content)
}