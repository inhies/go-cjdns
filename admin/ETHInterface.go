package admin

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
func (c *Conn) ETHInterface_beacon(iface int, state int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["interfaceNumber"] = iface
	args["state"] = state
	response, err = SendCmd(c, "ETHInterface_beacon", args)
	if err != nil {
		return
	}
	return
}

//Initiates a connection to the specified node
func (c *Conn) ETHInterface_beginConnection(iface int, mac string, pass string, pubkey string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["interfaceNumber"] = iface
	args["macAddress"] = mac
	args["password"] = pass
	args["publicKey"] = pubkey
	response, err = SendCmd(c, "ETHInterface_beginConnection", args)
	if err != nil {
		return
	}
	return
}

//ETHInterface_new creates a new ethernet interface
func (c *Conn) ETHInterface_new(device string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["bindDevice"] = device
	response, err = SendCmd(c, "ETHInterface_new", args)
	if err != nil {
		return
	}
	return
}
