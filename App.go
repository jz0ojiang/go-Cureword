package main

import (
	"flag"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/0ojixueseno0/go-cureword/mods"
	"github.com/robfig/cron"
)

// 定义flags
var Host = flag.String("host", "localhost:256", "Program run in host:port")
var Fastmode = flag.Bool("faststart", false, "是否快速启动")

func main() {
	flag.Parse() // 解析host, fastmode
	if !*Fastmode {
		fmt.Println("Cureword-API serve will start in 5 seconds (press Ctrl+C to Cancel)")
		time.Sleep(5 * time.Second)
	}
	// 连接日志文件
	mods.Linklog()
	defer mods.LogFile.Close()

	mods.Log("Serve Succeed run in " + *Host)

	// 数据库连接
	mods.ConnectDB()
	defer mods.Database.Close()

	// 每日重置API数量
	c := cron.New()
	c.AddFunc("0 0 4 * * *", mods.Resetr)
	c.Start()

	// 绑定函数
	http.HandleFunc("/", IndexHandler)
	http.HandleFunc("/api", CWAPI)
	http.HandleFunc("/test", posttest)
	err := http.ListenAndServe(*Host, nil)
	if err != nil {
		panic(err)
	}
	// time.Sleep(10 * time.Second)
}

func posttest(writer http.ResponseWriter, request *http.Request) {
	request.ParseForm()
	fmt.Println(request.Form)
	fmt.Println(reflect.TypeOf(request.Form["id"][0]))
	writer.Write([]byte(string(request.PostFormValue("word"))))
}

func IndexHandler(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	writer.Header().Set("content-type", "application/json")
	fmt.Fprintln(writer, mods.Rerrors[4])
}
func CWAPI(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	writer.Header().Set("content-type", "application/json")
	query := request.URL.Query()
	value, appid, secret := query.Get("value"), query.Get("appid"), query.Get("secret")
	if value == "" || appid == "" || secret == "" { //判断是否传入了需要的值
		fmt.Fprintln(writer, mods.Rerrors[0])
	} else {
		if !mods.CheckExist(appid) { // 判断Appid是否存在
			fmt.Fprintln(writer, mods.Rerrors[1])
		} else {
			// 解析传入数据 存入结构体变量
			request.ParseForm()
			var per mods.User = mods.GetUser(appid)
			per.Form = request.Form

			if secret == per.Secret { // 判断密钥是否正确
				fmt.Fprintln(writer, mods.ApiParse(per, request.Method, value))
			} else {
				fmt.Fprintln(writer, mods.Rerrors[1])
			}
		}
	}
}
