package mods

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/antlabs/strsim"
)

//TODO: blacklist

type submitJson struct {
	Content string `json:"message"`
	Contact string `json:"email"`
}

type responseJson struct {
	Code int    `json:"code"`
	Info string `json:"info"`
}

//检查文字是否与已有文字过度相似
func wordIsExist(word string) bool {
	tempwords, _ := GetTempWords()
	for _, v := range tempwords {
		if strsim.Compare(word, v.Content) >= 0.70 {
			return true
		}
	}
	return false
}

func submitTask(request *http.Request) []byte {
	var result responseJson
	defer request.Body.Close()
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		Log("Error:" + err.Error())
		result.Code = ERROR_POSTBODY
		result.Info = GetErrMsg(ERROR_POSTBODY) // 解析post.body出错
	} else {
		var data submitJson
		if err := json.Unmarshal(body, &data); err == nil {
			if !wordIsExist(data.Content) {
				if ListCountOfTempwords() <= 1000 {
					if err := AddNewTempWord(data.Content, data.Contact); err == nil {
						Log("[Submit] " + data.Content + "From: " + data.Contact + "-IP: " + request.RemoteAddr)
						result.Code = SUCCESS
						result.Info = "success uploaded"
					} else {
						result.Code = ERROR_DATABASE
						result.Info = GetErrMsg(ERROR_DATABASE) // 执行 addnewtempword 出错
					}
				} else {
					result.Code = ERROR_TEMPWORDSFULL
					result.Info = GetErrMsg(ERROR_TEMPWORDSFULL) // tempword超过1000条
				}
			} else {
				result.Code = ERROR_SIMILARWORD
				result.Info = GetErrMsg(ERROR_SIMILARWORD) // 和已有的过度相似
			}
		} else {
			Log("Error:" + err.Error())
			result.Code = ERROR_POSTBODY
			result.Info = GetErrMsg(ERROR_POSTBODY) // 解析json出错
		}
	}
	bytes, _ := json.Marshal(result)
	return bytes
}

func Submit(writer http.ResponseWriter, request *http.Request) {
	// writer.Header().Set("Access-Control-Allow-Origin", "*") //允许访问所有域
	// writer.Header().Set("Access-Control-Allow-Headers", "*")
	writer.Header().Set("Access-Control-Allow-Methods", "POST")
	writer.Header().Set("content-type", "application/json")

	writer.Write(submitTask(request))
}
