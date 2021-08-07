# Go-Cureword

![Cureword](https://img.shields.io/badge/Cureword-a2b8a5?style=flat&logo=Go) ![LICENSE](https://img.shields.io/badge/License-AGPL--3.0_License-yellow?style=flat) ![version](https://img.shields.io/badge/Version-0.2.0_R-blueviolet?style=flat) ![visit time](https://visitor-badge.glitch.me/badge?page_id=github-com-0ojixueseno0-go_cureword) [![wakatime](https://wakatime.com/badge/github/0ojixueseno0/go-Cureword.svg)](https://wakatime.com/badge/github/0ojixueseno0/go-Cureword)

## What do I use

![Data From WakaTime](https://gitee.com/im0o/photobed/raw/master/img/20210807153739.png)


## Sqlite

![](https://gitee.com/im0o/photobed/raw/master/img/20210807153209.png)

## API

API Doc: http://api.cureword.top:256

## Cli

Commands:

```bash
$ go-cureword

COMMANDS:
   account  API account operations
   run      Run API serve
   help, h  Shows a list of commands or help for one command

account:
  COMMANDS:
   add      add a new account
   list     list all accounts
   delete   delete an account
   set      set an account
   help, h  Shows a list of commands or help for one command

```

## Config

```yaml
host: 0.0.0.0
port: 256
```

## Permission

> 0. Sample permissions only return sample content.
>
> 1. Only 300 randget/get1st calls per day are allowed.
>
> 2. Allow randget/get1st/getword calls 500 times per day.
>
> 3. Can use admin dashboard(with google_auth)

## Errorno

此处为反馈中所有可能出现的类型与错误代码

|        类型        | 错误代码 |             描述              |
| :----------------: | :------: | :---------------------------: |
|      SUCCESS       |   200    |           正常获取            |
| ERROR_MISSING_DATA |   1001   |  缺失请求数据（检查你的url）  |
| ERROR_VERIFY_FAIL  |   1002   |   验证失败，请检查你的密钥    |
| ERROR_VALUE_ERROR  |   1003   | 请求方法异常，检查value内的值 |
|  ERROR_PERMISSION  |   1004   |           权限不足            |
|   ERROR_USECOUNT   |   1005   |     当日请求次数已达上限      |
|   ERROR_UNKNOWN    |   1006   |           未知错误            |
|   ERROR_DATABASE   |   2001   |        数据库连接出错         |


