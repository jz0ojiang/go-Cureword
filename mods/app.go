package mods

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

type AppResponseJson struct {
	Code int        `json:"code"`
	Info string     `json:"info"`
	Data []TempWord `json:"data"`
}

type appJson struct {
	Account string     `json:"account"`
	Token   string     `json:"accessToken"`
	Selects []TempWord `json:"selects"`
	Type    string     `json:"type"`
}

var PgLoc, _ = os.Executable()
var PgPath = filepath.Dir(PgLoc)

func appTask(request *http.Request) []byte {
	var result AppResponseJson
	var query url.Values = request.URL.Query()
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		Log("Error:" + err.Error())
		result.Code = ERROR_POSTBODY
		result.Info = GetErrMsg(ERROR_POSTBODY)
		result.Data = nil
	} else {
		switch query.Get("value") {
		case "submit": // 提交事件
			var data appJson
			if err := json.Unmarshal(body, &data); err == nil {
				if AccountIsExist(data.Account) {
					user := GetAccount(data.Account)
					gauth, err := NewGoogleAuth().GetCode(user.GoogleAuth)
					if gauth == data.Token && err == nil {
						switch data.Type {
						case "accept":
							if err := AcceptTempWords(data.Selects); err != nil {
								result.Code = ERROR_DATABASE
								result.Info = GetErrMsg(ERROR_DATABASE)
							} else {
								result.Code = SUCCESS
								result.Info = "Success accept words"
							}
						case "refuse":
							if err := DeleteTempWords(data.Selects); err != nil {
								result.Code = ERROR_DATABASE
								result.Info = GetErrMsg(ERROR_DATABASE)
							} else {
								result.Code = SUCCESS
								result.Info = "Success delete words"
							}
						default:
							result.Code = ERROR_VALUE_ERROR
							result.Info = GetErrMsg(ERROR_VALUE_ERROR)
							result.Data = nil
						}
					} else {
						result.Code = ERROR_VERIFY_FAIL
						result.Info = GetErrMsg(ERROR_VERIFY_FAIL)
						result.Data = nil
					}
				} else {
					result.Code = ERROR_VERIFY_FAIL
					result.Info = GetErrMsg(ERROR_VERIFY_FAIL)
					result.Data = nil
				}
			} else {
				Log("Error:" + err.Error())
				result.Code = ERROR_POSTBODY
				result.Info = GetErrMsg(ERROR_POSTBODY) // 解析json出错
				result.Data = nil
			}
		case "get": // 获取tempwords
			if tword, err := GetTempWords(); err != nil {
				result.Code = ERROR_DATABASE
				result.Info = GetErrMsg(ERROR_DATABASE)
				result.Data = nil
			} else {
				result.Code = SUCCESS
				result.Info = "Get TempWords"
				result.Data = tword
			}
		default:
			result.Code = ERROR_VALUE_ERROR
			result.Info = GetErrMsg(ERROR_VALUE_ERROR)
			result.Data = nil
		}
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

// route:admin 接口
func App(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set("Access-Control-Allow-Origin", "http://www.cureword.top") //允许访问所有域
	writer.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	writer.Header().Set("content-type", "application/json")
	writer.Write(appTask(request))
}
