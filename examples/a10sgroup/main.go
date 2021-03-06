package main

import (
	"fmt"
	"os"
	"strconv"

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

	serverCount, errServers := strconv.Atoi(os.Getenv("SERVERS"))
	if errServers != nil {
		fmt.Printf("%s: parsing SERVERS=[%s]: %v\n", me, os.Getenv("SERVERS"), errServers)
	}
	fmt.Printf("%s: serverCount=%d SERVERS=[%s]\n", me, serverCount, os.Getenv("SERVERS"))

	c := a10go.New(host, a10go.Options{Debug: debug})

	errLogin := c.Login(user, pass)
	if errLogin != nil {
		fmt.Printf("login failure: %v\n", errLogin)
		return
	}

	var serverList []string
	for i := 0; i < serverCount; i++ {
		s := fmt.Sprintf("a10sgroup_test%02d", i)
		a := fmt.Sprintf("99.99.99.%d", i)
		errCreate := c.ServerCreate(s, a, []string{"8888", "9999"})
		fmt.Printf("creating %d/%d server=[%s] error:%v\n", i, serverCount, s, errCreate)
		serverList = append(serverList, s)
	}

	var serverPortList []string
	for _, s := range serverList {
		serverPortList = append(serverPortList, s+",1111")
	}

	fmt.Printf("before service groups:\n")
	litter.Dump(c.ServiceGroupList())

	sgName := "a10sg_test00"

	create(c, sgName, serverPortList)
	create(c, sgName, serverPortList)

	fmt.Printf("after creating service groups:\n")
	litter.Dump(c.ServiceGroupList())

	update(c, sgName, serverPortList)

	fmt.Printf("after updating service groups - keeping members:\n")
	litter.Dump(c.ServiceGroupList())

	update(c, sgName, []string{})

	fmt.Printf("after updating service groups - removing members:\n")
	litter.Dump(c.ServiceGroupList())

	destroy(c, sgName)
	destroy(c, sgName)

	fmt.Printf("final service groups:\n")
	litter.Dump(c.ServiceGroupList())

	for i, s := range serverList {
		errDelete := c.ServerDelete(s)
		fmt.Printf("deleting %d/%d server=[%s] error:%v\n", i, serverCount, s, errDelete)
	}

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v\n", errLogout)
	}
}

func update(c *a10go.Client, sgName string, serverPortList []string) {
	proto := "2" // proto=TCP
	errUpdate := c.ServiceGroupUpdate(sgName, proto, serverPortList)
	fmt.Printf("updating service group=%s servers=%v error:%v\n", sgName, serverPortList, errUpdate)
}

func create(c *a10go.Client, sgName string, serverPortList []string) {
	proto := "2" // proto=TCP
	errCreate := c.ServiceGroupCreate(sgName, proto, serverPortList)
	fmt.Printf("creating service group=%s servers=%v error:%v\n", sgName, serverPortList, errCreate)
}

func destroy(c *a10go.Client, sgName string) {
	errDel := c.ServiceGroupDelete(sgName)
	fmt.Printf("deleting service group=%s error:%v\n", sgName, errDel)
}
