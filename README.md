# Go-Cureword

![Cureword](https://img.shields.io/badge/Cureword-fe8000?style=flat&logo=Go) ![LICENSE](https://img.shields.io/badge/License-AGPL--3.0_License-yellow?style=flat) ![version](https://img.shields.io/badge/Version-0.1.1_R-blueviolet?style=flat) ![visit time](https://visitor-badge.glitch.me/badge?page_id=github-com-0ojixueseno0-go_cureword)

## What do I use

- Go

  cron

  go-sqlite3

  strsim


- Sqlite3


## Sqlite

```
File ::(Appdata.db)
|| - Table
|| -> main
||    || -> id     integer
||    || -> Word   text
||
|| -> token
||    || -> id          integer
||    || -> appid       text
||    || -> secret      text
||    || -> permission  integer
||    || -> usecount    integer
```

## API

```
Getword - Get All words in database
    useage: [GET/POST] host/api?value=getword&appid=#appid&secret=#secret
    response body:
      {
        "status": "OK",
        "info": "get all words",
        "words": [
            "这是一条例子",
            "这是一条毒鸡汤(?)",
            "这是一条鸡汤",
            "这是示例中最新的句子"
        ]
      }

Get1st - Get the latest word
    useage: [GET/POST] host/api?value=get1st&appid=#appid&secret=#secret
    response body:
      {
        "status": "OK",
        "info": "get the latest word",
        "words": "这是示例中最新的句子"
      }

RandomGet - random get a word
    useage: [GET/POST] host/api?value=randget&appid=#appid&secret=#secret
    response body:
      {
        "status": "OK",
        "info": "random get a word",
        "words": "这是一条毒鸡汤(?)"
      }
```

---

```
Upload - upload a word to database
    useage: [POST] host/api?value=upload&appid=#appid&secret=#secret
    request header:
      Content-Type => application/x-www-form-urlencoded
    request body:
      word:#Cureword
    response body:
      {
        "status": "OK",
        "info": "uploaded to the database (admin will check it later)",
        "word": "#Cureword"
      }

Delete - delete a word with id
    useage: [POST] host/api?value=delete&appid=#appid&secret=#secret
    request header:
      Content-Type => application/x-www-form-urlencoded
    request body:
      id:#uid
    response body:
      {
        "status": "OK",
        "info": "Deleted data where id = #uid"
	  }

```

### Permission

> 0. Sample permissions only return sample content.
>
> 1. Only 300 randget/get1st calls per day are allowed.
>
> 2. Allow randget/get1st/getword calls 500 times per day.
>
> 3. Allow POST upload calls 800 times per day.
>
> 4. Allow POST delete unlimited number of calls

### ErrorCode

```

  50000 - missing some data in url
  50001 - sign verify failed
  50002 - value error
  50003 - overly similar strings
  50005 - data dont have key 'word'
  50006 - id does not exist(in data)
  50007 - id does not exist(in database)
  50008 - delete failed
  50009 - unknown error
  50010 - Permission Denied
  50011 - api usage has reached the upper limit
  50012 - Unknown token
  50013 - Invalid Connect

```
