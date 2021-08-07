package cli

import (
	"bytes"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"strconv"

	"github.com/0ojixueseno0/go-cureword/mods"
	"github.com/go-basic/uuid"
	"github.com/urfave/cli/v2"
)

// permission:
//   0. 示例权限 仅返回示例内容
//   1. 仅允许randget/get1st         每日调用次数300次
//   2. 允许randget/get1st/getword   每日调用次数500次
//   3. 带有GoogleAuth的管理员账号

func input(prompt string) string {
	var text string
	fmt.Print(prompt)
	fmt.Scanln(&text)
	return text
}

func createRandomString(len int) string {
	var container string
	var str = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
	b := bytes.NewBufferString(str)
	length := b.Len()
	bigInt := big.NewInt(int64(length))
	for i := 0; i < len; i++ {
		randomInt, _ := rand.Int(rand.Reader, bigInt)
		container += string(str[randomInt.Int64()])
	}
	return container
}

func generateGoogleAuthQR(appid string) string {
	fmt.Println("请打开链接使用GoogleAuth扫描二维码导入动态令牌")
	google_auth := mods.NewGoogleAuth().GetSecret()
	fmt.Println(mods.NewGoogleAuth().GetQrcodeUrl("[Cureword-cli]"+appid, google_auth))
	return google_auth
}

func AddAccount(c *cli.Context) error {
	fmt.Println("[Cureword-cli] 进入api账号创建向导")
	var user mods.User
	user.Appid = input("\n- 账号(appid[string])(留空自动生成):")
	if user.Appid == "" {
		user.Appid = uuid.New()
	}
	user.Secret = input("\n- 密码(secret[string])(留空自动生成):")
	if user.Secret == "" {
		user.Secret = createRandomString(10)
	}
	perm := input("\n- 权限[int(0-3)](default: 1):")
	if perm == "" {
		perm = "1"
	}
	var err error
	user.Perm, err = strconv.Atoi(perm)
	if 0 > user.Perm || user.Perm > 3 {
		return errors.New("权限不合法")
	}
	if err != nil {
		return err
	}
	if user.Perm == 3 {
		user.GoogleAuth = generateGoogleAuthQR(user.Appid)

	}
	if mods.AccountIsExist(user.Appid) {
		return errors.New("该Appid已经存在 请重新键入一个新的Appid")
	}
	result := mods.DB.Create(&user)
	// fmt.Println(result.Error == nil)
	return result.Error
}

func ListAccount(c *cli.Context) error {
	var users []mods.User
	result := mods.DB.Find(&users)
	if result.Error != nil {
		return result.Error
	}
	fmt.Printf("CurewordAPI Account List (Total: %d)\n[序号] Appid - Secret - 权限 - 使用次数\n", result.RowsAffected)
	for i, user := range users {
		fmt.Printf("[%d] %s - %s - %d - %d\n", i, user.Appid, user.Secret, user.Perm, user.Usecount)
	}
	return nil
}

func DeleteAccount(c *cli.Context) error {
	var user mods.User
	if c.Args().Len() != 1 {
		fmt.Println("[Cureword-cli] Usage: <command> delete appid(string)\nTips: You can use <command> list to list all accounts")
		return nil
	} else {
		user.Appid = c.Args().Get(0)
		if mods.AccountIsExist(user.Appid) {
			if input("输入('ConfirmDelete')确认删除账号:") == "ConfirmDelete" {
				fmt.Println("account deleted")
				mods.DB.Where("appid = ?", user.Appid).Delete(&user)
				return nil
			} else {
				return errors.New("account does not exist")
			}
		}
	}
	return errors.New("Could not found the account Named " + user.Appid)
}

func SetAccount(c *cli.Context) error {
	if c.Args().Len() != 1 {
		fmt.Println("[Cureword-cli] Usage: <command> set <appid(string)>\nTips: You can use <command> list to list all accounts")
		return nil
	} else {
		appid := c.Args().Get(0)
		if mods.AccountIsExist(appid) {
			var user mods.User
			var original mods.User = mods.GetAccount(appid)
			fmt.Println("[Cureword-cli] 进入api账号设置向导(留空不进行更改)")
			user.Appid = input(fmt.Sprintf("\n- Appid(%s):", original.Appid))
			if user.Appid == "" {
				user.Appid = original.Appid
			}
			user.Secret = input(fmt.Sprintf("\n- Secret(%s):", original.Secret))
			if user.Secret == "" {
				user.Secret = original.Secret
			}
			perm := input(fmt.Sprintf("\n- Perm(%d):", original.Perm))
			if perm == "" {
				user.Perm = original.Perm
			} else {
				var err error
				user.Perm, err = strconv.Atoi(perm)
				if err != nil {
					return err
				}
				if 0 > user.Perm || user.Perm > 3 {
					return errors.New("权限不合法")
				}
			}
			usecount := input(fmt.Sprintf("\n- Usecount(%d):", original.Usecount))
			if usecount == "" {
				user.Usecount = original.Usecount
			} else {
				var err error
				user.Usecount, err = strconv.Atoi(usecount)
				if err != nil {
					return err
				}
			}
			if user.Perm == 3 && original.Perm == 3 {
				user.GoogleAuth = input("\n- 是否重新生成GoogleAuth二维码?(y/N)")
				if user.GoogleAuth == "Y" || user.GoogleAuth == "y" {
					user.GoogleAuth = generateGoogleAuthQR(user.Appid)
				}
			} else if user.Perm == 3 && original.Perm != 3 {
				user.GoogleAuth = generateGoogleAuthQR(user.Appid)
			} else {
				mods.DB.Model(&user).Where("appid = ?", original.Appid).Select("google_auth").Updates(mods.User{GoogleAuth: ""})
			}
			if user.Perm == 0 {
				mods.DB.Model(&user).Where("appid = ?", original.Appid).Select("perm").Updates(mods.User{Perm: 0})
			}
			mods.DB.Model(&user).Where("appid = ?", original.Appid).Updates(user)
			fmt.Println(user)
		} else {
			return errors.New("account does not exist")
		}
	}
	return nil
}
