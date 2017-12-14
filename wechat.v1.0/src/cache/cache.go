package cache
//版本
const VERSION = "1.0"
//消息队列长度
const CAPACITY = 1024

import (
	"net"
	"fmt"
	"errors"
	"io/ioutil"
	"regexp"
	"strings"
)

type Request struct{
	method string
	key string
	value string
	ch chan string
}

type Cache struct{
	//存放数据数组
	Data map[string]string
	//处理请求chan
	ReqC chan Request
}

func (cache Cache)checkerr(err interface{}) {
	if err != nil{
		errors.New("An Error Occured!")
	}
}
func (cache Cache)reply(conn net.Conn,result string){
	conn.Write([]byte(result))
	conn.Close();
}

func (cache *Cache)Start(ntp,addr string){
	//建立连接
	tcpAddr,err := net.ResolveTCPAddr(ntp,addr)
	cache.checkerr(err)
	tcpListener,err := net.ListenTCP(ntp,tcpAddr)
	cache.checkerr(err)
	cache.ReqC = make(chan Request,CAPACITY)
	cache.Data = make(map[string]string)
	go func() {
		for {
			cache.getData()
		}
	}()
	for {
		conn,err := tcpListener.Accept()
		cache.checkerr(err)
		go cache.handleConnection(conn)
	}
}

func (cache Cache)handleConnection(conn net.Conn){
	r,err := ioutil.ReadAll(conn)
	cache.checkerr(err)
	strRequest := string(r)
	strRequest = strings.TrimSpace(strRequest)
	arrRequest := regexp.MustCompile("\\s").Split(strRequest,3)
	var req Request
	var c chan string
	c = make(chan string)
	if arrRequest[0]=="set"{
		if len(arrRequest) < 3{
			cache.reply(conn,"nil")
			return
		}
		req = Request{
			method:"set",
			key:arrRequest[1],
			value:arrRequest[2],
			ch:c,
		}
	} else {
		if len(arrRequest) < 2{
			cache.reply(conn,"nil")
			return
		}
		req = Request{
			method:arrRequest[0],
			key:arrRequest[1],
			value:"",
			ch:c,
		}
	}
	fmt.Print("Request:"+strRequest+" ")
	cache.ReqC <- req
	response := <- req.ch
	fmt.Println("Response:"+response)
	cache.reply(conn,response)
}

func (cache *Cache)getData() {
	request := <-cache.ReqC
	if request.method == "get"{
		if _,ok := cache.Data[request.key];ok{
			request.ch <- cache.Data[request.key]
		} else {
			request.ch <- "nil"
		}
	} else if request.method == "set" {
		cache.Data[request.key] = request.value
		request.ch <- "ok"
	} else if request.method == "delete" {
		if _,ok := cache.Data[request.key];ok {
			delete(cache.Data,request.key)
		}
		request.ch <- "ok"
	}
}

//func main() {
//	cache := Cache{}
//	cache.Start("tcp",":3001")
//}