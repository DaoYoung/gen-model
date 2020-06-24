package handler

import (
    "fmt"
    "os"
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
        if len(*dealTable.Columns) == 0 {
            fmt.Println("empty table: " + tn)
            continue
        }
        cmdRequest.Wg.Add(1)
        go structWrite(dealTable, cmdRequest)
    }
    cmdRequest.Wg.Wait()
    os.Exit(0)
}

func structWrite(dealTable *dealTable, cmdRequest *CmdRequest) {
    structName := camelString(dealTable.TableName)
    absOutPutPath, packageName := cmdRequest.getAbsPathAndPackageName()
    fileName := checkFile(absOutPutPath + "/" + structName + ".go")
    str := "package " + packageName + "\n\n"
    columnProcessor := columnProcess(dealTable.Columns, cmdRequest)
    str += columnProcessor.ImportSegment
    str += "\ntype " + structName + " struct {"
    str += columnProcessor.AttrSegment
    str += "\n}\n\n"
    str += "func (model *" + structName + ") TableName() string {\n    return \"" + dealTable.TableName + "\"\n}"
    strmodel := fmt.Sprintf("%s", str)
    defer cmdRequest.Wg.Done()
    fmt.Print("\ncreate model " + fileName)
    err := writeFile(fileName, strmodel)
    if err != nil {
        fmt.Print(" failed/n")
        fmt.Println(err.Error())
    }else{
        fmt.Print(" success")
    }
}

func columnProcess(columns *[]SchemaColumn, cmdRequest *CmdRequest) *columnProcessor {
    columnProcessor := &(columnProcessor{})
    var importPackages []string
    var primary string
    for _, column := range *columns {
        var structTags []string
        structAttr := camelString(column.ColumnName)
        fieldType, needPackage := mysqlTypeToGoType(column.DataType, column.IsNull(), cmdRequest.Gen.HasGureguNullPackage)
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

