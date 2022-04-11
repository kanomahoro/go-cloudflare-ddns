package cf

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Record struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}
type Records struct {
	Success bool     `json:"success"`
	Result  []Record `json:"result"`
}
type Cord struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Content  string `json:"content"`
	TTL      int    `json:"ttl"`
	Priority int    `json:"priority"`
	Proxied  bool   `json:"proxied"`
}
type Status struct {
	Success bool   `json:"success"`
	Result  Record `json:"result"`
}

//第一个返回值为记录是否存在，第二个返回值为记录本身
func GetRecords(ZoneName string, ZoneID string, APIKey string) (bool, Record, error) {
	var API string = "https://api.cloudflare.com/client/v4/zones/" + ZoneID + "/dns_records"
	req, _ := http.NewRequest("GET", API, nil)
	data := Records{Success: false, Result: nil}
	client := http.Client{}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIKey)
	response, err := client.Do(req)
	if err != nil {
		return false, Record{}, err
	}
	defer response.Body.Close()
	byteData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(byteData, &data)
	for i, record := range data.Result {
		if record.Name == ZoneName {
			return true, data.Result[i], nil
		}
	}
	return false, Record{}, nil
}

//*更新DNS记录，返回值为是否成功
func UpdateRecord(ZoneName string, ZoneID string, RecordID string, IP string, IPV6 bool, Proxied bool, APIKey string) (bool, error) {
	var API string = "https://api.cloudflare.com/client/v4/zones/" + ZoneID + "/dns_records/" + RecordID
	var Type string
	if IPV6 {
		Type = "AAAA"
	} else {
		Type = "A"
	}
	cord := Cord{Type: Type, Name: ZoneName, Content: IP, TTL: 1, Priority: 10, Proxied: Proxied}
	status := Status{}
	client := http.Client{}
	params, _ := json.Marshal(cord)
	req, _ := http.NewRequest("PUT", API, bytes.NewReader(params))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIKey)
	response, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()
	byteData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(byteData, &status)
	if !status.Success || status.Result.Content != IP {
		return false, nil
	}
	return true, nil
}

//*创建DNS记录，返回值为是否成功
func CreateRecord(ZoneName string, ZoneID string, IP string, IPV6 bool, Proxied bool, APIKey string) (bool, string, error) {
	var API string = "https://api.cloudflare.com/client/v4/zones/" + ZoneID + "/dns_records"
	var Type string
	if IPV6 {
		Type = "AAAA"
	} else {
		Type = "A"
	}
	cord := Cord{Type: Type, Name: ZoneName, Content: IP, TTL: 1, Priority: 10, Proxied: Proxied}
	status := Status{}
	client := http.Client{}
	params, _ := json.Marshal(cord)
	req, _ := http.NewRequest("POST", API, bytes.NewReader(params))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+APIKey)
	response, err := client.Do(req)
	if err != nil {
		return false, "", err
	}
	defer response.Body.Close()
	byteData, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(byteData, &status)
	if !status.Success || status.Result.Content != IP {
		return false, "", nil
	}
	return true, status.Result.ID, nil
}
