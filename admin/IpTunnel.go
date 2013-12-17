package admin

import (
	"github.com/inhies/go-cjdns/key"
	"net"
)

func (c *Conn) IpTunnel_allowConnection(publicKey key.Public, ip4Address, ip6Address net.IP) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_allowConnections",
		Args: &struct {
			Ip4    net.IP     `bencode:"ip4Address"`
			Ip6    net.IP     `bencode:"ip6Address"`
			PubKey key.Public `bencode:"publicKeyOfAuthorizedNode"`
		}{ip4Address, ip6Address, publicKey}})

	return err
}

func (c *Conn) IpTunnel_connectTo(publicKey key.Public) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_connectTo",
		Args: &struct {
			PubKey key.Public `bencode:"publicKeyOfNodeToConnectTo"`
		}{publicKey}})

	return err
}

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

func (c *Conn) IpTunnel_removeConnection(connection int) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_removeConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	return err
}

func (c *Conn) IpTunnel_showConnection(connection int) error {
	_, err := c.sendCmd(&request{AQ: "IpTunnel_showConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	return err
}
