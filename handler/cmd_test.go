package handler

import (
	"github.com/spf13/viper"
	"testing"
)

func TestSetDataByViper(t *testing.T) {
	g := &CmdRequest{}
	expected := "school"
	viper.Set("gen.searchTableName", expected)
	g.SetDataByViper()
	if g.Gen.SearchTableName != expected {
		t.Errorf("val is %s; expected %s", g.Gen.SearchTableName, expected)
	}
}
func TestCreateModelStruct(t *testing.T) {

	// g.CreateModelStruct()
}
