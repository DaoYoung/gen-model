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
}

func structWrite(dealTable dealTable, cmdRequest *CmdRequest) {
    var paper string
    defer cmdRequest.Wg.Done()
    defer func() {
        fmt.Print(paper)
    }()
    structName := camelString(dealTable.TableName + cmdRequest.Gen.ModelSuffix)
    paper = "\ncreate struct " + structName + ".go"
    absOutPutPath, packageName := cmdRequest.getAbsPathAndPackageName()
    fileName, err := mkGolangFile(absOutPutPath, structName)
    if err != nil {
        paper += err.Error()
        return
    }
    str := "package " + packageName + "\n\n"
    columnProcessor := columnProcess(dealTable.Columns, cmdRequest)
    str += columnProcessor.ImportSegment
    str += "\ntype " + structName + " struct {"
    str += columnProcessor.AttrSegment
    str += "\n}\n\n"
    str += "func (model *" + structName + ") TableName() string {\n    return \"" + dealTable.TableName + "\"\n}"
    err = writeFile(fileName, fmt.Sprintf("%s", str))
    if err != nil {
        paper += " failed!!! " + err.Error()
    } else {
        paper += " success."
        if cmdRequest.Gen.PersistType == persistLocal && cmdRequest.Gen.SourceType != sourceLocal {
            paper += " create mapper " + structName + "FieldMapper.yaml"
            mapFileName := filepath.Join(absOutPutPath, structName+"FieldMapper.yaml")
            err = genMapYaml(dealTable.TableName, mapFileName, columnProcessor)
            if err != nil {
                paper += " failed!!! " + err.Error()
            } else {
                paper += " success."
            }
        }
    }
}

func columnProcess(columns *[]SchemaColumn, cmdRequest *CmdRequest) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    var importPackages []string
    var primary string
    // columnProcessor.Attrs = make(map[string]string)
    for _, column := range *columns {
        var structTags []string
        structAttr := camelString(column.ColumnName)
        fieldType, needPackage := mysqlTypeToGoType(column.DataType, column.IsNull(), cmdRequest.Gen.HasGureguNullPackage)
        nameAndType := fieldNameAndType{}
        nameAndType[structAttr] = fieldType
        columnProcessor.Attrs = append(columnProcessor.Attrs, nameAndType)
        if needPackage != "" && !containString(importPackages, needPackage) {
            importPackages = append(importPackages, needPackage)
        }
        if cmdRequest.Gen.HasGormTag {
            primary = ""
            if column.ColumnKey == "PRI" {
                primary = ";primary_key"
            }
            structTags = append(structTags, fmt.Sprintf("gorm:\"column:%s%s\"", column.ColumnName, primary))
        }
        if cmdRequest.Gen.HasJsonTag {
            structTags = append(structTags, fmt.Sprintf("json:\"%s\"", lcfirst(structAttr)))
        }
        columnProcessor.AttrSegment += fmt.Sprintf("\n    %s %s `%s`",
            structAttr,
            fieldType,
            strings.Join(structTags, " "))
    }
    if len(importPackages) > 0 {
        sort.Strings(importPackages)
        columnProcessor.ImportSegment = "import (\n"
        for _, p := range importPackages {
            columnProcessor.ImportSegment += "    \"" + p + "\"\n"
        }
        columnProcessor.ImportSegment += ")\n"
    }
    return columnProcessor
}
func columnProcessYaml(cmdRequest *CmdRequest ,mapfileName, modelPath string) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    var importPackages []string
    fieldMap := readYamlMap(mapfileName, modelPath)
    columnProcessor.TableName = fieldMap.TableName
    for _, fieldNameAndType := range fieldMap.Fields {
        var structTags []string
        structAttr, fieldType := fieldNameAndType.getValues()
        fmt.Println(structAttr,fieldType)
        needPackage := getImportPackage(fieldType)
        columnProcessor.Attrs = append(columnProcessor.Attrs, fieldNameAndType)
        if needPackage != "" && !containString(importPackages, needPackage) {
            importPackages = append(importPackages, needPackage)
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
    if len(importPackages) > 0 {
        sort.Strings(importPackages)
        columnProcessor.ImportSegment = "import (\n"
        for _, p := range importPackages {
            columnProcessor.ImportSegment += "    \"" + p + "\"\n"
        }
        columnProcessor.ImportSegment += ")\n"
    }
    return columnProcessor
}

func mkStructFromYaml(cmdRequest *CmdRequest, mapfileName, packageName, modelPath string)  {
    var paper string
    defer cmdRequest.Wg.Done()
    defer func() {
        fmt.Print(paper)
    }()
    structName := strings.TrimSuffix(mapfileName,"FieldMapper")
    paper = "\ncreate struct " + structName + ".go"
    fileName, err := mkGolangFile(modelPath, structName)
    if err != nil {
        paper += err.Error()
        return
    }
    str := "package " + packageName + "\n\n"
    columnProcessor := columnProcessYaml(cmdRequest,mapfileName,modelPath)
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
    }

}
