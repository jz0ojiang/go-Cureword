package mods

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/antlabs/strsim"
	_ "github.com/mattn/go-sqlite3"
)

// permission:
//   0. 示例权限 仅返回示例内容
//   1. 仅允许randget/get1st         每日调用次数300次
//   2. 允许randget/get1st/getword   每日调用次数500次
//   3. 允许POST upload              每日调用次数800次
//   4. 允许POST delete              不限调用次数

var PgLoc, _ = os.Executable()
var PgPath = filepath.Dir(PgLoc)
var Rerrors = [...]string{
	`{"status": "ERROR", "error_id": 50000, "error_info": "missing some data in url"}`,
	`{"status": "ERROR", "error_id": 50001, "error_info": "sign verify failed"}`,
	`{"status": "ERROR","error_id": 50002,"error_info": "value error"}`,
	`{"status": "ERROR", "error_id": 50003, "error_info": "overly similar strings"}`,
	`{"title": "go-Curewords API","interfaces": [{"method": "get","path": "/api","value": "?value=getword/get1st/randget&appid=youappid&secret=youkey","response": [{"info": "String","status": "String","words": ["string"]}]}]}`,
	`{"status": "ERROR", "error_id": 50005, "error_info": "data dont have key 'word'"}`,
	`{"status": "ERROR", "error_id": 50006, "error_info": "id does not exist(in data)"}`,
	`{"status": "ERROR", "error_id": 50007, "error_info": "id does not exist(in database)"}`,
	`{"status": "ERROR", "error_id": 50008, "error_info": "delete failed"}`,
	`{"status": "ERROR", "error_id": 50009, "error_info": "unknown error"}`,
	`{"status": "ERROR", "error_id": 50010, "error_info": "Permission Denied"}`,
	`{"status": "ERROR", "error_id": 50011, "error_info": "api usage has reached the upper limit"}`,
	`{"status": "ERROR", "error_id": 50012, "error_info": "Unknown token"}`,
	`{"status": "ERROR", "error_id": 50013, "error_info": "Invalid Connect"}`,
}

type User struct {
	Appid    string
	Secret   string
	Perm     int
	Usecount int
	Form     url.Values
}

var Database *sql.DB

func checkErr(err error) {
	if err != nil {
		Loger.Panic(err)
	}
}

func ConnectDB() {
	Database, _ = sql.Open("sqlite3", filepath.Join(PgPath, "Appdata.db"))
}

// 传入指令 输出结果（ONLY SELECT COMMAND）
func runSqlite(command string) []string {
	rows, err := Database.Query(command)
	checkErr(err)
	results := make([]string, 0)
	for rows.Next() {
		var result string
		err := rows.Scan(&result)
		checkErr(err)
		results = append(results, result)
	}
	return results
}

func CheckExist(appid string) bool {
	b, err := strconv.Atoi(runSqlite(fmt.Sprintf(`SELECT COUNT(appid) FROM token WHERE appid = '%s'`, appid))[0])
	checkErr(err)
	return (b != 0)
}

func GetUser(appid string) User {
	var per User
	per.Appid = appid
	per.Secret = runSqlite(fmt.Sprintf("SELECT secret FROM token WHERE appid = '%s'", appid))[0]
	per.Perm, _ = strconv.Atoi(runSqlite(fmt.Sprintf("SELECT permission FROM token WHERE appid = '%s'", appid))[0])
	per.Usecount, _ = strconv.Atoi(runSqlite(fmt.Sprintf("SELECT usecount FROM token WHERE appid = '%s'", appid))[0])
	return per
}

func AddCount(per User) {
	if per.Appid != "web" {
		prepare, err := Database.Prepare("UPDATE token SET usecount = ? WHERE appid = ?")
		if prepare == nil || err != nil {
			panic(err)
		}
		_, err = prepare.Exec((per.Usecount + 1), per.Appid)
		checkErr(err)
		Log(fmt.Sprintf("Appid => %s Used API", per.Appid))
	}
}

func wordIsExist(word string) bool {
	var tempwords []string = runSqlite(`SELECT Word FROM token`)
	for _, v := range tempwords {
		if strsim.Compare(word, v) >= 0.70 {
			return true
		}
	}
	return false
}

