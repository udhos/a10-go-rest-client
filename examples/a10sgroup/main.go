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

	fmt.Printf("before service groups:\n")
	sgroups := c.ServiceGroupList()
	litter.Dump(sgroups)

	sgName := "a10sg_test00"

	errCreate := c.ServiceGroupCreate(sgName, []string{"s2", "s3"})
	fmt.Printf("creating service group: error:%v\n", errCreate)

	fmt.Printf("after service groups:\n")
	sgroups = c.ServiceGroupList()
	litter.Dump(sgroups)

	errDel := c.ServiceGroupDelete(sgName)
	fmt.Printf("deleting server: error:%v\n", errDel)

	fmt.Printf("final service groups:\n")
	sgroups = c.ServiceGroupList()
	litter.Dump(sgroups)

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v", errLogout)
	}
}
