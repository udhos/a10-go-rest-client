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

	fmt.Printf(me + " version 0.0\n")

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

	groupKey := "SGROUPS" // GROUP is a bash reserved var
	groupValue := os.Getenv(groupKey)
	groupCount, errGroups := strconv.Atoi(groupValue)
	if errGroups != nil {
		fmt.Printf("%s: parsing %s=[%s]: %v\n", me, groupKey, groupValue, errGroups)
	}
	fmt.Printf("%s: groupCount=%d %s=[%s]\n", me, groupCount, groupKey, groupValue)

	c := a10go.New(host, a10go.Options{Debug: debug})

	errLogin := c.Login(user, pass)
	if errLogin != nil {
		fmt.Printf("login failure: %v\n", errLogin)
		return
	}

	fmt.Printf("\n#### before:\n\n")
	litter.Dump(c.VirtualServerList())
	litter.Dump(c.ServiceGroupList())
	litter.Dump(c.ServerList())

	groupList, serverList, _ := createGroups(c, groupCount, serverCount)
	if groupList == nil {
		fmt.Printf("empty group list\n")
		return
	}

	var virtualPorts []string
	for _, g := range groupList {
		virtualPorts = append(virtualPorts, g+",8888")
	}

	vServerName := "a10vserver_vs00"

	createVS(c, vServerName, virtualPorts)
	createVS(c, vServerName, virtualPorts)

	fmt.Printf("\n#### after:\n\n")
	litter.Dump(c.VirtualServerList())
	litter.Dump(c.ServiceGroupList())
	litter.Dump(c.ServerList())

	deleteVS(c, vServerName)
	deleteVS(c, vServerName)

	deleteGroups(c, groupList)
	deleteServers(c, serverList)

	fmt.Printf("\n#### final:\n\n")
	litter.Dump(c.VirtualServerList())
	litter.Dump(c.ServiceGroupList())
	litter.Dump(c.ServerList())

	errLogout := c.Logout()
	if errLogout != nil {
		fmt.Printf("logout failure: %v\n", errLogout)
	}
}

func createVS(c *a10go.Client, name string, ports []string) {
	errCreate := c.VirtualServerCreate(name, "88.88.88.88", ports)
	fmt.Printf("creating virtuaServer=%s: %v\n", name, errCreate)
}

func deleteVS(c *a10go.Client, name string) {
	errDelete := c.VirtualServerDelete(name)
	fmt.Printf("deleting virtualServer=%s: %v\n", name, errDelete)
}

func createGroups(c *a10go.Client, groupCount int, serverCount int) ([]string, []string, []string) {

	var groupList, totalServerList, totalPortList []string

	for i := 0; i < groupCount; i++ {

		prefixServer := fmt.Sprintf("a10vgroup_sg%02d_server", i)
		prefixAddr := fmt.Sprintf("99.99.%d.", i)

		// create servers for group
		serverList, serverPortList := createServers(c, prefixServer, prefixAddr, serverCount)
		if serverList == nil {
			fmt.Printf("empty server list\n")
			return nil, nil, nil
		}

		// create group
		sgName := fmt.Sprintf("a10vgroup_sg%02d", i)
		proto := "2" // proto=TCP
		errCreate := c.ServiceGroupCreate(sgName, proto, serverPortList)
		fmt.Printf("creating %d/%d service_group=[%s] error:%v\n", i, groupCount, sgName, errCreate)
		if errCreate != nil {
			return nil, nil, nil
		}

		groupList = append(groupList, sgName)
		totalServerList = append(totalServerList, serverList...)
		totalPortList = append(totalPortList, serverPortList...)
	}

	return groupList, totalServerList, totalPortList
}

func deleteGroups(c *a10go.Client, groupList []string) {
	for i, g := range groupList {
		errDelete := c.ServiceGroupDelete(g)
		fmt.Printf("deleting %d/%d group=[%s] error:%v\n", i, len(groupList), g, errDelete)
	}
}

func createServers(c *a10go.Client, prefixServer, prefixAddr string, serverCount int) ([]string, []string) {
	var serverList []string
	var serverPortList []string

	for i := 0; i < serverCount; i++ {
		s := fmt.Sprintf(prefixServer+"%02d", i)
		a := fmt.Sprintf(prefixAddr+"%d", i)
		errCreate := c.ServerCreate(s, a, []string{"8888", "9999"})
		fmt.Printf("creating %d/%d server=[%s] error:%v\n", i, serverCount, s, errCreate)
		if errCreate != nil {
			return nil, nil
		}
		serverList = append(serverList, s)
	}

	for _, s := range serverList {
		serverPortList = append(serverPortList, s+",1111")
	}

	return serverList, serverPortList
}

func deleteServers(c *a10go.Client, serverList []string) {
	for i, s := range serverList {
		errDelete := c.ServerDelete(s)
		fmt.Printf("deleting %d/%d server=[%s] error:%v\n", i, len(serverList), s, errDelete)
	}
}
