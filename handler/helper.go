package handler

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
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
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
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

func writeFile(fileName, content string) error {
	var err error
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0755)
	f.Truncate(0)
	defer f.Close()
	if err != nil {
		return err
	}
	_, err = f.Write([]byte(content))
	if err != nil {
		return err
	}
	return nil
}
func spellGolangFile(outPutPath, structName string) (fileName string, err error) {
	goFileName := structName + ".go"
	fileName = filepath.Join(outPutPath, goFileName)
	if isExist(fileName) && !viper.GetBool("forceCover") {
		return "", errors.New(" failed!!! " + goFileName + " has exist, set --forceCover=true if you want cover.")
	}
	return fileName, nil
}

// Welcome slogan
func Welcome() {
	slogan := "                                            __     __\r\n   ____ ____  ____     ____ ___  ____  ____/ /__  / /\r\n  / __ `/ _ \\/ __ \\   / __ `__ \\/ __ \\/ __  / _ \\/ /\r\n / /_/ /  __/ / / /  / / / / / / /_/ / /_/ /  __/ /\r\n \\__, /\\___/_/ /_/  /_/ /_/ /_/\\____/\\__,_/\\___/_/\r\n/____/\r\n"
	fmt.Println(slogan)
}

func printMessageAndExit(msg string) {
	printErrorMsg(msg)
	exitWithCode(1)
}

func printErrorAndExit(err error) {
	printErrorMsg(err.Error())
	exitWithCode(1)
}

func printErrorMsg(msg interface{}) {
	fmt.Println("\n\noccur error:")
	fmt.Print("  ")
	fmt.Print(msg)
	fmt.Print("\n")
}

func exitWithCode(c int) {
	if viper.GetBool("is_test") {
		fmt.Print("  os exit")
	} else {
		os.Exit(c)
	}
}
func isFileNameMatch(pattern, suffix, fileName string) bool {
	fileName = strings.TrimSuffix(fileName, YamlMap)
	pattern = camelString(pattern) + suffix
	if strings.Contains(pattern, "*") {
		isMatch, _ := regexp.MatchString(strings.Replace(pattern, "*", "(.*)", -1), fileName)
		return isMatch
	}
	return fileName == pattern
}
