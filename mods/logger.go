package mods

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

var Loger *log.Logger
var LogFile *os.File

func Linklog() {
	path := filepath.Join(PgPath, "logs/Curewords-"+string(time.Now().Format("2006-01-02"))+".log")
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			_ = os.Mkdir("./logs", 0777)
			file, err = os.Create(path)
			if err != nil {
				panic(err)
			}
		} else {
			panic(err)
		}
	}
	LogFile = file
	Loger = log.New(LogFile, "[CW]", log.LstdFlags|log.Lshortfile|log.LUTC)
}

func Log(msg string) {
	Loger.Println(msg)
	time := time.Now().Format("15:04:05")
	fmt.Printf("[%s] %s\n", time, msg)
}
