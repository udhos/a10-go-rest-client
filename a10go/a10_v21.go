package a10go

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"unicode"
)

// Client is an api client
type Client struct {
	host      string  // api host
	sessionID string  // session id
	opt       Options // client options
}

// FuncPrintf is function type for debug Printf
type FuncPrintf func(format string, v ...interface{})

// Options specify parameters for the api client
type Options struct {
	Debug       bool       // enable debugging
	DebugPrintf FuncPrintf // custom Printf function for debugging
}

func (c *Client) debugf(format string, v ...interface{}) {
	if c.opt.Debug {
		c.opt.DebugPrintf("DEBUG "+format, v...)
	}
}

// New creates api client
func New(host string, options Options) *Client {
	if options.DebugPrintf == nil {
		options.DebugPrintf = log.Printf // default debug Printf
	}
	return &Client{host: host, opt: options}
}

// Login opens a new session
func (c *Client) Login(username, password string) error {
	var errAuth error
	c.sessionID, errAuth = a10v21Auth(c.host, username, password)
	return errAuth
}

// Logout closes an existing session
func (c *Client) Logout() error {
	return a10v21Close(c.host, c.sessionID)
}

// Get calls http GET for an specific api method
func (c *Client) Get(method string) ([]byte, error) {
	return a10SessionGet(c.debugf, c.host, method, c.sessionID)
}

// Post calls http POST for an specific api method
func (c *Client) Post(method, body string) ([]byte, error) {
	return a10SessionPost(c.debugf, c.host, method, c.sessionID, body)
}

/*
// Delete calls http DELETE for an specific api method
func (c *Client) Delete(method, body string) ([]byte, error) {
	return a10SessionDelete(c.debugf, c.host, method, c.sessionID, body)
}
*/

// ServerList retrieves the full server list
func (c *Client) ServerList() []A10Server {
	return a10ServerList(c.debugf, c.host, c.sessionID)
}

// ServerCreate creates new server
func (c *Client) ServerCreate(name, host string, ports []string) error {

	format := `{
            "server": {
                "name": "%s",
                "host": "%s",
                "status": 1,
		"port_list": [%s]
            }
        }
`

	portList := ""
	for _, p := range ports {
		portName, portProto := splitPortProto(c.debugf, p)
		portFmt := portFormat(portName, portProto)
		if portList == "" {
			portList = portFmt
			continue
		}
		portList += "," + portFmt
	}

	payload := fmt.Sprintf(format, name, host, portList)

	body, errPost := c.Post("slb.server.create", payload)

	c.debugf("ServerCreate: reqPayload=[%s] respBody=[%s] error=[%v]", payload, body, errPost)

	return errPost
}

func splitPortProto(debugf FuncPrintf, portProto string) (string, string) {
	s := strings.FieldsFunc(portProto, isSep)
	count := len(s)
	switch {
	case count < 1:
		proto := "2"
		debugf("splitPortProto(%s): defaulting to port protocol=%s", portProto, proto)
		return "", proto
	case count < 2:
		proto := "2"
		debugf("splitPortProto(%s): defaulting to port protocol=%s", portProto, proto)
		return s[0], proto
	}
	return s[0], s[1]
}

func portFormat(port, protocol string) string {
	return fmt.Sprintf(`{"port_num": %s, "protocol": %s}`, port, protocol)
}

// ServerDelete deletes an existing server
func (c *Client) ServerDelete(name string) error {

	format := `{ "server": { "name": "%s" } }`

	payload := fmt.Sprintf(format, name)

	body, errDelete := c.Post("slb.server.delete", payload)

	c.debugf("ServerDelete: reqPayload=[%s] respBody=[%s] error=[%v]", payload, body, errDelete)

	return errDelete
}

// ServiceGroupList retrieves the full server group list
func (c *Client) ServiceGroupList() []A10ServiceGroup {
	return a10ServiceGroupList(c.debugf, c.host, c.sessionID)
}

// ServiceGroupCreate creates new service group
func (c *Client) ServiceGroupCreate(name, protocol string, members []string) error {

	format := `{
            "service_group": {
                "name": "%s",
                "protocol": %s,
		"member_list": [%s]
            }
        }
`

	serverList := ""
	for _, s := range members {
		serverName, serverPort, serverProto := splitServerPortProto(c.debugf, s)
		serverFmt := serverFormat(serverName, serverPort, serverProto)
		if serverList == "" {
			serverList = serverFmt
			continue
		}
		serverList += "," + serverFmt
	}

	payload := fmt.Sprintf(format, name, protocol, serverList)

	body, errPost := c.Post("slb.service_group.create", payload)

	c.debugf("ServiceGroupCreate: reqPayload=[%s] respBody=[%s] error=[%v]", payload, body, errPost)

	return errPost
}

