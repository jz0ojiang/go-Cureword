package mods

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/kylelemons/go-gypsy/yaml"
	"github.com/robfig/cron"
	"github.com/urfave/cli/v2"
)

var sampleWords = [...]string{
	"哪怕是示例鸡汤，也是那么的美丽",
	"就算世界毁灭，我也是一条毒鸡汤",
	"用来提供说明，也是一种值得骄傲的事情",
	"叭叭叭叭叭，我是一条鸡汤",
	"Cureword - 作者: 0o酱",
}

type configguration struct {
	host string
	port int64
}

type ResponseJson struct {
	Code int         `json:"code"`
	Info string      `json:"info"`
	Data interface{} `json:"data"`
}

type words struct {
	Word []string `json:"words"`
}

// 初始化API服务
func initAPI() {
	// 连接日志文件
	Linklog()
	// 每日重置API数量
	count := cron.New()
	count.AddFunc("0 0 4 * * *", ResetCount)
	count.Start()
}

func readSettings() configguration {
	config, err := yaml.ReadFile("config.yml")
	if err != nil {
		Log("读取config.yml时出错 将使用默认配置")
		return configguration{host: "0.0.0.0", port: 256}
	}
	host, err := config.Get("host")
	if err != nil {
		Log("读取config.yml时出错 将使用默认配置")
		return configguration{host: "0.0.0.0", port: 256}
	}
	port, err := config.GetInt("port")
	if err != nil {
		Log("读取config.yml时出错 将使用默认配置")
		return configguration{host: "0.0.0.0", port: 256}
	}
	return configguration{host: host, port: port}
}

// API核心

func userExist(appid string, secret string) bool {
	return AccountIsExist(appid) && GetAccount(appid).Secret == secret
}

func dataOperation(value string, user User) ResponseJson {
	switch value {
	case "getlast":
		if user.Perm < 1 {
			return ResponseJson{
				Code: ERROR_PERMISSION,
				Info: GetErrMsg(ERROR_PERMISSION),
				Data: struct{}{},
			}
		}
		word, err := GetNewest()
		if err != nil {
			Log("Error:" + err.Error())
			return ResponseJson{
				Code: ERROR_DATABASE,
				Info: GetErrMsg(ERROR_DATABASE),
				Data: struct{}{},
			}
		}
		var data words
		data.Word = append(data.Word, word.Word)
		return ResponseJson{
			Code: SUCCESS,
			Info: "Get newest word",
			Data: data,
		}
	case "randget":
		if user.Perm < 1 {
			return ResponseJson{
				Code: ERROR_PERMISSION,
				Info: GetErrMsg(ERROR_PERMISSION),
				Data: struct{}{},
			}
		}
		word, err := RandomGet()
		if err != nil {
			Log("Error:" + err.Error())
			return ResponseJson{
				Code: ERROR_DATABASE,
				Info: GetErrMsg(ERROR_DATABASE),
				Data: struct{}{},
			}
		}
		var data words
		data.Word = append(data.Word, word.Word)
		return ResponseJson{
			Code: SUCCESS,
			Info: "Random get a word",
			Data: data,
		}
	case "getword":
		if user.Perm < 2 {
			return ResponseJson{
				Code: ERROR_PERMISSION,
				Info: GetErrMsg(ERROR_PERMISSION),
				Data: struct{}{},
			}
		}
		word, err := GetWords()
		if err != nil {
			Log("Error:" + err.Error())
			return ResponseJson{
				Code: ERROR_DATABASE,
				Info: GetErrMsg(ERROR_DATABASE),
				Data: struct{}{},
			}
		}
		var data words
		for _, v := range word {
			data.Word = append(data.Word, v.Word)
		}
		return ResponseJson{
			Code: SUCCESS,
			Info: "Get all words",
			Data: data,
		}
	default:
		return ResponseJson{
			Code: ERROR_VALUE_ERROR,
			Info: GetErrMsg(ERROR_VALUE_ERROR),
			Data: struct{}{},
		}
	}
}

func usecountOver(user User) bool {
	return (user.Usecount >= 300 && user.Perm == 1) || (user.Usecount >= 500 && user.Perm == 2)
}

func apiTask(query url.Values) []byte {
	var result ResponseJson
	value, appid, secret := query.Get("value"), query.Get("appid"), query.Get("secret")
	if len(value+appid+secret) == 0 {
		result.Code = ERROR_MISSING_DATA
		result.Info = GetErrMsg(ERROR_MISSING_DATA)
		result.Data = struct{}{}
	} else {
		if userExist(appid, secret) {
			var user User = GetAccount(appid)
			if user.Perm == 0 {
				switch value {
				case "getlast":
					result.Code = SUCCESS
					result.Info = "Get newest word"
					var data words
					data.Word = append(data.Word, sampleWords[0])
					result.Data = data
				case "randget":
					result.Code = SUCCESS
					result.Info = "Random get a word"
					var data words
					data.Word = append(data.Word, sampleWords[rand.Intn(len(sampleWords))])
					result.Data = data
				case "getword":
					result.Code = SUCCESS
					result.Info = "Get all words"
					var data words
					data.Word = append(data.Word, sampleWords[:]...)
					result.Data = data
				default:
					result.Code = ERROR_VALUE_ERROR
					result.Info = GetErrMsg(ERROR_VALUE_ERROR)
					result.Data = struct{}{}
				}
			} else {
				if usecountOver(user) && user.Appid != "web" {
					result.Code = ERROR_USECOUNT
					result.Info = GetErrMsg(ERROR_USECOUNT)
					result.Data = struct{}{}
				} else {
					AddCount(appid)
					result = dataOperation(value, user)
				}
			}
		} else {
			result.Code = ERROR_VERIFY_FAIL
			result.Info = GetErrMsg(ERROR_VERIFY_FAIL)
			result.Data = struct{}{}
		}
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

func api(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Methods", "GET")
	writer.Header().Set("content-type", "application/json")
	query := request.URL.Query()
	writer.Write(apiTask(query))
}

// 管理面板
func admin(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Methods", "GET")
	http.ServeFile(writer, request, "./html/admin.html")
}

// 运行API服务
func Run(c *cli.Context) error {
	initAPI()
	defer LogFile.Close()
	settings := readSettings()
	fmt.Println("Cureword-API serve will start in 5 seconds (press Ctrl+C to Cancel)")
	time.Sleep(5 * time.Second)
	host := fmt.Sprintf("%s:%d", settings.host, settings.port)
	Log(fmt.Sprintf("Cureword-API will run in %s:%d\n", settings.host, settings.port))

	// 绑定函数
	http.Handle("/", http.FileServer(http.Dir("./doc")))
	http.HandleFunc("/api", api)
	http.HandleFunc("/app", App)
	http.HandleFunc("/admin", admin)
	http.HandleFunc("/submit", Submit)

	err := http.ListenAndServe(host, nil)
	return err
}
