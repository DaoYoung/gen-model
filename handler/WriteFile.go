package handler

import (
	"log"
	"os"
	"strings"
)

type DealTable struct {
	TableName string
	Columns   *[]SchemaColumn
}

func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}
func jsonWrite(data []byte) {
	fp, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	_, err = fp.Write(data)
	if err != nil {
		log.Fatal(err)
	}
}
func structWrite(dealTable *DealTable, genRequest *GenRequest) {
	structName := camelString(dealTable.TableName)
	paths := strings.Split(genRequest.OutPutPath, "/")
	packageName := paths[len(paths)-1]
	fileName := genRequest.OutPutPath + "/" + structName + ".go"
	log.Println(fileName)
	fp, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	str := "package " + packageName
	str += "type " + structName + " struct {"
	str += "}"
	str += "func (tc *" + structName + ") TableName() string {return \"" + dealTable.TableName + "\"}"

	_, err = fp.Write([]byte(str))
	if err != nil {
		log.Fatal(err)
	}
}
