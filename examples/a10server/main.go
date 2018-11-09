package main

import (
	"fmt"
	"os"
	"strings"

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

	serverName := os.Getenv("SERVER_NAME")
	if serverName == "" {
		serverName = "a10server_test00"
	}
	fmt.Printf("%s: serverName=%s SERVER_NAME=[%s]\n", me, serverName, os.Getenv("SERVER_NAME"))

	portList := os.Getenv("PORTS")
	if portList == "" {
		portList = "8888 9999"
	}
	ports := strings.Fields(portList)
	fmt.Printf("%s: ports=%v PORTS=[%s]\n", me, ports, os.Getenv("PORTS"))

	c := a10go.New(host, a10go.Options{Debug: debug})

	errLogin := c.Login(user, pass)
	if errLogin != nil {
		fmt.Printf("login failure: %v", errLogin)
		return
	}

	fmt.Printf("\nbefore servers:\n")
	litter.Dump(c.ServerList())

	create(c, serverName, ports)
	create(c, serverName, ports)

	fmt.Printf("\nafter servers:\n")
	litter.Dump(c.ServerList())

	update(c, "intentional-non-existant-server-name", ports)

	p7 := []string{"7777"}
	update(c, serverName, p7) // will add the port to list

	fmt.Printf("\nafter updating ports=%v:\n", p7)
	litter.Dump(c.ServerList())

	update(c, serverName, nil) // will clear port list

	fmt.Printf("\nafter updating ports=nil:\n")
	litter.Dump(c.ServerList())

	update(c, serverName, p7)

	fmt.Printf("\nafter updating ports=%v:\n", p7)
	litter.Dump(c.ServerList())

	destroy(c, serverName)
	destroy(c, serverName)

	fmt.Printf("\nfinal servers:\n")
	litter.Dump(c.ServerList())

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v", errLogout)
	}
}

func create(c *a10go.Client, serverName string, ports []string) {
	errCreate := c.ServerCreate(serverName, "99.99.99.99", ports)
	fmt.Printf("creating server=%s ports=%v error:%v\n", serverName, ports, errCreate)
}

func update(c *a10go.Client, serverName string, ports []string) {
	errUpdate := c.ServerUpdate(serverName, "99.99.99.99", ports)
	fmt.Printf("updating server=%s ports=%v error:%v\n", serverName, ports, errUpdate)
}

func destroy(c *a10go.Client, serverName string) {
	errDel := c.ServerDelete(serverName)
	fmt.Printf("deleting server=%s error:%v\n", serverName, errDel)
}
