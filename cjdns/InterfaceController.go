package cjdns

// Returns stats on currently connected peers
func InterfaceController_peerStats(user *Conn, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})

	args["page"] = page

	response, err = SendCmd(user, "InterfaceController_peerStats", args)
	if err != nil {
		return
	}
	return
}

//InterfaceController_disconnectPeer(pubkey)

