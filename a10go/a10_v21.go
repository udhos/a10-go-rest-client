package a10go

import (
	"encoding/json"
	"fmt"
	"log"
)

type A10Client struct {
	host      string
	sessionId string
}

func New(host string) *A10Client {
	return &A10Client{host: host}
}

func (c *A10Client) Login(username, password string) error {
	var errAuth error
	c.sessionId, errAuth = a10v21Auth(c.host, username, password)
	return errAuth
}

func (c *A10Client) Logout() error {
	return a10v21Close(c.host, c.sessionId)
}

func (c *A10Client) Get(method string) ([]byte, error) {
	return a10SessionGet(c.host, method, c.sessionId)
}

func (c *A10Client) ServerList() []A10Server {
	return a10ServerList(c.host, c.sessionId)
}

func (c *A10Client) ServiceGroupList() []A10ServiceGroup {
	return a10ServiceGroupList(c.host, c.sessionId)
}

func (c *A10Client) VirtualServerList() []A10VServer {
	return a10VirtualServerList(c.host, c.sessionId)
}

type A10VServer struct {
	Name          string
	Address       string
	Port          string
	ServiceGroups []string
}

type A10ServiceGroup struct {
	Name    string
	Members []A10SGMember
}

type A10SGMember struct {
	Name string
	Port string
}

type A10Server struct {
	Name  string
	Host  string
	Ports []string
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

func a10v21urlSession(host, method, sessionId string) string {
	return a10v21url(host, method) + "&session_id=" + sessionId
}

func mapGetStr(tab map[string]interface{}, key string) string {
	value, found := tab[key]
	if !found {
		log.Printf("mapGetStr: key=[%s] not found", key)
		return ""
	}
	str, isStr := value.(string)
	if !isStr {
		log.Printf("mapGetStr: key=[%s] non-string value: [%v]", key, value)
		return ""
	}
	return str
}

func mapGetValue(tab map[string]interface{}, key string) string {
	value, found := tab[key]
	if !found {
		log.Printf("mapGetValue: key=[%s] not found", key)
		return ""
	}
	return fmt.Sprintf("%v", value)
}

func a10ServerList(host, sessionId string) []A10Server {
	var list []A10Server

	servers, errGet := a10SessionGet(host, "slb.server.getAll", sessionId)
	if errGet != nil {
		return list
	}

	//log.Printf("servers: [%s]", string(servers))

	sList := jsonExtractList(servers, "server_list")
	if sList == nil {
		return list
	}

	for _, s := range sList {
		sMap, isMap := s.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(sMap, "name")
		host := mapGetStr(sMap, "host")
		server := A10Server{Name: name, Host: host}

		log.Printf("server: %s", name)

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
			portNum := mapGetValue(pMap, "port_num")
			server.Ports = append(server.Ports, portNum)
		}

		list = append(list, server)
	}

	return list
}

func a10ServiceGroupList(host, sessionId string) []A10ServiceGroup {
	var list []A10ServiceGroup

	groups, errGet := a10SessionGet(host, "slb.service_group.getAll", sessionId)
	if errGet != nil {
		return list
	}

	//log.Printf("groups: [%s]", string(groups))

	sgList := jsonExtractList(groups, "service_group_list")
	if sgList == nil {
		return list
	}

	for _, sg := range sgList {
		sgMap, isMap := sg.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(sgMap, "name")
		group := A10ServiceGroup{Name: name}

		log.Printf("service group: %s", name)

		memberList := sgMap["member_list"]
		mList, isList := memberList.([]interface{})
		if isList {
			for _, m := range mList {
				mMap, isMMap := m.(map[string]interface{})
				if !isMMap {
					continue
				}
				memberName := mapGetStr(mMap, "server")
				memberPort := mapGetValue(mMap, "port")
				member := A10SGMember{Name: memberName, Port: memberPort}
				group.Members = append(group.Members, member)
			}
		}

		list = append(list, group)
	}

	return list
}

func a10VirtualServerList(host, sessionId string) []A10VServer {
	var list []A10VServer

	bodyVirtServers, errGet := a10SessionGet(host, "slb.virtual_server.getAll", sessionId)
	if errGet != nil {
		return list
	}

	vsList := jsonExtractList(bodyVirtServers, "virtual_server_list")
	if vsList == nil {
		return list
	}

	for _, vs := range vsList {
		vsMap, isMap := vs.(map[string]interface{})
		if !isMap {
			continue
		}

		name := mapGetStr(vsMap, "name")
		addr := mapGetStr(vsMap, "address")

		log.Printf("virtual server: %s", name)

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
			pStr := mapGetValue(pMap, "port")
			sGroup := mapGetStr(pMap, "service_group")

			vServer.Port = pStr
			vServer.ServiceGroups = append(vServer.ServiceGroups, sGroup)
			log.Printf("virtual server: %s service_group=%s", name, sGroup)
		}

		list = append(list, vServer)
	}

	return list
}

func jsonExtractList(body []byte, listName string) []interface{} {
	me := "extractList"
	tab := map[string]interface{}{}
	errJson := json.Unmarshal(body, &tab)
	if errJson != nil {
		log.Printf(me+": list=%s json error: %v", listName, errJson)
		return nil
	}
	list, found := tab[listName]
	if !found {
		log.Printf(me+": list=%s not found", listName)
		return nil
	}
	slice, isSlice := list.([]interface{})
	if !isSlice {
		log.Printf(me+": list=%s not an slice", listName)
		return nil
	}
	return slice
}

func a10SessionGet(host, method, sessionId string) ([]byte, error) {
	me := "a10SessionGet"
	api := a10v21urlSession(host, method, sessionId)
	body, err := httpGet(api)
	if err != nil {
		log.Printf(me+": api=[%s] error: %v", api, err)
	}
	return body, err
}

func a10v21Close(host, sessionId string) error {

	api := a10v21urlSession(host, "session.close", sessionId)

	format := `{"session_id": "%s"}`
	payload := fmt.Sprintf(format, sessionId)

	_, errPost := httpPostString(api, "application/json", payload)

	return errPost
}

func a10v21Auth(host, username, password string) (string, error) {

	body, errAuth := v21auth(host, username, password)
	if errAuth != nil {
		return "", errAuth
	}

	response := map[string]interface{}{}

	errJson := json.Unmarshal(body, &response)
	if errJson != nil {
		return "", errJson
	}

	id, found := response["session_id"]
	if !found {
		return "", fmt.Errorf("auth response missing session_id")
	}

	session_id, isStr := id.(string)
	if !isStr {
		return "", fmt.Errorf("auth session_id not a string")
	}

	return session_id, nil
}

func v21auth(host, username, password string) ([]byte, error) {

	api := a10v21url(host, "authenticate")

	format := `{ "username": "%s", "password": "%s" }`
	payload := fmt.Sprintf(format, username, password)

	return httpPostString(api, "application/json", payload)
}
