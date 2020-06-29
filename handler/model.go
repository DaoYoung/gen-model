package handler

import (
    "fmt"
    "sort"
    "strings"
    "path/filepath"
)

type columnProcessor struct {
    AttrSegment   string
    ImportSegment string
    Attrs         []fieldNameAndType
    TableName string
    ImportPackages []string
}
func (columnProcessor *columnProcessor) buildImportSegment() {
    if len(columnProcessor.ImportPackages) > 0 {
        sort.Strings(columnProcessor.ImportPackages)
        columnProcessor.ImportSegment = "import (\n"
        for _, p := range columnProcessor.ImportPackages {
            columnProcessor.ImportSegment += "    \"" + p + "\"\n"
        }
        columnProcessor.ImportSegment += ")\n"
    }
}

func mkStructFromSelfTable(tableName string, cmdRequest *CmdRequest) {
    defer cmdRequest.Wg.Done()
    dealTable := &(dealTable{})
    dealTable.TableName = tableName
    dealTable.Columns = getOneTableColumns(cmdRequest.Db.Database, tableName)
    if len(*dealTable.Columns) == 0 {
        fmt.Println("empty table: " + tableName)
        return
    }
    structName := camelString(dealTable.TableName + cmdRequest.Gen.ModelSuffix)
    modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
    columnProcessor := getProcessorSelfTable(dealTable, cmdRequest)
    outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)

}

func getProcessorSelfTable(dealTable *dealTable, cmdRequest *CmdRequest) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    columnProcessor.TableName = dealTable.TableName
    columns := *dealTable.Columns
    for _, column := range columns {
        structAttr := camelString(column.ColumnName)
        fieldType := mysqlTypeToGoType(column.DataType, column.IsNull(), cmdRequest.Gen.HasGureguNullPackage)
        nameAndType := fieldNameAndType{}
        nameAndType[structAttr] = fieldType
        oneFieldProcess(columnProcessor, nameAndType, cmdRequest)
    }
    columnProcessor.buildImportSegment()
    return columnProcessor
}
func getProcessorYaml(cmdRequest *CmdRequest ,mapfileName, modelPath string) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    fieldMap := readYamlMap(mapfileName, modelPath)
    columnProcessor.TableName = fieldMap.TableName
    for _, fieldNameAndType := range fieldMap.Fields {
        oneFieldProcess(columnProcessor, fieldNameAndType, cmdRequest)
    }
    columnProcessor.buildImportSegment()
    return columnProcessor
}

func oneFieldProcess(columnProcessor *columnProcessor, fieldNameAndType fieldNameAndType, cmdRequest *CmdRequest) {
    var structTags []string
    structAttr, fieldType := fieldNameAndType.getValues()
    needPackage := getImportPackage(fieldType)
    columnProcessor.Attrs = append(columnProcessor.Attrs, fieldNameAndType)
    if needPackage != "" && !containString(columnProcessor.ImportPackages, needPackage) {
        columnProcessor.ImportPackages = append(columnProcessor.ImportPackages, needPackage)
    }
    if cmdRequest.Gen.HasGormTag {
        structTags = append(structTags, fmt.Sprintf("gorm:\"column:%s\"", snakeString(structAttr)))
    }
    if cmdRequest.Gen.HasJsonTag {
        structTags = append(structTags, fmt.Sprintf("json:\"%s\"", lcfirst(structAttr)))
    }
    columnProcessor.AttrSegment += fmt.Sprintf("\n    %s %s `%s`",
        structAttr,
        fieldType,
        strings.Join(structTags, " "))
}
func beforeMkStruct(cmdRequest *CmdRequest)  {
    defer cmdRequest.Wg.Done()
}
func outputStruct(cmdRequest *CmdRequest, columnProcessor *columnProcessor,modelPath,packageName,structName string)  {
    var paper string
    defer func() {
        fmt.Print(paper)
    }()
    paper = "\ncreate struct " + structName + ".go"
    fileName, err := mkGolangFile(modelPath, structName)
    if err != nil {
        paper += err.Error()
        return
    }
    str := "package " + packageName + "\n\n"
    str += columnProcessor.ImportSegment
    str += "\ntype " + structName + " struct {"
    str += columnProcessor.AttrSegment
    str += "\n}\n\n"
    str += "func (model *" + structName + ") TableName() string {\n    return \"" + columnProcessor.TableName + "\"\n}"
    err = writeFile(fileName, fmt.Sprintf("%s", str))
    if err != nil {
        paper += " failed!!! " + err.Error()
    } else {
        paper += " success."
        if cmdRequest.Gen.PersistType == persistLocal && cmdRequest.Gen.SourceType != sourceLocal {
            paper += " create mapper " + structName + YamlMap + YamlExt
            mapFileName := filepath.Join(modelPath, structName + YamlMap + YamlExt)
            err = genMapYaml(columnProcessor.TableName, mapFileName, columnProcessor)
            if err != nil {
                paper += " failed!!! " + err.Error()
            } else {
                paper += " success."
            }
        }
    }
}
func mkStructFromYaml(cmdRequest *CmdRequest, mapfileName, packageName, modelPath string)  {
    defer cmdRequest.Wg.Done()
    structName := strings.TrimSuffix(mapfileName, YamlMap)
    columnProcessor := getProcessorYaml(cmdRequest,mapfileName,modelPath)
    outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)




}
