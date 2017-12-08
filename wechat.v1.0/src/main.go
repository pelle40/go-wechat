package main

import (
	"config"
	//"encoding/json"
	"encoding/hex"
	"net/http"
	"io"
	"sort"
	"crypto/sha1"
)

func WechatHandle(w http.ResponseWriter,r *http.Request){
	r.ParseForm();
	if r.Method == "GET"{
		fieldlist := []string{"account","signature","timestamp","nonce","echostr"}
		for _,v := range fieldlist{
			value := r.Form[v]
			if value == nil{
				io.WriteString(w,"Miss param: "+v)
				return
			}
		}
		account := r.FormValue("account")
		signature := r.FormValue("signature")
		timestamp := r.FormValue("timestamp")
		nonce := r.FormValue("nonce")
		echostr := r.FormValue("echostr")

		if config.WechatConfig[account] == nil {
			io.WriteString(w,"Wrong param: account")
			return
		}

		keys := []string{config.WechatConfig[account].InterfaceToken,timestamp,nonce}
		sort.Strings(keys)
		strsha1 := ""
		for _,v := range keys {
			strsha1 += v
		}
		h := sha1.New()
		h.Write([]byte(strsha1))
		shabyte := h.Sum(nil)
		if hex.EncodeToString(shabyte) == signature {
			io.WriteString(w,echostr)
		} else {
			io.WriteString(w,"signatrue error")
		}
	} else if r.Method == "POST" {

	}
}
func main() {
	http.HandleFunc("/wechat",WechatHandle)
	http.ListenAndServe(":3000",nil)
}