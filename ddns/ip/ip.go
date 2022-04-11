package ip

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type IP struct {
	IP    string `json:"IP"`
	IP_CN string `json:"processedString"`
}

//中科大的接口，国外解析可能不稳定
const API_V4_CN string = "https://test.ustc.edu.cn/backend/getIP.php"
const API_V6_CN string = "https://test6.ustc.edu.cn/backend/getIP.php"

//ipify的接口，国内解析可能不稳定
const API_V4 string = "https://api.ipify.org/?format=json"
const API_V6 string = "https://api6.ipify.org/?format=json"

//获取本机公网IPV4地址
func GetIPV4(CN bool) (IP, error) {
	var res IP = IP{"", ""}
	var API string = API_V4
	if CN {
		API = API_V4_CN
	}
	response, err := http.Get(API)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	byteData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(byteData, &res)
	return res, nil
}

//获取本机公网IPV6地址
func GetIPV6(CN bool) (IP, error) {
	var res IP = IP{"", ""}
	var API string = API_V6
	if CN {
		API = API_V6_CN
	}
	response, err := http.Get(API)
	if err != nil {
		return res, err
	}
	defer response.Body.Close()
	byteData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(byteData, &res)
	return res, nil
}
