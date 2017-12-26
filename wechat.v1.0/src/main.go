package main

import (
	"os"
	"path"
	"wechat"
	"os/exec"
	"net/http"
	"path/filepath"
)


func main() {
	if os.Getppid()!=1{
		filePath,_ := filepath.Abs(os.Args[0])
		folder,_ := os.Getwd()
		logfile := path.Join(folder,'wechat.server.log')
		logfile,err := os.OpenFile(logfile,os.O_CREATE,0666)
		if err!=nil{
			return
		}
		cmd := exec.Command(filePath)
		cmd.Stdin = logfile
		cmd.Stdout = logfile
		cmd.Stderr = logfile
		cmd.Start()
		return
	}
	http.HandleFunc("/wechat",wechat.AccessHandle)
	http.ListenAndServe(":3000",nil)
}