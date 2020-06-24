package handler

import (
    "log"
    "strconv"
    "os"
    "github.com/spf13/viper"
    "path/filepath"
)

var YamlFile = ".gen-model"

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
    content += "  isLowerCamelCaseJson: " + strconv.FormatBool(cmdRequest.Gen.IsLowerCamelCaseJson) + " # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'\n"
    content += "  hasGormTag: " + strconv.FormatBool(true) + " # gorm tag, `gorm:\"column:name\"`\n"
    content += "  hasJsonTag: " + strconv.FormatBool(true) + " # json tag, `json:\"age\"`\n"
    content += "  hasGureguNullPackage: " + strconv.FormatBool(cmdRequest.Gen.HasGureguNullPackage) + " # have package: \"gopkg.in/guregu/null.v3\"\n"
    fileName := cmdRequest.getOutPutPath() + "/" + YamlFile + ".yaml"
    fileName =filepath.FromSlash(fileName)
    log.Println("GenConfigYaml: ", fileName)
    if isExist(fileName) && !viper.GetBool("force-over") {
        log.Println("you have config file: " + filepath.FromSlash(fileName) + ", \nset falg --force-over=true if you want cover")
        os.Exit(1)
    }
    writeFile(fileName, content)
}