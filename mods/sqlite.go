package mods

import (
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

type Word struct {
	Id   int
	Word string
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

// func checkErr(err error) {
// 	if err != nil {
// 		Loger.Panic(err)
// 	}
// }

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

// 获取所有鸡汤
func GetWords() ([]Word, error) {
	var words []Word
	result := DB.Find(&words)
	return words, result.Error
}

// 随机获取一条鸡汤
func RandomGet() (Word, error) {
	var word Word
	result := DB.Order("random() desc").Find(&word)
	return word, result.Error
}

// 获取最新的一条鸡汤
func GetNewest() (Word, error) {
	var word Word
	result := DB.Last(&word)
	return word, result.Error
}

// 检查文字是否与已有文字过度相似
// func wordIsExist(word string) bool {
// 	tempwords, _ := GetWords()
// 	for _, v := range tempwords {
// 		if strsim.Compare(word, v.Word) >= 0.70 {
// 			return true
// 		}
// 	}
// 	return false
// }