func splitServerPortProto(debugf FuncPrintf, serverPort string) (string, string, string) {
	s := strings.FieldsFunc(serverPort, isSep)
	count := len(s)
	if count < 1 {
		proto := "2"
		debugf("splitServerPortProto(%s): count=%d defaulting to port protocol=%s", serverPort, count, proto)
		return "", "", proto
	}
	if count < 2 {
		proto := "2"
		debugf("splitServerPortProto(%s): count=%d defaulting to port protocol=%s", serverPort, count, proto)
		return s[0], "", proto
	}
	if count < 3 {
		proto := "2"
		debugf("splitServerPortProto(%s): count=%d defaulting to port protocol=%s", serverPort, count, proto)
		return s[0], s[1], proto
	}
	return s[0], s[1], s[2]
}

func isSep(c rune) bool {
	return c == ',' || unicode.IsSpace(c)
}

func serverFormat(name, port, proto string) string {
	return fmt.Sprintf(`{"server": "%s", "port": %s, "protocol": %s}`, name, port, proto)
}

// ServiceGroupDelete deletes an existing service group
func (c *Client) ServiceGroupDelete(name string) error {

	format := `{ "name": "%s" }`

	payload := fmt.Sprintf(format, name)

	body, errDelete := c.Post("slb.service_group.delete", payload)

	c.debugf("ServiceGroupDelete: reqPayload=[%s] respBody=[%s] error=[%v]", payload, body, errDelete)

	return errDelete
}

// VirtualServerList retrieves the full virtual server list
func (c *Client) VirtualServerList() []A10VServer {
	return a10VirtualServerList(c.debugf, c.host, c.sessionID)
}

// A10VServer is a virtual server
type A10VServer struct {
	Name          string
	Address       string
	Port          string
	ServiceGroups []string
}

// A10ServiceGroup is a service group
type A10ServiceGroup struct {
	Name     string
	Protocol string
	Members  []A10SGMember
}

// A10SGMember is a service group member
type A10SGMember struct {
	Name     string
	Port     string
	Protocol string
}

// A10Server is a server
type A10Server struct {
	Name  string
	Host  string
	Ports []A10Port
}

// A10Port defines port/protocol
type A10Port struct {
	Number   string
	Protocol string
}

// V3:
//
// Source: https://github.com/a10networks/tps-scripts/blob/master/axapi_curl_example.txt
//
// curl -k -X POST -H 'content-type: application/json' -d '{"credentials": {"username": "admin", "password": "a10"}}' 'https://192.168.199.152/axapi/v3/auth'
//
// V2:
//
// Source: https://www.a10networks.com/resources/articles/axapi-python
//
// https://10.255.255.6/services/rest/V2/?method=authenticate&username=admin&password=a10&format=json
//
// V2.1:
//
// Source: https://github.com/a10networks/acos-client/blob/master/acos_client/v21/session.py
//
// url:       /services/rest/v2.1/?format=json&method=authenticate
// post body: { "username": username, "password": password }

func a10v21url(host, method string) string {
	return "https://" + host + "/services/rest/v2.1/?format=json&method=" + method
}

func a10v21urlSession(host, method, sessionID string) string {
	return a10v21url(host, method) + "&session_id=" + sessionID
}

func mapGetStr(debugf FuncPrintf, tab map[string]interface{}, key string) string {
	value, found := tab[key]
	if !found {
		debugf("mapGetStr: key=[%s] not found", key)
		return ""
	}
	str, isStr := value.(string)
	if !isStr {
		debugf("mapGetStr: key=[%s] non-string value: [%v]", key, value)
		return ""
	}
	return str
}

