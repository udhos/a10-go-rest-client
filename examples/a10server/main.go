package main

import (
	"fmt"
	"os"

	"github.com/sanity-io/litter"
	"github.com/udhos/a10-go-rest-client/a10go"
)

func main() {
	me := os.Args[0]
	if len(os.Args) != 4 {
		fmt.Printf("usage:   %s host         username password\n", me)
		fmt.Printf("example: %s 10.255.255.6 admin    a10\n", me)
		return
	}

	host := os.Args[1]
	user := os.Args[2]
	pass := os.Args[3]

	debug := os.Getenv("DEBUG") != ""

	fmt.Printf("%s: debug=%v DEBUG=[%s]\n", me, debug, os.Getenv("DEBUG"))

	c := a10go.New(host, a10go.Options{Debug: debug})

	errLogin := c.Login(user, pass)
	if errLogin != nil {
		fmt.Printf("login failure: %v", errLogin)
		return
	}

	fmt.Printf("before servers:\n")
	servers := c.ServerList()
	litter.Dump(servers)

	errCreate := c.ServerCreate("a10server_test00", "99.99.99.99", []string{"8888", "9999"})
	fmt.Printf("creating server: error:%v\n", errCreate)

	fmt.Printf("after servers:\n")
	servers = c.ServerList()
	litter.Dump(servers)

	errDel := c.ServerDelete("a10server_test00")
	fmt.Printf("deleting server: error:%v\n", errDel)

	fmt.Printf("final servers:\n")
	servers = c.ServerList()
	litter.Dump(servers)

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v", errLogout)
	}
}
