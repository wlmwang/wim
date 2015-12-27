package sys

import (
	"fmt"
	"github.com/bitly/go-simplejson"
	"io/ioutil"
	"os"
)

type Config struct {
	simplejson.Json
}

var BaseConf *Config
var LangConf *Config

func NewConfig(filename string) *Config {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("打开配置文件：", err)
		os.Exit(1)
		return nil
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("读配置文件：", err)
		os.Exit(2)
		return nil
	}
	js := new(Config)
	err = js.UnmarshalJSON(content)
	if err != nil {
		fmt.Println("解析配置文件：", err)
		os.Exit(1)
		return nil
	}
	return js
}
