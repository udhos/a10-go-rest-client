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

	fmt.Printf("virtual servers:\n")
	vServers := c.VirtualServerList()
	litter.Dump(vServers)

	fmt.Printf("service groups:\n")
	sgroups := c.ServiceGroupList()
	litter.Dump(sgroups)

	fmt.Printf("servers:\n")
	servers := c.ServerList()
	litter.Dump(servers)

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v", errLogout)
	}
}
