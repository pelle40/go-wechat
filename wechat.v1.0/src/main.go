package main

import (
	"fmt"
	"config"
)

func main() {
	fmt.Println("Hello world!")
	fmt.Println("Version "+config.VERSION)
	fmt.Printf("%s",config.VERSION)
}