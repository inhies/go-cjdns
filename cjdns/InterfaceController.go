package cjdns

// Returns stats on currently connected peers
func (a *Admin) InterfaceController_peerStats(user *Admin, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})

	args["page"] = page

	response, err = SendCmd(a, "InterfaceController_peerStats", args)
	if err != nil {
		return
	}
	return
}

//InterfaceController_disconnectPeer(pubkey)

