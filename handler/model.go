package handler

import (
	"fmt"
	"github.com/spf13/viper"
	"path/filepath"
	"sort"
	"strings"
)

type columnProcessor struct {
	AttrSegment    string
	ImportSegment  string
	Attrs          []fieldNameAndType
	TableName      string
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

func mkStructFromGenTable(tableName string, cmdRequest *CmdRequest) {
	defer cmdRequest.Wg.Done()
	structName := camelString(tableName + cmdRequest.Gen.ModelSuffix)
	mapSlice, err := findStructMapper(cmdRequest.Db.Database, tableName, structName)
	if err != nil {
		fmt.Println(tableName + ": " + err.Error())
		return
	}
	modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
	columnProcessor := getProcessorGenTable(tableName, mapSlice, cmdRequest)
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
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
	structName := camelString(tableName + cmdRequest.Gen.ModelSuffix)
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
		fieldType := mysqlTypeToGoType(column.DataType, column.isNull(), cmdRequest.Gen.HasGureguNullPackage)
		nameAndType := fieldNameAndType{}
		nameAndType[structAttr] = fieldType
		oneFieldProcess(columnProcessor, nameAndType, cmdRequest)
	}
	columnProcessor.buildImportSegment()
	return columnProcessor
}
func getProcessorYaml(cmdRequest *CmdRequest, mapfileName, modelPath string) *columnProcessor {
	columnProcessor := &(columnProcessor{})
	var fm *fieldMap
	var ok bool
	if fm, ok = viper.Get("mock_map").(*fieldMap); !ok {
		fm = readYamlMap(mapfileName, modelPath)
	}
	columnProcessor.TableName = fm.TableName
	for _, fieldNameAndType := range fm.Fields {
		oneFieldProcess(columnProcessor, fieldNameAndType, cmdRequest)
	}
	columnProcessor.buildImportSegment()
	return columnProcessor
}
func getProcessorGenTable(tableName string, mapSlice *[]structMapper, cmdRequest *CmdRequest) *columnProcessor {
	columnProcessor := &(columnProcessor{})
	columnProcessor.TableName = tableName
	for _, sm := range *mapSlice {
		nameAndType := fieldNameAndType{}
		nameAndType[sm.ModelFieldName] = sm.ModelFieldType
		oneFieldProcess(columnProcessor, nameAndType, cmdRequest)
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
func saveStructMappers(cmdRequest *CmdRequest, columnProcessor *columnProcessor, structName, modelPath string) (paper string) {
	var err error
	hasBuildAction := false
	if cmdRequest.Gen.isBuildLocalMapper() {
		hasBuildAction = true
		paper += " create mapper " + structName + YamlMap + YamlExt
		mapFileName := filepath.Join(modelPath, structName+YamlMap+YamlExt)
		err = genMapYaml(mapFileName, columnProcessor)
	}
	if cmdRequest.Gen.isBuildGenTable() {
		hasBuildAction = true
		initGenDb()
		paper += " mapper sql: insert into gen_model.struct_mappers"
		err = createOrUpdateMappers(viper.GetString("mysql.database"), structName, columnProcessor)
	}
	if hasBuildAction {
		if err != nil {
			paper += " failed!!! " + err.Error()
		} else {
			paper += " success."
		}
	}
	return paper
}
func outputStruct(cmdRequest *CmdRequest, columnProcessor *columnProcessor, modelPath, packageName, structName string) {
	var paper string
	defer func() {
		fmt.Print(paper)
	}()
	paper = "\ncreate struct " + structName + ".go"
	fileName, existErr := mkGolangFile(modelPath, structName)

	str := "package " + packageName + "\n\n"
	str += columnProcessor.ImportSegment
	str += "\ntype " + structName + " struct {"
	str += columnProcessor.AttrSegment
	str += "\n}\n\n"
	str += "func (model *" + structName + ") TableName() string {\n    return \"" + columnProcessor.TableName + "\"\n}"
	if existErr != nil {
		paper += existErr.Error()
		paper += "\n\n------- print " + structName + " start -------\n"
		paper += str
		paper += "\n\n------- print " + structName + " end -------\n"
		return
	}
	err := writeFile(fileName, fmt.Sprintf("%s", str))
	if err != nil {
		paper += " failed!!! " + err.Error()
	} else {
		paper += " success."
		paper += saveStructMappers(cmdRequest, columnProcessor, structName, modelPath)
	}
}
func mkStructFromYaml(cmdRequest *CmdRequest, mapfileName, packageName, modelPath string) {
	defer cmdRequest.Wg.Done()
	structName := strings.TrimSuffix(mapfileName, YamlMap)
	columnProcessor := getProcessorYaml(cmdRequest, mapfileName, modelPath)
	outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}
