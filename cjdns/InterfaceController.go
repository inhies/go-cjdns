package cjdns

type PeerStats struct {
	Last        int64
	SwitchLabel string
	BytesOut    int64
	PublicKey   string
	State       string
	IsIncoming  bool
	BytesIn     int64
}

// Returns stats on currently connected peers
func InterfaceController_peerStats(user *Conn, page int) (response []PeerStats, err error) {
	args := make(map[string]interface{})

	args["page"] = page

	data, err := SendCmd(user, "InterfaceController_peerStats", args)
	if err != nil {
		return
	}

	// Convert the map to a slice of structs.
	// This should be fixed so ALL functions return structs... eventually...
	response = make([]PeerStats, 0)
	for _, peer := range data["peers"].([]interface{}) {
		info := peer.(map[string]interface{})
		var incoming bool
		if info["isIncoming"].(int64) > 0 {
			incoming = true
		}
		peer := PeerStats{
			Last:        info["last"].(int64),
			BytesIn:     info["bytesIn"].(int64),
			BytesOut:    info["bytesOut"].(int64),
			IsIncoming:  incoming,
			State:       info["state"].(string),
			PublicKey:   info["publicKey"].(string),
			SwitchLabel: info["switchLabel"].(string),
		}
		response = append(response, peer)
	}
	return
}

//InterfaceController_disconnectPeer(pubkey)
