package admin

import (
	"github.com/inhies/go-cjdns/key"
	"net"
)

type IPTunnel struct{ client *Client }

func (i *IPTunnel) AllowConnection(publicKey *key.Public, addr net.IP) (err error) {
	if b := addr.To4; b != nil {
		_, err = i.client.sendCmd(&request{AQ: "IpTunnel_allowConnection",
			Args: &struct {
				Ip     net.IP      `bencode:"ip4Address"`
				PubKey *key.Public `bencode:"publicKeyOfAuthorizedNode"`
			}{addr, publicKey}})
	} else {
		_, err = i.client.sendCmd(&request{AQ: "IpTunnel_allowConnections",
			Args: &struct {
				Ip     net.IP      `bencode:"ip6Address"`
				PubKey *key.Public `bencode:"publicKeyOfAuthorizedNode"`
			}{addr, publicKey}})
	}
	return
}

func (i *IPTunnel) ConnectTo(publicKey *key.Public) error {
	_, err := i.client.sendCmd(&request{AQ: "IpTunnel_connectTo",
		Args: &struct {
			PubKey *key.Public `bencode:"publicKeyOfNodeToConnectTo"`
		}{publicKey}})

	return err
}

// IpTunnel_listConnections returns a list of all current IP tunnels
func (i *IPTunnel) ListConnections() (tunnelIndexes []int, err error) {
	resp := new(struct {
		List []int
	})

	var pack *packet
	pack, err = i.client.sendCmd(&request{AQ: "IpTunnel_listConnections"})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp.List, err
}

func (i *IPTunnel) RemoveConnection(connection int) error {
	_, err := i.client.sendCmd(&request{AQ: "IpTunnel_removeConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	return err
}

type IpTunnelConnection struct {
	Ip4Address *net.IP
	Ip6Address *net.IP
	Key        *key.Public
	Outgoing   bool
}

func (i *IPTunnel) ShowConnection(connection int) (*IpTunnelConnection, error) {
	resp := new(IpTunnelConnection)

	pack, err := i.client.sendCmd(&request{AQ: "IpTunnel_showConnection",
		Args: &struct {
			Connection int `bencode:"connection"`
		}{connection}})
	if err == nil {
		err = pack.Decode(resp)
	}
	return resp, err
}
