package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Traces struct {
	AcceptTime string `json:"AcceptTime"`
	AcceptStation string `json:"AcceptStation"`
	Remark string `json:"Remark"`
}

type K struct {
	EBusinessID string   `json:"EBusinessID"`
	Traces []Traces      `json:"Traces"`
	OrderCode string     `json:"OrderCode"`
	ShipperCode string   `json:"ShipperCode"`
	LogisticCode string  `json:"LogisticCode"`
	Success bool         `json:"Success"`
	Reason string        `json:"Reason"`
	State string         `json:"State"`
}

const  (
	EBusinessID = ""
	RequestType = "1002"
	DataType = "2-json"
	RequestUrl = "http://sandboxapi.kdniao.com:8080/kdniaosandbox/gateway/exterfaceInvoke.json"
	AppKEY = ""
	ShipSfCode = "SF"
)

func main()  {
    var OrderCode,LogisticCode string
	OrderCode = ""
	LogisticCode = ""
	niao,_:= kuaiDiNiao(OrderCode,ShipSfCode,LogisticCode)
    fmt.Println(niao)
}

func kuaiDiNiao(OrderCode,ShipperCode,LogisticCode string) (k K,err error) {
	reqMap := make(map[string]string)

	reqMap["OrderCode"] = OrderCode
	reqMap["ShipperCode"] = ShipperCode
	reqMap["LogisticCode"] = LogisticCode

	requestContent,_ := json.Marshal(reqMap)

	md5Cont := md5.New()
	md5Cont.Write([]byte(string(requestContent) + AppKEY))
	toString := hex.EncodeToString(md5Cont.Sum(nil))

	RequestData := url.QueryEscape(string(requestContent))

	DataSign := url.QueryEscape(base64.StdEncoding.EncodeToString([]byte(toString)))

	data := url.Values{
		"RequestData" : {RequestData},
		"EBusinessID" : {EBusinessID},
		"RequestType" : {RequestType},
		"DataType" : {DataType},
		"DataSign" : {DataSign},
	}

	resp, err := http.PostForm(RequestUrl, data)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	_ = json.Unmarshal(body, &k)
	if !k.Success {
		print(k.Reason)
		return k,errors.New(k.Reason)
	}
	return k,nil

}
