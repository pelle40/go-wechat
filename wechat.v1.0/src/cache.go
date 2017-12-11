package main

import (
	"net"
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
)

var Data map[string]string
type Request struct{
	method string
	key string
	value string
	c chan string
}
var RequestChain chan Request

func handleCon(conn net.Conn,req chan Request){
	receive,_ := ioutil.ReadAll(conn)
	strRequest := string(receive)
	defer conn.Close()

	strRequest = strings.TrimSpace(strRequest)
	arrRequest := regexp.MustCompile("\\s+").Split(strRequest,3)
	var request Request
	res := make(chan string)
	if arrRequest[0] == "set"{
		if len(arrRequest) < 3{
			conn.Write([]byte("missing args"))
			return
		}
		request = Request{
			arrRequest[0],
			arrRequest[1],
			arrRequest[2],
			res,
		}
	} else {
		if len(arrRequest) < 2{
			conn.Write([]byte("missing args"))
			return
		}
		request = Request{
			arrRequest[0],
			arrRequest[1],
			"",
			res,
		}
	}
	fmt.Println(strRequest)
	req <- request
	response := <-res
	fmt.Println(response)
	conn.Write([]byte(response))
}

func main(){
	tcpAddr,err := net.ResolveTCPAddr("tcp",":3001")
	if err != nil{
		fmt.Println("Create TCP Addr error")
		return
	}
	tcpListener,err := net.ListenTCP("tcp",tcpAddr)
	if err != nil{
		fmt.Println("ListenTCP error")
		return
	}
	Data = make(map[string]string)
	RequestChain = make(chan Request,1024)
	//处理请求
	for {
		conn,err := tcpListener.Accept()
		if err == nil {
			go handleCon(conn,RequestChain)
		}
	}
	//处理缓存
	for {
		request := <-RequestChain
		if request.method == "get"{
			if _,ok := Data[request.key];ok{
				request.c <- Data[request.key]
			} else {
				request.c <- ""
			}
		} else if request.method == "set" {
			Data[request.key] = request.value
			request.c <- "1"
		} else if request.method == "delete" {
			if _,ok := Data[request.key];ok {
				delete(Data,request.key)
			}
			request.c <- "1"
		}
	}
}