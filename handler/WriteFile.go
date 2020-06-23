package handler

import (
    "log"
    "os"
    "strings"
    "strconv"
)
var Yamlfile = ".gen-model"
type DealTable struct {
    TableName string
    Columns   *[]SchemaColumn
}

func camelString(s string) string {
    data := make([]byte, 0, len(s))
    j := false
    k := false
    num := len(s) - 1
    for i := 0; i <= num; i++ {
        d := s[i]
        if k == false && d >= 'A' && d <= 'Z' {
            k = true
        }
        if d >= 'a' && d <= 'z' && (j || k == false) {
            d = d - 32
            j = false
            k = true
        }
        if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
            j = true
            continue
        }
        data = append(data, d)
    }
    return string(data[:])
}
func jsonWrite(data []byte) {
    fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        log.Fatal(err)
    }
    defer fp.Close()
    _, err = fp.Write(data)
    if err != nil {
        log.Fatal(err)
    }
}
func structWrite(dealTable *DealTable, genRequest *GenRequest) {
    structName := camelString(dealTable.TableName)
    paths := strings.Split(genRequest.OutPutPath, "/")
    packageName := paths[len(paths)-1]
    fileName := genRequest.OutPutPath + "/" + structName + ".go"
    log.Println(fileName)
    fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
    if err != nil {
        log.Fatal(err)
    }
    defer fp.Close()
    str := "package " + packageName
    str += "type " + structName + " struct {"
    str += "}"
    str += "func (tc *" + structName + ") TableName() string {return \"" + dealTable.TableName + "\"}"

    _, err = fp.Write([]byte(str))
    if err != nil {
        log.Fatal(err)
    }
}
func GenConfigYaml(genRequest *GenRequest) {
    log.Println(genRequest)
    content :=
`mysql:
  host: ` + genRequest.DbConfig.Host + `
  database: ` + genRequest.DbConfig.Database + `
  port: ` + strconv.Itoa(genRequest.DbConfig.Port) + `
  username: ` + genRequest.DbConfig.Username + `
  password: ` + genRequest.DbConfig.Password + `
searchTableName: ` + genRequest.SearchTableName + ` # support patten with '*'
outPutPath: ` + genRequest.OutPutPath + ` # file path
isLowerCamelCaseJson: ` + strconv.FormatBool(genRequest.IsLowerCamelCaseJson) + ` # true: model json tag use lower camelcase, like 'camelCase', not like 'CamelCase'
`
    fileName := genRequest.getOutPutPath()+ "/" +  Yamlfile + ".yaml"
    log.Println("fileName: ", fileName)
    if isExist(fileName) {
        log.Println("you must delete config file: "+fileName+", if you want create new one")
    }
    mkdir(genRequest.getOutPutPath())
    f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
    defer f.Close()
    if err != nil {
        log.Println(err.Error())
    } else {
        _, err = f.Write([]byte(content))
    }
}
func mkdir(path string)  {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.Mkdir(path, 0777)
        os.Chmod(path, 0777)
    }
}
func isExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}