package main

import (
    "github.com/spf13/cobra"
    "fmt"
    "bytes"
    "io/ioutil"
    "testing"
)
var in string
func NewRootCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "hugo",
        Short: "Hugo is a very fast static site generator",
        RunE: func(cmd *cobra.Command, args []string) error {
            fmt.Fprintf(cmd.OutOrStdout(), in)
            return nil
        },
    }
    cmd.Flags().StringVar(&in, "in", "", "This is a very important input.")
    return cmd
}
func Test_Main(t *testing.T)  {
    cmd := NewRootCmd()
    b := bytes.NewBufferString("")
    cmd.SetOut(b)
    cmd.SetArgs([]string{"--in", "testisawesome"})
    cmd.Execute()
    out, err := ioutil.ReadAll(b)
    if err != nil {
        t.Fatal(err)
    }
    if string(out) != "testisawesome" {
        t.Fatalf("expected \"%s\" got \"%s\"", "testisawesome", string(out))
    }
}
