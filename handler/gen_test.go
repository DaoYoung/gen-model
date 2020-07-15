package main

import (
    "testing"
    "fmt"
    "github.com/DaoYoung/gen-model/cmd"
    "github.com/DaoYoung/gen-model/handler"
)

func TestGenerateStructFromSelfTable(t *testing.T)  {
    dealTable := &(handler.DealTable{})
    dealTable.TableName = "student"
    studentColumns := []handler.SchemaColumn{
        {ColumnName:"id",ColumnKey:"PRI",DataType:"int(10)", IsNullable:"NO"},
        {ColumnName:"real_name",DataType:"varchar", IsNullable:"NO", ColumnComment:"姓名"},
        {ColumnName:"age",DataType:"mediumint", IsNullable:"NO", ColumnComment:"年龄"},
        {ColumnName:"sex",DataType:"tinyint", IsNullable:"NO", ColumnComment:"性别 1：男 2：女"},
        {ColumnName:"birthday",DataType:"date", IsNullable:"YES", ColumnComment:"生日"},

    }
    dealTable.Columns = &studentColumns
    if len(*dealTable.Columns) == 0 {
        fmt.Println("empty table: " + tableName)
        return
    }
    structName := camelString(tableName + cmdRequest.Gen.ModelSuffix)
    modelPath, packageName := cmdRequest.getAbsPathAndPackageName()
    columnProcessor := getProcessorSelfTable(dealTable, cmdRequest)
    outputStruct(cmdRequest, columnProcessor, modelPath, packageName, structName)
}
func BenchmarkHello(b *testing.B) {
    for i := 0; i < b.N; i++ {
        fmt.Sprintf("hello")
    }
}
func ExampleHello() {
    cmd.InitConfig()
    cmd.CmdRequest.SetDataByViper()
    cmd.CmdRequest.CreateModelStruct()
    // Output: hello
}