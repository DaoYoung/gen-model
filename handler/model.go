package handler

import (
    "fmt"
    "os"
    "sort"
    "strings"
    "path/filepath"
)

type columnProcessor struct {
    AttrSegment   string
    ImportSegment string
    Attrs         map[string]string
}

func Table2struct(cmdRequest *CmdRequest) {
    tables := cmdRequest.getTables()
    dealTable := &(dealTable{})
    for _, tn := range tables {
        dealTable.TableName = tn
        dealTable.Columns = getOneTableColumns(cmdRequest.Db.Database, tn)
        if len(*dealTable.Columns) == 0 {
            fmt.Println("empty table: " + tn)
            continue
        }
        cmdRequest.Wg.Add(1)
        go structWrite(*dealTable, cmdRequest)
    }
    cmdRequest.Wg.Wait()
    os.Exit(0)
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
    columnProcessor.Attrs = make(map[string]string)
    for _, column := range *columns {
        var structTags []string
        structAttr := camelString(column.ColumnName)
        fieldType, needPackage := mysqlTypeToGoType(column.DataType, column.IsNull(), cmdRequest.Gen.HasGureguNullPackage)
        columnProcessor.Attrs[structAttr] = fieldType
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
