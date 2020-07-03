package main

import (
    "testing"
    "fmt"
    "github.com/DaoYoung/gen-model/cmd"
)

// func TestCreateBySelfTable(t *testing.T)  {
//     cmd.InitConfig()
//     cmd.CmdRequest.SetDataByViper()
//     cmd.CmdRequest.CreateModelStruct()
// }
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