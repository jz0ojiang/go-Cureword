package mods

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	Appid      string
	Secret     string
	Perm       int
	Usecount   int
	GoogleAuth string
}

func (User) TableName() string {
	return "users"
}

type Word struct {
	Id   int
	Word string
}

func (Word) TableName() string {
	return "words"
}

type TempWord struct {
	Id      int    `json:"id"`      // id
	Content string `json:"content"` // 内容
	Contact string `json:"contact"` // 联系方式
	Time    string `json:"time"`    // 添加时间
}

func (TempWord) TableName() string {
	return "tempwords"
}

var DB *gorm.DB = Connect()

// 连接数据库
func Connect() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("Appdata.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

// 获取待审核的鸡汤总数
func ListCountOfTempwords() int64 {
	result := DB.Find(&TempWord{})
	return result.RowsAffected
}

// 添加一条待审核的鸡汤
func AddNewTempWord(word string, contact string) error {
	var tword TempWord
	tword.Content = word
	tword.Contact = contact
	tword.Time = time.Now().Format("2006-01-02 15:04:05")
	result := DB.Create(&tword)
	return result.Error
}

// 通过待审核的鸡汤
func AcceptTempWords(twords []TempWord) error {
	for _, tword := range twords {
		var word Word
		word.Word = tword.Content
		result := DB.Create(&word)
		if result.Error != nil {
			return result.Error
		}
	}
	return DeleteTempWords(twords)
}

// 删除待审核的鸡汤
func DeleteTempWords(twords []TempWord) error {
	for _, tword := range twords {
		result := DB.Delete(&tword)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// 从tempwords获取内容
func GetTempWords() ([]TempWord, error) {
	var words []TempWord
	result := DB.Find(&words)
	return words, result.Error
}

// 添加使用次数
func AddCount(appid string) error {
	if AccountIsExist(appid) {
		user := GetAccount(appid)
		user.Usecount += 1
		result := DB.Model(&user).Where("appid = ?", user.Appid).Updates(user)
		return result.Error
	}
	return nil
}

// 检测账户是否存在
func AccountIsExist(appid string) bool {
	var user User
	result := DB.First(&user, "appid = ?", appid)
	return result.RowsAffected != 0
}

// 获取用户信息
func GetAccount(appid string) User {
	var user User
	DB.First(&user, "appid = ?", appid)
	return user
}

//* AutoRun
func ResetCount() {
	DB.Model(User{}).Updates(User{Usecount: 0})
}

/**
* @api {Get} /api 获取所有鸡汤
* @apiGroup GetAllWords
* @apiDescription 获取所有鸡汤 以[]string形式返回
*
* @apiParam {String} value 参数(getlast|randget|getword)
* @apiParam {String} appid 请求ID
* @apiParam {String} secret 请求密钥
* @apiParamExample {json} 参数示例:
* {
*	  "value": "getword",
*	  "appid": "test",
*	  "secret": "test"
* }
*
* @apiError (FAIL) {Number} code 错误码
* @apiError (FAIL) {String} info 错误信息
* @apiErrorExample 错误响应示例
* HTTP/1.1 200 OK
* {
*	  "code": 1002,
*	  "info": "验证失败",
*	  "data": {}
* }
*
* @apiSuccess (Success) {Number} code 状态码
* @apiSuccess (Success) {String} info 状态码
* @apiSuccess (Success) {Json} data 数据
* @apiSuccess (Success) {[]String} data.words 内容
* @apiSuccessExample 成功响应示例
* HTTP/1.1 200 OK
* {
* 	"code": 200,
* 	"info": "Get all words",
* 	"data": {
* 		"words": [
* 			"哪怕是示例鸡汤，也是那么的美丽",
* 			"就算世界毁灭，我也是一条毒鸡汤",
* 			"用来提供说明，也是一种值得骄傲的事情",
* 			"叭叭叭叭叭，我是一条鸡汤",
* 			"Cureword - 作者: 0o酱"
* 		]
* 	}
* }
 */
func GetWords() ([]Word, error) {
	var words []Word
	result := DB.Find(&words)
	return words, result.Error
}

/**
* @api {Get} /api 随机获取一条鸡汤
* @apiGroup RandomGetWord
* @apiDescription 随机获取一条鸡汤 以[]string形式返回
*
* @apiParam {String} value 参数(getlast|randget|getword)
* @apiParam {String} appid 请求ID
* @apiParam {String} secret 请求密钥
* @apiParamExample {json} 参数示例:
* {
*	  "value": "randget",
*	  "appid": "test",
*	  "secret": "test"
* }
*
* @apiError (FAIL) {Number} code 错误码
* @apiError (FAIL) {String} info 错误信息
* @apiErrorExample 错误响应示例
* HTTP/1.1 200 OK
* {
*	  "code": 1002,
*	  "info": "验证失败",
*	  "data": {}
* }
*
* @apiSuccess (Success) {Number} code 状态码
* @apiSuccess (Success) {String} info 状态码
* @apiSuccess (Success) {Json} data 数据
* @apiSuccess (Success) {[]String} data.words 内容
* @apiSuccessExample 成功响应示例
* HTTP/1.1 200 OK
* {
* 	"code": 200,
* 	"info": "Random get a word",
* 	"data": {
* 		"words": [
* 			"用来提供说明，也是一种值得骄傲的事情"
* 		]
* 	}
* }
 */
func RandomGet() (Word, error) {
	var word Word
	result := DB.Order("random() desc").Find(&word)
	return word, result.Error
}

/**
* @api {Get} /api 获取最新的一条鸡汤
* @apiGroup GetNewestWord
* @apiDescription 获取最新的一条鸡汤 以[]string形式返回
*
* @apiParam {String} value 参数(getlast|randget|getword)
* @apiParam {String} appid 请求ID
* @apiParam {String} secret 请求密钥
* @apiParamExample {json} 参数示例:
* {
*	  "value": "randget",
*	  "appid": "test",
*	  "secret": "test"
* }
*
* @apiError (FAIL) {Number} code 错误码
* @apiError (FAIL) {String} info 错误信息
* @apiErrorExample 错误响应示例
* HTTP/1.1 200 OK
* {
*	  "code": 1002,
*	  "info": "验证失败",
*	  "data": {}
* }
*
* @apiSuccess (Success) {Number} code 状态码
* @apiSuccess (Success) {String} info 状态码
* @apiSuccess (Success) {Json} data 数据
* @apiSuccess (Success) {[]String} data.words 内容
* @apiSuccessExample 成功响应示例
* HTTP/1.1 200 OK
* {
* 	"code": 200,
* 	"info": "Get newest word",
* 	"data": {
* 		"words": [
* 			"哪怕是示例鸡汤，也是那么的美丽"
* 		]
* 	}
* }
 */
func GetNewest() (Word, error) {
	var word Word
	result := DB.Last(&word)
	return word, result.Error
}
