package handler

import (
    "strings"
    "os"
    "unicode"
    "log"
)

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

func lcfirst(str string) string {
    for i, v := range str {
        return string(unicode.ToLower(v)) + str[i+1:]
    }
    return ""
}

func containString(s []string, e string) bool {
    for _, a := range s {
        if strings.EqualFold(a, e) {
            return true
        }
    }
    return false
}

func mkdir(path string) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        os.Mkdir(path, 0777)
        os.Chmod(path, 0777)
    }
}
func isExist(path string) bool {
    _, err := os.Stat(path)
    return err == nil || os.IsExist(err)
}

func writeFile(fileName, content string)  {
    f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
    f.Truncate(0)
    defer f.Close()
    if err != nil {
        log.Println(err.Error())
        os.Exit(1)
    } else {
        _, err = f.Write([]byte(content))
        if err != nil {
            log.Fatal(err)
            os.Exit(1)
        }
    }
}