func mapGetValue(debugf FuncPrintf, tab map[string]interface{}, key string) string {
	value, found := tab[key]
	if !found {
		debugf("mapGetValue: key=[%s] not found", key)
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func a10ServerList(debugf FuncPrintf, host, sessionID string) []A10Server {
	var list []A10Server

	servers, errGet := a10SessionGet(debugf, host, "slb.server.getAll", sessionID)
	if errGet != nil {
		return list
	}

	sList := jsonExtractList(debugf, servers, "server_list")
	if sList == nil {
		return list
	}

	for _, s := range sList {
		sMap, isMap := s.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(debugf, sMap, "name")
		host := mapGetStr(debugf, sMap, "host")
		server := A10Server{Name: name, Host: host}

		debugf("server: %s", name)

		portList := sMap["port_list"]
		pList, isList := portList.([]interface{})
		if !isList {
			continue
		}
		for _, p := range pList {
			pMap, isPMap := p.(map[string]interface{})
			if !isPMap {
				continue
			}
			portNum := mapGetValue(debugf, pMap, "port_num")
			proto := mapGetValue(debugf, pMap, "protocol")
			server.Ports = append(server.Ports, A10Port{Number: portNum, Protocol: proto})
		}

		list = append(list, server)
	}

	return list
}

func a10ServiceGroupList(debugf FuncPrintf, host, sessionID string) []A10ServiceGroup {
	var list []A10ServiceGroup

	groups, errGet := a10SessionGet(debugf, host, "slb.service_group.getAll", sessionID)
	if errGet != nil {
		return list
	}

	sgList := jsonExtractList(debugf, groups, "service_group_list")
	if sgList == nil {
		return list
	}

	for _, sg := range sgList {
		sgMap, isMap := sg.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(debugf, sgMap, "name")
		protocol := mapGetValue(debugf, sgMap, "protocol")
		group := A10ServiceGroup{Name: name, Protocol: protocol}

		debugf("service group: %s", name)

		memberList := sgMap["member_list"]
		mList, isList := memberList.([]interface{})
		if isList {
			for _, m := range mList {
				mMap, isMMap := m.(map[string]interface{})
				if !isMMap {
					continue
				}
				memberName := mapGetStr(debugf, mMap, "server")
				memberPort := mapGetValue(debugf, mMap, "port")
				member := A10SGMember{Name: memberName, Port: memberPort}
				group.Members = append(group.Members, member)
			}
		}

		list = append(list, group)
	}

	return list
}

func a10VirtualServerList(debugf FuncPrintf, host, sessionID string) []A10VServer {
	var list []A10VServer

	bodyVirtServers, errGet := a10SessionGet(debugf, host, "slb.virtual_server.getAll", sessionID)
	if errGet != nil {
		return list
	}

	vsList := jsonExtractList(debugf, bodyVirtServers, "virtual_server_list")
	if vsList == nil {
		return list
	}

	for _, vs := range vsList {
		vsMap, isMap := vs.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(debugf, vsMap, "name")
		addr := mapGetStr(debugf, vsMap, "address")

		debugf("virtual server: %s", name)

		vServer := A10VServer{Name: name, Address: addr}

		portList := vsMap["vport_list"]
		pList, isList := portList.([]interface{})
		if !isList {
			continue
		}
		for _, vp := range pList {
			pMap, isPMap := vp.(map[string]interface{})
			if !isPMap {
				continue
			}
			pStr := mapGetValue(debugf, pMap, "port")
			sGroup := mapGetStr(debugf, pMap, "service_group")

			vServer.Port = pStr
			vServer.ServiceGroups = append(vServer.ServiceGroups, sGroup)
			debugf("virtual server: %s service_group=%s", name, sGroup)
		}

		list = append(list, vServer)
	}

	return list
}

func jsonExtractList(debugf FuncPrintf, body []byte, listName string) []interface{} {
	me := "extractList"
	tab := map[string]interface{}{}
	errJSON := json.Unmarshal(body, &tab)
	if errJSON != nil {
		log.Printf(me+": list=%s json error: %v", listName, errJSON)
		return nil
	}
	list, found := tab[listName]
	if !found {
		debugf(me+": list=%s not found", listName)
		return nil
	}
	slice, isSlice := list.([]interface{})
	if !isSlice {
		debugf(me+": list=%s not an slice", listName)
		return nil
	}
	return slice
}

func a10SessionGet(debugf FuncPrintf, host, method, sessionID string) ([]byte, error) {
	me := "a10SessionGet"
	api := a10v21urlSession(host, method, sessionID)
	debugf(me+": url=[%s]", api)
	body, err := httpGet(api)
	if err != nil {
		debugf(me+": api=[%s] error: %v", api, err)
	}
	return body, err
}

func a10SessionPost(debugf FuncPrintf, host, method, sessionID, body string) ([]byte, error) {
	me := "a10SessionPost"
	api := a10v21urlSession(host, method, sessionID)
	debugf(me+": url=[%s]", api)
	respBody, err := httpPostString(api, contentTypeJSON, body)
	if err != nil {
		debugf(me+": api=[%s] error: %v", api, err)
	}
	return respBody, err
}

/*
func a10SessionDelete(debugf FuncPrintf, host, method, sessionID, body string) ([]byte, error) {
	me := "a10SessionDelete"
	api := a10v21urlSession(host, method, sessionID)
	respBody, err := httpDeleteString(api, contentTypeJSON, body)
	if err != nil {
		debugf(me+": api=[%s] error: %v", api, err)
	}
	return respBody, err
}
*/

const contentTypeJSON = "application/json"

func a10v21Close(host, sessionID string) error {

	api := a10v21urlSession(host, "session.close", sessionID)

	format := `{"session_id": "%s"}`
	payload := fmt.Sprintf(format, sessionID)

	_, errPost := httpPostString(api, contentTypeJSON, payload)

	return errPost
}

func a10v21Auth(host, username, password string) (string, error) {

	body, errAuth := v21auth(host, username, password)
	if errAuth != nil {
		return "", errAuth
	}

	response := map[string]interface{}{}

	errJSON := json.Unmarshal(body, &response)
	if errJSON != nil {
		return "", errJSON
	}

	id, found := response["session_id"]
	if !found {
		return "", fmt.Errorf("auth response missing session_id")
	}

	sessionID, isStr := id.(string)
	if !isStr {
		return "", fmt.Errorf("auth session_id not a string")
	}

	return sessionID, nil
}

func v21auth(host, username, password string) ([]byte, error) {

	api := a10v21url(host, "authenticate")

	format := `{ "username": "%s", "password": "%s" }`
	payload := fmt.Sprintf(format, username, password)

	return httpPostString(api, contentTypeJSON, payload)
}
