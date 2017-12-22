package main

import (
	"fmt"
	"net"
	"bufio"
	"os"
	"io"
	"strings"
)

func main() {
	raddr,_ := net.ResolveTCPAddr("tcp",":3001")
	conn,_ := net.DialTCP("tcp",nil,raddr)
	args := strings.Join(os.Args[1:]," ")
	fmt.Println(args)
	conn.Write([]byte(args))
	conn.CloseWrite();
	data,err := bufio.NewReader(conn).ReadString('\n')
	if err != nil && err != io.EOF{
		fmt.Println(err)
		return
	}
	fmt.Println(data)
	fmt.Println("End")
	conn.Close()
}