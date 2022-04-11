package main

import (
	"ddns/cf"
	"ddns/config"
	"ddns/ip"
	"fmt"
	"os"
	"time"

	"github.com/robfig/cron"
)

//cloudflare ddns client By:Harry 2022/4/10
var staus chan bool
var Config config.Config
var InitStatus bool = false
var currect string
var currectID string

const RETRY int = 5

func loop() {
	currectTime := time.Now()
	var temp string
	fmt.Println(currectTime)
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("开始检测公网IP是否变动")
	var i int
	for i = 1; i <= RETRY; i++ {
		if Config.IPV6 {
			fmt.Printf("第%v次尝试获取本机公网IP:\n", i)
			address, err := ip.GetIPV6(Config.CN)
			if err != nil {
				println("失败")
				continue
			}
			if Config.CN {
				temp = address.IP_CN
			} else {
				temp = address.IP
			}
			fmt.Printf("%v\n", temp)
			break
		}
		fmt.Printf("第%v次尝试获取本机公网IP:\n", i)
		address, err := ip.GetIPV4(Config.CN)
		if err != nil {
			println("失败")
			continue
		}
		if Config.CN {
			temp = address.IP_CN
		} else {
			temp = address.IP
		}
		fmt.Printf("%v\n", temp)
		break
	}
	if i > RETRY {
		fmt.Println("本次更新失败")
	} else if temp == currect {
		fmt.Println("公网IP未变动,跳过本次更新")
	} else {
		for i = 1; i <= RETRY; i++ {
			fmt.Printf("第%v次尝试更新DNS记录:\n", i)
			up, err := cf.UpdateRecord(Config.Name, Config.ZeroID, currectID, temp, Config.IPV6, Config.Proxied, Config.Token)
			if !up || err != nil {
				println("失败")
				continue
			}
			break
		}
		if i > RETRY {
			fmt.Println("本次更新失败")
		} else {
			currect = temp
			fmt.Println("更新DNS记录成功")
		}
	}
	fmt.Println("---------------------------------------------------------------------")
	staus <- true
}
func initialization() {
	var i int
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("主程序开始初始化")
	fmt.Println("---------------------------------------------------------------------")
	for i = 1; i <= RETRY; i++ {
		if Config.IPV6 {
			fmt.Printf("第%v次尝试获取本机公网IP:\n", i)
			address, err := ip.GetIPV6(Config.CN)
			if err != nil {
				println("失败")
				continue
			}
			if Config.CN {
				currect = address.IP_CN
			} else {
				currect = address.IP
			}
			fmt.Printf("%v\n", currect)
			var j int
			for j = 1; j <= RETRY; j++ {
				fmt.Printf("第%v次尝试获取DNS记录:\n", j)
				code, record, err := cf.GetRecords(Config.Name, Config.ZeroID, Config.Token)
				if err != nil {
					println("失败")
					continue
				}
				//记录已存在并且不一致
				if code && record.Content != currect {
					fmt.Printf("%v: %v\n", record.Name, record.Content)
					fmt.Println("当前IP与DNS记录不一致,开始更新DNS")
					var k int
					for k = 1; k <= RETRY; k++ {
						fmt.Printf("第%v次尝试更新DNS记录:\n", k)
						up, err := cf.UpdateRecord(Config.Name, Config.ZeroID, record.ID, currect, Config.IPV6, Config.Proxied, Config.Token)
						if !up || err != nil {
							println("失败")
							continue
						}
						fmt.Println("更新DNS记录成功")
						currectID = record.ID
						InitStatus = true
						break
					}
					//记录不存在
				} else if !code {
					fmt.Println("DNS记录不存在,尝试创建DNS记录")
					var k int
					for k = 1; k <= RETRY; k++ {
						fmt.Printf("第%v次尝试创建DNS记录:\n", k)
						up, id, err := cf.CreateRecord(Config.Name, Config.ZeroID, currect, Config.IPV6, Config.Proxied, Config.Token)
						if !up || err != nil {
							println("失败")
							continue
						}
						fmt.Println("创建DNS记录成功")
						currectID = id
						InitStatus = true
						break
					}
					//存在且一致
				} else {
					fmt.Printf("%v: %v\n", record.Name, record.Content)
					currectID = record.ID
					InitStatus = true
					break
				}
				break
			}
			break
		}
		fmt.Printf("第%v次尝试获取本机公网IP:\n", i)
		address, err := ip.GetIPV4(Config.CN)
		if err != nil {
			println("失败")
			continue
		}
		if Config.CN {
			currect = address.IP_CN
		} else {
			currect = address.IP
		}
		fmt.Printf("%v\n", currect)
		var j int
		for j = 1; j <= RETRY; j++ {
			fmt.Printf("第%v次尝试获取DNS记录:\n", j)
			code, record, err := cf.GetRecords(Config.Name, Config.ZeroID, Config.Token)
			if err != nil {
				println("失败")
				continue
			}
			//记录已存在并且不一致

			if code && record.Content != currect {
				fmt.Printf("%v: %v\n", record.Name, record.Content)
				fmt.Println("当前IP与DNS记录不一致,开始更新DNS")
				var k int
				for k = 1; k <= RETRY; k++ {
					fmt.Printf("第%v次尝试更新DNS记录:\n", k)
					up, err := cf.UpdateRecord(Config.Name, Config.ZeroID, record.ID, currect, Config.IPV6, Config.Proxied, Config.Token)
					if !up || err != nil {
						println("失败")
						continue
					}
					fmt.Println("更新DNS记录成功")
					currectID = record.ID
					InitStatus = true
					break
				}
				//记录不存在
			} else if !code {
				fmt.Println("DNS记录不存在,尝试创建DNS记录")
				var k int
				for k = 1; k <= RETRY; k++ {
					fmt.Printf("第%v次尝试创建DNS记录:\n", k)
					up, id, err := cf.CreateRecord(Config.Name, Config.ZeroID, currect, Config.IPV6, Config.Proxied, Config.Token)
					if !up || err != nil {
						println("失败")
						continue
					}
					fmt.Println("创建DNS记录成功")
					InitStatus = true
					currectID = id
					break
				}
				//存在且一致
			} else {
				fmt.Printf("%v: %v\n", record.Name, record.Content)
				currectID = record.ID
				InitStatus = true
				break
			}
			break
		}
		break
	}
	if !InitStatus {
		fmt.Println("初始化失败,请检查网络设置")
		os.Exit(0)
	}
	fmt.Println("初始化成功")
	fmt.Println("---------------------------------------------------------------------")
}
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
func main() {
	args := os.Args
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("             Cloudflare DDNS Client By:Harry v1.0")
	fmt.Println("---------------------------------------------------------------------")
	fmt.Println("           github.com/kanomahoro/go-cloudflare-ddns")
	//命令行参数第一项必定为程序的绝对路径
	if len(args) == 1 {
		//尝试从当前目录加载config.json
		status, _ := exists("config.json")
		if status {
			data, err := config.LoadConfig("config.json")
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			Config = data
		} else {
			fmt.Println("没有找到配置文件")
			err := config.CreateConfig()
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			fmt.Println("已在当前目录创建默认配置文件config.json,请重新运行")
			return
		}
	} else if len(args) == 2 {
		//从给定的配置文件启动
		status, _ := exists(args[1])
		if status {
			data, err := config.LoadConfig(args[1])
			if err != nil {
				fmt.Println(err.Error())
				return
			}
			Config = data
		} else {
			fmt.Println("给定的配置文件不存在")
			return
		}
	} else {
		//参数不合法
		fmt.Println("无效的参数")
		return
	}
	//初始化
	initialization()
	timer := cron.New()
	err := timer.AddFunc(Config.Cron, loop)
	//校验cron表达式
	if err != nil {
		fmt.Println("cron表达式不合法,请检测配置文件")
	}
	timer.Start()
	//阻塞主线程，直到用户手动停止程序
	for {
		<-staus
	}
}
