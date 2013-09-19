package admin

import (
	"fmt"
	"time"
)

//ETHInterface_beacon will set the specified beacon state on the specified interface
//State is any of the following: 
//0 -- Disabled.
//1 -- Accept beacons, this will cause cjdns to accept incoming
//beacon messages and try connecting to the sender.
//2 -- Accept and send beacons, this will cause cjdns to broadcast
//messages on the local network which contain a randomly
//generated per-session password, other nodes which have this
//set to 1 or 2 will hear the beacon messages and connect
//automatically.
func ETHInterface_beacon(user *Admin, iface int, state int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["interfaceNumber"] = iface
	args["state"] = state
	response, err = SendCmd(user, "ETHInterface_beacon", args)
	if err != nil {
		return
	}
	return
}

//Initiates a connection to the specified node
func ETHInterface_beginConnection(user *Admin, iface int, mac string, pass string, pubkey string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["interfaceNumber"] = iface
	args["macAddress"] = mac
	args["password"] = pass
	args["publicKey"] = pubkey
	response, err = SendCmd(user, "ETHInterface_beginConnection", args)
	if err != nil {
		return
	}
	return
}

//ETHInterface_new creates a new ethernet interface
func ETHInterface_new(user *Admin, device string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["bindDevice"] = device
	response, err = SendCmd(user, "ETHInterface_new", args)
	if err != nil {
		return
	}
	return
}

//RouterModule_lookup returns a single path for an address. Not sure what this is used for
func RouterModule_lookup(user *Admin, address string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["address"] = address
	response, err = SendCmd(user, "RouterModule_lookup", args)
	if err != nil {
		return
	}
	return
}

// Core_exit tells cjdns to shutdown
func Core_exit(user *Admin) (response map[string]interface{}, err error) {
	response, err = SendCmd(user, "Core_exit", nil)
	if err != nil {
		return
	}

	return
}

// Security_noFiles removes the ability to create new files from cjdns.
// It is recommended to always set this.
func Security_noFiles(user *Admin) (response map[string]interface{}, err error) {
	response, err = SendCmd(user, "Security_noFiles", nil)
	if err != nil {
		return
	}
	return
}

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

// InterfaceController_peerStats returns a page of stats about direct peers
func InterfaceController_peerStats(user *Admin, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["page"] = page

	response, err = SendCmd(user, "InterfaceController_peerStats", args)
	if err != nil {
		return
	}
	return
}

// IpTunnel_listConnections returns a list of all current IP tunnels
func IpTunnel_listConnections(user *Admin) (response map[string]interface{}, err error) {
	response, err = SendCmd(user, "IpTunnel_listConnections", nil)
	if err != nil {
		return
	}
	return
}

// GetFunctions returns all available functions that cjdns supports
func Admin_availableFunctions(user *Admin, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})

	args["page"] = page

	response, err = SendCmd(user, "Admin_availableFunctions", args)
	if err != nil {
		return
	}
	return
}

// Subscribes you to receive logging updates based on the parameters that are set.
// Returns a map[string]interface channel where all logging data will be sent
// and the stream ID cjdns uses to identify the subscription.
func AdminLog_subscribe(user *Admin, file string, level string, line int) (channel chan map[string]interface{}, streamId string, err error) {

	args := make(map[string]interface{})

	args["file"] = file
	args["level"] = level
	args["line"] = line

	response, err := SendCmd(user, "AdminLog_subscribe", args)

	if err != nil {
		return
	}

	streamId = response["streamId"].(string)
	channel = make(chan map[string]interface{}, 100) // use buffered channel to avoid blocking reader.
	user.Mu.Lock()
	user.Channels[streamId] = channel
	user.Mu.Unlock()
	return

}

// Removes the logging subscription so that you no longer recieve log info.
func AdminLog_unsubscribe(user *Admin, streamid string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["streamId"] = streamid
	response, err = SendCmd(user, "AdminLog_unsubscribe", args)
	if err != nil {
		return
	}
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

// Pings the specified IPv6 address and will timeout if it takes longer than the specified timeout period.
func RouterModule_pingNode(user *Admin, addr string, timeout int) (data PingResponse, err error) {
	args := make(map[string]interface{})
	args["path"] = addr
	args["timeout"] = timeout
	response, err := SendCmd(user, "RouterModule_pingNode", args)

	if err != nil {
		return
	}

	if _, ok := response["error"]; ok { //check if an error was sent
		data.Error = response["error"].(string)

	} else if response["result"] == "timeout" { //check if we had a timeout
		data.Time = response["ms"].(int64)
		data.Result = response["result"].(string)

	} else { //everything must be fine!
		data.Time = response["ms"].(int64)
		data.Result = response["result"].(string)
		data.Version = response["version"].(string)
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

// Checks with cjdns to see if asynchronous communication is allowed
func Admin_asyncEnabled(user *Admin) (enabled bool, err error) {
	response, err := SendCmd(user, "Admin_asyncEnabled", nil)
	if err != nil {
		return
	}

	if response["asyncEnabled"].(int64) == 1 {
		enabled = true
	}

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
