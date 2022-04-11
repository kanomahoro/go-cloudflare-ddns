package config

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
)

type Config struct {
	Token   string `json:"token"`   //cloudflare global api key
	ZeroID  string `json:"zeroid"`  //目标域名的ID
	Name    string `json:"name"`    //子域的名字
	Cron    string `json:"cron"`    //用于实现自动更新的cron表达式
	CN      bool   `json:"cn"`      //中国加速
	IPV6    bool   `json:"ipv6"`    //是否解析IPV6
	Proxied bool   `json:"proxied"` //是否启用CDN
}

//加载配置文件
func LoadConfig(file string) (Config, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return Config{"", "", "", "", false, false, false}, errors.New("无法载入配置文件")
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	if config.ZeroID == "" || config.Token == "" || config.Name == "" || config.Cron == "" {
		return config, errors.New("无效的配置文件")
	}
	return config, nil
}

//创建默认配置文件
func CreateConfig() error {
	jsonFile, err := os.OpenFile("config.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, _ := json.Marshal(Config{"", "", "", "", false, false, false})
	jsonFile.Write(byteValue)
	return nil
}
