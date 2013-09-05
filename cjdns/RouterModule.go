package cjdns

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

