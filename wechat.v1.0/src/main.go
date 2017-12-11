package main

import (
	"os"
	"wechat"
	"os/exec"
	"net/http"
	"path/filepath"
)


func main() {
	if os.Getppid()!=1{
		filePath,_ := filepath.Abs(os.Args[0])
		cmd := exec.Command(filePath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Start()
		return
	}
	http.HandleFunc("/wechat",wechat.AccessHandle)
	http.ListenAndServe(":3000",nil)
}