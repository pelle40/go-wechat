package main
//版本
//const VERSION = "1.0"

import (
	"os"
	"net"
	"fmt"
	"path"
	"time"
	"errors"
	"regexp"
	"strings"
	"os/exec"
	"io/ioutil"
	"path/filepath"
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
	//更新数据到文件
	UpFlag chan int
	//数据文件
	dataFile string
}

func (cache *Cache)Start(ntp,addr,datafile string,capcity int){
	//建立连接
	tcpAddr,err := net.ResolveTCPAddr(ntp,addr)
	cache.checkerr(err)
	tcpListener,err := net.ListenTCP(ntp,tcpAddr)
	cache.checkerr(err)
	cache.ReqC = make(chan Request,capcity)
	cache.Data = make(map[string]string)
	cache.UpFlag = make(chan int,1)
	cache.dataFile = datafile
	go func() {
		for {
			cache.getData()
		}
	}()
	go func() {
		for{
			_ := <-cache.UpFlag
			cache.updateDataToFile()
		}
	}
	for {
		conn,err := tcpListener.Accept()
		cache.checkerr(err)
		go cache.handleConnection(conn)
	}
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
		cache.UpFlag <- 1
		request.ch <- "ok"
	} else if request.method == "delete" {
		if _,ok := cache.Data[request.key];ok {
			delete(cache.Data,request.key)
		}
		cache.UpFlag <- 1
		request.ch <- "ok"
	}
}
func (cache *Cache)updateDataToFile() {
	//数据文件
	datafile,err := os.OpenFile(cache.dataFile,os.O_TRUNC,0666)
	if err!=nil{
		fmt.Println("Create or Open data file error!")
		return
	}
	for k,v := range cache.Data{
		datafile.WriteString(k+","+v+"\n")
	}
	datafile.Close()
}

func main() {
	folder,_ := os.Getwd()
	if os.Getppid()!=1{
		filepath,_ := filepath.Abs(os.Args[0])
		cmd := exec.Command(filepath)
		//日志文件
		filename := path.Join(folder,"log"+time.Now().Format("2006-01-02")+".log")
		openfile,err := os.OpenFile(filename,os.O_CREATE,0666)
		if err!=nil{
			return nil,errors.New("Create file error")
		}
		defer logfile.Close()
		//标准输入输出
		cmd.Stdin = logfile
		cmd.Stdout = logfile
		cmd.Stderr = logfile
		cmd.Start()
		return
	}
	cache := Cache{}
	datafilename := path.Join(folder,'data')
	cache.Start("tcp",":3001",datafilename,1024)
}