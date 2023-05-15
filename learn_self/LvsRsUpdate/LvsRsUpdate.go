package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Api_key: EP
//Secure_key: di@ep$di
var (
	ApiKey  string = "EP"
	secret  string = "di@ep$di"
	sginUri string = "/lvs/outer/rs_update"

	rsUpdateUri     string = "http://autoproxy.sys.xiaojukeji.com:8009/lvs/outer/rs_update"
	rsAddDelUri     string = "http://autoproxy.sys.xiaojukeji.com:8009/lvs/outer/rs_add_del"
	rsAddDelsginUri string = "/lvs/outer/rs_update"
)

func ComputeHmac256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	return base64.RawURLEncoding.EncodeToString(h.Sum(nil))

}

func main() {
	rsAddDel()
}

func updateRsConfig() {
	// 取当前时间戳 并转化为字符串
	currentTime := time.Now().Unix()
	currentTimeStamp := strconv.FormatInt(currentTime, 10)
	fmt.Println(currentTimeStamp)
	fmt.Printf("%s\n", currentTimeStamp)

	message := fmt.Sprintf("%s_%s_%s_%s", ApiKey, secret, sginUri, currentTimeStamp)
	fmt.Println(message)

	rsListSign := ComputeHmac256(message, secret)
	// 加密签名
	fmt.Println(rsListSign)

	url := fmt.Sprintf("%s?key=%s&sign=%s&timestamp=%s", rsUpdateUri, ApiKey, rsListSign, currentTimeStamp)
	fmt.Println(url)

	//type MemberDetail struct {
	//	Cluster            string `json:"cluster,omitempty"`
	//	Warehouse          string `json:"warehouse,omitempty"`
	//	Product            string `json:"product,omitempty"`
	//	ProductName        string `json:"productname,omitempty"`
	//	Ip                 string `json:"ip,omitempty"`
	//	Port               string `json:"port,omitempty"`
	//	Protocol           string `json:"protocol,omitempty"`
	//	Checkrs            string `json:"check,omitempty"`
	//	Path               string `json:"path,omitempty"`
	//	Connect_timeout    int    `json:"connect_timeout,omitempty"`
	//	Connect_port       string `json:"connect_port,omitempty"`
	//	Weight             int    `json:"weight,omitempty"`
	//	Nb_get_retry       int    `json:"nb_get_retry,omitempty"`       /* dgw连接rs的重试次数 */
	//	Inhibit_on_failure bool   `json:"inhibit_on_failure,omitempty"` /* 以前给LVS的参数,已废除 */
	//	Algorithm          string `gorm:"column:lb_algo" json:"lb_algo,omitempty"`
	//	PortRange          string `gorm:"column:port_range" json:"portRange,omitempty"` /* port范围 */
	//
	//	PortInt     *uint32 `json:"-"`
	//	ConnectPort *uint32 `json:"-"`
	//}
	//
	//type VipReqAPI struct {
	//	Ip        string `json:"ip"`
	//	Port      string `json:"port"`
	//	Protocol  string `json:"protocol"`
	//	PortRange string `gorm:"column:port_range" json:"portRange,omitempty"` /* port范围 */
	//}
	//
	//type MemUP struct { /* 一个rs可以挂载到多个vs下,故目测表示为member中的每个rs都属于Vip中指定的vs下 */
	//	Vip    []VipReqAPI
	//	Member []MemberDetail
	//}
	//
	//type MemberUpdateReqAPI struct {
	//	Cluster     string
	//	Warehouse   string
	//	Product     string
	//	ProductName string `json:"productname,omitempty"`
	//	Member      []MemUP
	//}

	// "reason":[{"ip":"10.90.49.45","port":"3307","vip":[{"cluster":"ilvs02_GZ-YS","warehouse":"GZ-YS","product":"001069","ip":"10.88.151.136","port":"3306",
	//"protocol":"TCP","check":"tcp","connect_timeout":5,"connect_port":"3307","weight":100,"nb_get_retry":3,"lb_algo":"rr"}]}]}

	payload := strings.NewReader(`{
    "warehouse": "GZ-YS",
    "cluster": "ilvs02_GZ-YS",
    "product": "001069",
	"productname": "001069",
    "member" : [
        {
            "vip": [
                {
                    "ip": "10.88.151.136",
                    "port": "3306",
                    "protocol": "tcp"
                }
            ],
            "member": [
                {
                    "ip": "10.90.49.44",
                    "port": "3306",
                    "weight": 100
                }
            ]
        }
    ]
}`)
	client := &http.Client{}
	req, err := http.NewRequest("POST", rsUpdateUri, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

}

func rsAddDel() {
	// 取当前时间戳 并转化为字符串
	currentTime := time.Now().Unix()
	currentTimeStamp := strconv.FormatInt(currentTime, 10)
	fmt.Println(currentTimeStamp)
	fmt.Printf("%s\n", currentTimeStamp)

	message := fmt.Sprintf("%s_%s_%s_%s", ApiKey, secret, rsAddDelsginUri, currentTimeStamp)
	fmt.Println(message)

	rsAddDelSign := ComputeHmac256(message, secret)
	// 加密签名
	fmt.Println(rsAddDelSign)

	url := fmt.Sprintf("%s?key=%s&sign=%s&timestamp=%s", rsAddDelUri, ApiKey, rsAddDelSign, currentTimeStamp)
	fmt.Println(url)
	// "reason":[{"ip":"10.90.49.45","port":"3307","vip":[{"cluster":"ilvs02_GZ-YS","warehouse":"GZ-YS","product":"001069","ip":"10.88.151.136","port":"3306",
	//"protocol":"TCP","check":"tcp","connect_timeout":5,"connect_port":"3307","weight":100,"nb_get_retry":3,"lb_algo":"rr"}]}]}

	payload := strings.NewReader(`{
    "cluster":"ilvs02_GZ-YS",
    "warehouse":"GZ-YS",
    "product":"001069",
    "productname": "001069",
    "vip":[
        {
            "ip":"10.88.151.136",
            "port":"3306",
            "protocol":"TCP",
            "member":{
                "add":[
                    {
                        "ip":"10.79.23.45",
                        "port":"5306"
                    }
                ],
                "del":[
                    {
                        "ip":"10.79.23.46",
                        "port":"5306"
                    }
                ]
            }
        }
    ]
}
`)

	client := &http.Client{}
	req, err := http.NewRequest("POST", rsAddDelUri, payload)
	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
	//fmt.Printf("%T",body)

	type Response struct {
		Status string `json:"status"`
		VIP    string `json:"vip"`
		Reason string `json:"reason"`
	}

	var response Response
	err = json.Unmarshal([]byte(body), &response)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)

	if response.Status=="suc" {
		log.Println("成功")
	} else {
		log.Println("失败")
	}

}
