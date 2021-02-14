package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"srvguide/global"
	"srvguide/internal/model"
	"srvguide/router"
	"srvguide/util"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/go-irain/logger"

	"github.com/Unknwon/goconfig"
)

var config *goconfig.ConfigFile

func init() {
	var err error

	config, err = goconfig.LoadConfigFile("config/config.ini")
	if err != nil {
		panic(err.Error())
	}

	global.Config = config
	fmt.Println("init configuration with release mode;")

	// 判断log文件夹，不存在则创建
	tempDir := config.MustValue("log", "path")
	var isExistFlag bool
	isExistFlag, err = util.PathExists(tempDir)
	if err != nil {
		panic(err.Error())
	}
	if false == isExistFlag {
		err := os.MkdirAll(tempDir, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	}

	// 设置日志前台可见
	consoleFlag := false
	if config.MustInt("log", "console") == 1 {
		consoleFlag = true
	}
	log.SetConsole(consoleFlag)

	// 根据配置文件设置日志等级
	logLevel := log.ERROR
	switch config.MustValue("log", "level") {
	case "Debug":
		logLevel = log.DEBUG
	case "Info":
		logLevel = log.INFO
	case "Warn":
		logLevel = log.WARN
	case "Error":
		logLevel = log.ERROR
	case "Fatal":
		logLevel = log.FATAL
	default:
	}
	log.SetLevel(logLevel)

	// 根据配置文件，设置日志路径，日志名，日志切割大小限制
	log.SetRollingFile(config.MustValue("log", "path"),
		config.MustValue("log", "filename"),
		int32(config.MustInt("log", "num", 10)),
		int64(config.MustInt("log", "max", 50)), log.MB)
	log.JSON(false)

	// 获取当前工作目录
	log.Debug(os.Getwd())
}

// main ...
func main() {
	// database, err := db.CreateDatabase()
	// if err != nil {
	// 	log.Fatal("Database connection failed: %s", err.Error())
	// }
	model.InitDB(config)

	// app := &app.App{
	// 	Router: mux.NewRouter().StrictSlash(true),
	// 	// Database: database,
	// }
	// app.SetupRouter()

	// log.Fatal(http.ListenAndServe(":8080", app.Router))

	//初始化gin服务
	svr := gin.Default()

	//加载预处理中间件
	svr.Use(requestMiddleHandler)

	//加载响应处理中间件
	svr.Use(responseMiddleHandler)

	//初始化路由
	router.Init(svr)

	//启动服务
	err := svr.Run(config.MustValue("server", "ip") + ":" + config.MustValue("server", "port"))
	if err != nil {
		panic(err.Error())
	}
}

// 预处理中间件
func requestMiddleHandler(c *gin.Context) {
	c.Set("start_time", time.Now())

	//生成logId
	logid := util.GenerateLogid()
	err := c.Request.ParseForm()
	if err != nil {
		panic(err.Error())
	}
	log.Debug(logid, "**************** Client : ",
		c.Request.RemoteAddr, " [", c.Request.Method, "] ", c.Request.URL.Path, " ****************")
	log.Debug(logid, "request : [", c.Request.URL, "] ", c.Request.PostForm)

	//设置logId
	c.Set("logid", logid)
	c.Next()
}

// 响应客户端中间件
func responseMiddleHandler(c *gin.Context) {
	c.Next()
	logid, _ := c.Get("logid")
	err, exist := c.Get("err")
	//如果发生错误，则返回错误信息
	if exist {
		errMap := map[string]interface{}{"status": 1, "message": err.(string)}
		c.JSON(http.StatusOK, errMap)
		resByte, err := json.Marshal(errMap)
		if err != nil {
			log.Error("json marshal failed")
		}
		log.Debug(logid, "Response : ", string(resByte))
	} else { //否则返回正确结果
		res, exist := c.Get("res")
		if exist {
			c.JSON(http.StatusOK, res)
			resByte, err := json.Marshal(res)
			if err != nil {
				log.Error("json marshal failed")
			}
			log.Debug(logid, "Response : ", string(resByte))
		}
	}
	startTime, _ := c.Get("start_time")
	timeEnd := time.Now()
	duration := timeEnd.Sub(startTime.(time.Time))
	log.Debug(logid, "**************** Cost Duration : ", duration.String(), " ****************\r\n")
}
