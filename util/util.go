package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"srvguide/global"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

//GenerateLogid 生成logid
func GenerateLogid() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return strconv.FormatInt(time.Now().UnixNano(), 10) + fmt.Sprintf(":%d", r.Intn(10000))
}

//String interface 转string
func String(v interface{}) string {
	return fmt.Sprintf("%v", v)
}

//Int64 字符串转int64
func Int64(s string) int64 {
	if s == "" {
		return int64(0)
	}
	n, err := strconv.ParseInt(s, 10, 0)
	if err != nil {
		return 0
	}
	return n
}

//Int 字符转整形
func Int(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return n
}

//GetRandomString 获取随机字符串
func GetRandomString(num int, str ...string) string {
	s := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if len(str) > 0 {
		s = str[0]
	}
	l := len(s)
	r := rand.New(rand.NewSource(getRandSeek()))
	var buf bytes.Buffer
	for i := 0; i < num; i++ {
		x := r.Intn(l)
		buf.WriteString(s[x : x+1])
	}
	return buf.String()
}

var (
	randSeek = int64(1)
	l        sync.Mutex
)

func getRandSeek() int64 {
	l.Lock()
	if randSeek >= 100000000 {
		randSeek = 1
	}
	randSeek++
	l.Unlock()
	return time.Now().UnixNano() + randSeek
}

//Logid 获取logid
func Logid(c *gin.Context) string {
	logid, ok := c.Get("logid")
	if !ok {
		return ""
	}

	return logid.(string)
}

type structInData struct {
	Code int `json:"code"`
	Data struct {
		InTime     int    `json:"in_time"`
		ParkCode   string `json:"park_code"`
		VplNumber  string `json:"vpl_number"`
		RecordID   int64  `json:"record_id"`
		InArmCode  string `json:"in_arm_code"`
		TraceUk    int    `json:"trace_uk"`
		DepotID    int    `json:"depot_id"`
		VplType    int    `json:"vpl_type"`
		VplColor   int    `json:"vpl_color"`
		IsNoPlate  bool   `json:"is_no_plate"`
		UpdateTime int    `json:"update_time"`
		Info       struct {
			Color int    `json:"color"`
			Logo  string `json:"logo"`
			Model int    `json:"model"`
			Power int    `json:"power"`
			Type  int    `json:"type"`
		} `json:"info"`
		CarType int `json:"car_type"`
	} `json:"data"`
	Message string `json:"message"`
}

func GetSiteData(vpl string, parkCode string) string {
	client := &http.Client{}
	resp, err := client.Get(global.Config.MustValue("third_url", "url_get_in_data") + "?vpl=" + vpl + "&park=" + parkCode)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var inData structInData
	err = json.Unmarshal(body, &inData)
	if err != nil {
		return ""
	}

	return inData.Data.InArmCode
}
