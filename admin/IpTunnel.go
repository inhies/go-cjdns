package admin

// IpTunnel_listConnections returns a list of all current IP tunnels
func (c *Conn) IpTunnel_listConnections() (tunnelIndexes []int, err error) {
	resp := new(struct {
		List []int
	})

	var pack *packet
	pack, err = c.sendCmd(&request{AQ: "IpTunnel_listConnections"})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp.List, err
}

//IpTunnel_allowConnection(publicKeyOfAuthorizedNode, ip6Address=0, ip4Address=0)
//IpTunnel_connectTo(publicKeyOfNodeToConnectTo)
//IpTunnel_removeConnection(connection)
//IpTunnel_showConnection(connection)
