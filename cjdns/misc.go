package cjdns

import (
	"fmt"
	"time"
)

// Memory is supposed to return information on cjdns's memory use but it currently
// just crashes it.
func Memory(user *Admin) (memory int64, err error) {
	response, err := SendCmd(user, "memory", nil)
	if err != nil {
		return
	}
	memory = response["bytes"].(int64)
	return
}

// This will return a page from cjdns's routing table.
func NodeStore_dumpTable(user *Admin, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["page"] = page
	response, err = SendCmd(user, "NodeStore_dumpTable", args)
	if err != nil {
		return
	}
	return
}

// Requests a cookie from cjdns and returns it.
func ReqCookie(user *Admin) (cookie string, err error) {
	response, err := SendCmd(user, "cookie", nil)
	if err != nil {
		return
	}
	cookie = response["cookie"].(string)
	return
}

// Sends a ping to cjdns and returns true if a pong was received
// before the specified timeout.
func SendPing(user *Admin, timeout time.Duration) (bool, error) {
	rChan := make(chan map[string]interface{}, 1)
	go func() {
		response, err := SendCmd(user, "ping", nil)
		if err != nil {
			return
		}
		rChan <- response
	}()

	timeout *= 1000 * 1000

	reply := make(map[string]interface{})
	var err error
	var ok bool
	select {
	case reply, ok = <-rChan:
		if !ok {
			fmt.Errorf("error reading ping response from cjdns.")
		}
	case <-time.After(timeout):
		err = fmt.Errorf("cjdns is not responding, you may need to restart it.")
	}

	if err != nil {
		return false, err
	}
	if reply["q"] != "pong" {
		err := fmt.Errorf("Did not recieve pong from cjdns.")
		return false, err
	}
	return true, nil
}
