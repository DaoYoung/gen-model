package main

import (
    "testing"
    "os/exec"
    "os"
    "log"
    "fmt"
)
// func TestCreateBySelfTable(t *testing.T)  {
//     cmd.InitConfig()
//     cmd.CmdRequest.SetDataByViper()
//     cmd.CmdRequest.CreateModelStruct()
// }
func Crasher() {
    fmt.Println("Going down in flames!")
    os.Exit(1)

}
func TestCrasher(t *testing.T) {
    log.Println(1111)
    if os.Getenv("BE_CRASHER") == "1" {
        log.Println(1111)
        Crasher()
        return
    }
    cmd := exec.Command(os.Args[0], "-test.run=TestCrasher")
    cmd.Env = append(os.Environ(), "BE_CRASHER=1")
    err := cmd.Run()
    if e, ok := err.(*exec.ExitError); ok && !e.Success() {
        t.Log( e)
        return
    }
    t.Fatalf("process ran with err %v, want exit status 1", err)
}