func updateWord(word string) string {
	if wordIsExist(word) {
		return Rerrors[3]
	}
	prepare, err := Database.Prepare("INSERT INTO main (Word) VALUES ('?');")
	if prepare == nil || err != nil {
		panic(err)
	}
	_, err = prepare.Exec(word)
	checkErr(err)
	Log(`Get => Updated a word\n--> ` + word)
	return fmt.Sprintf(`{
		"status": "OK",
		"info": "uploaded to the database (admin will check it later)",
		"word": "%s"
	}`, word)
}

func deleteWord(id string) string {
	iid, err := strconv.Atoi(id)
	if err != nil {
		return Rerrors[6]
	}
	idExist := func() bool {
		idc, _ := strconv.Atoi(runSqlite(fmt.Sprintf("SELECT COUNT(id) FROM main WHERE id = %d", iid))[0])
		return idc != 0
	}
	if !idExist() {
		return Rerrors[7]
	}
	stmt, err := Database.Prepare("DELETE FROM main WHERE id =?")
	checkErr(err)
	_, err = stmt.Exec(id)
	checkErr(err)
	if idExist() {
		return Rerrors[8]
	}
	return fmt.Sprintf(`{
		"status": "OK",
		"info": "Deleted data where id = %s"
	}`, id)
}

func ApiParse(per User, method string, value string) string {
	switch { // 判断permission && count
	case per.Perm == 0: // 示例操作部分
		sample := []string{"这是一条例子", "这是一条毒鸡汤(?)", "这是一条鸡汤", "这是示例中最新的句子"}
		switch value {
		case "get1st":
			return fmt.Sprintf(`{
				"status": "OK",
				"info": "get the latest word",
				"words": "%s"
			}`, sample[len(sample)-1])
		case "randget":
			rand.Seed(time.Now().UnixNano())
			return fmt.Sprintf(`{
				"status": "OK",
				"info": "random get a word",
				"words": "%s"
			}`, sample[rand.Intn(len(sample))])
		case "getword":
			return fmt.Sprintf(`{
				"status": "OK",
				"info": "get all words",
				"words": %s
			}`, list2str(sample))
		default:
			return Rerrors[2]
		}
	case per.Perm == 1:
		if per.Usecount >= 300 {
			return Rerrors[11]
		}
	case per.Perm == 2:
		if per.Usecount >= 500 {
			return Rerrors[11]
		}
	case per.Perm == 3:
		if per.Usecount >= 800 {
			return Rerrors[11]
		}
	}
	// var values map[string]string
	values := make(map[string]bool)
	values["get1st"], values["randget"], values["getword"], values["upload"], values["delete"] = true, true, true, true, true
	switch {
	case value == "get1st" && per.Perm >= 1:
		AddCount(per)
		return fmt.Sprintf(`{
			"status": "OK",
			"info": "get the latest word",
			"words": "%s"
		}`, perttier(runSqlite("SELECT Word FROM main ORDER BY id DESC LIMIT 1")[0]))
	case value == "randget" && per.Perm >= 1:
		AddCount(per)
		return fmt.Sprintf(`{
			"status": "OK",
			"info": "random get a word",
			"words": "%s"
		}`, perttier(runSqlite("SELECT Word FROM main ORDER BY RANDOM() DESC LIMIT 1")[0]))
	case value == "getword" && per.Perm >= 2:
		AddCount(per)
		return fmt.Sprintf(`{
			"status": "OK",
			"info": "get all words",
			"words": %s
		}`, list2str(runSqlite("SELECT Word FROM main")))
	case value == "upload" && method == "POST" && per.Perm >= 3:
		AddCount(per)
		if _, ok := per.Form["word"]; ok {
			return updateWord(perttier(per.Form["word"][0]))
		} else {
			return Rerrors[5]
		}
	case value == "delete" && method == "POST" && per.Perm == 4:
		if _, ok := per.Form["id"]; ok {
			return deleteWord(per.Form["id"][0])
		} else {
			return Rerrors[6]
		}
	default:
		if _, ok := values[value]; ok {
			return Rerrors[10]
		} else {
			return Rerrors[2]
		}
	}
}

// lazy as me
func list2str(input []string) string {
	var output string
	for i := 0; i < len(input); i++ {
		output += fmt.Sprintf(`"%s",`+"\n", input[i])
	}
	return "[" + output[0:len(output)-2] + "]"
}

func perttier(str string) string {
	str = strings.Replace(str, " ", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Replace(str, "\r", "", -1)
	return str
}

func Resetr() {
	pre, err := Database.Prepare("UPDATE token SET usecount=0")
	checkErr(err)
	pre.Exec()
}
