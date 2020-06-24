package handler

import (
    "fmt"
    "log"
    "os"
    "path/filepath"
    "sort"
    "strings"
)

type columnProcessor struct {
    AttrSegment   string
    ImportSegment string
}

func Table2struct(cmdRequest *CmdRequest) {
    tables := cmdRequest.getTables()
    dealTable := &(dealTable{})
    for _, tn := range tables {
        dealTable.TableName = tn
        dealTable.Columns = getOneTableColumns(cmdRequest.Db.Database, tn)
        structWrite(dealTable, cmdRequest)
    }
}


func columnProcess(columns *[]SchemaColumn, hasGormAnnotation, hasJsonAnnotation, hasGureguNullPackage bool) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    var importPackages []string
    for _, column := range *columns {
        var annotations []string
        fieldType, needPackage := mysqlTypeToGoType(column.DataType, column.IsNull(), hasGureguNullPackage)
        if needPackage != "" && !containString(importPackages, needPackage) {
            importPackages = append(importPackages, needPackage)
        }
        if hasGormAnnotation {
            primary := ""
            if column.ColumnKey == "PRI" {
                primary = ";primary_key"
            }
            annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s\"", column.ColumnName, primary))
        }
        if hasJsonAnnotation {
            annotations = append(annotations, fmt.Sprintf("json:\"%s\"", camelString(column.ColumnName)))
        }
        columnProcessor.AttrSegment += fmt.Sprintf("\n    %s %s `%s`",
            camelString(column.ColumnName),
            fieldType,
            strings.Join(annotations, " "))
    }
    if len(importPackages) > 0 {
        sort.Strings(importPackages)
        columnProcessor.ImportSegment = "import (\n"
        for _, p := range importPackages {
            columnProcessor.ImportSegment += "    \"" + p + "\"\n"
        }
        columnProcessor.ImportSegment += ")"
    }
    return columnProcessor
}

func structWrite(dealTable *dealTable, cmdRequest *CmdRequest) {
    structName := camelString(dealTable.TableName)
    absPath, err := filepath.Abs(cmdRequest.Gen.OutPutPath)
    if err != nil {
        log.Println("error OutPutPath: " + cmdRequest.Gen.OutPutPath)
        os.Exit(0)
    }
    if !isExist(absPath) {
        log.Println("OutPutPath not exist: " + absPath)
        os.Exit(0)
    }
    packageName := ""
    appPath, err := os.Getwd()
    fileName := absPath + "/" + structName + ".go"
    if absPath == appPath {
        packageName = "main"
    } else {
        _, packageName = filepath.Split(fileName)

    }
    fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
    fp.Truncate(0)
    if err != nil {
        log.Fatal(err)
    }
    defer fp.Close()
    str := "package " + packageName + "\n\n"

    gormAnnotation := true
    jsonAnnotation := true
    columnProcessor := columnProcess(dealTable.Columns, gormAnnotation, jsonAnnotation, true)
    str += columnProcessor.ImportSegment
    str += "\ntype " + structName + " struct {"
    str += columnProcessor.AttrSegment
    str += "\n}\n\n"
    str += "func (model *" + structName + ") TableName() string {\n    return \"" + dealTable.TableName + "\"\n}"
    strmodel := fmt.Sprintf("%s", str)
    _, err = fp.Write([]byte(strmodel))
    if err != nil {
        log.Fatal(err)
    }
}
