package admin

import "github.com/inhies/go-cjdns/key"

// UDPInterface_beginConnection starts a direct connection to another node.
// Note that returned error only pertains to loading connection details,
// and will not convey the state of the connection itself.
//
// address has the form host:port.
func (a *Conn) UDPInterface_beginConnection(pubkey key.Public, address string, interfaceNumber int, password string) error {
	var (
		args = &struct {
			Address        string     `bencode:"address"`
			IntefaceNumber int        `bencode:"interfaceNumber,omitempty"`
			Password       string     `bencode:"password"`
			PublicKey      key.Public `bencode:"address"`
		}{address, interfaceNumber, password, pubkey}
		req  = &request{AQ: "UDPInterface_beginConnection", Args: args}
		resp = new(struct{ InterfaceNumber int })

		pack *packet
		err  error
	)

	if pack, err = a.sendCmd(req); err == nil {
		err = pack.Decode(resp)
	}
	return err
}

// UDPInterface_new creates a new UDPInterface which is either bound to an address/port or not and returns an index number for the interface.
//
// laddr has the form host:port, if host is unspecified, it is assumed to be `0.0.0.0`.
func (a *Conn) UDPInterface_new(laddr string) (interfaceNumber int, err error) {
	var (
		args = new(struct {
			Addr string `bencode:"bindAddress"`
		})
		req  = &request{AQ: "UDPInterface_new", Args: args}
		resp = new(struct{ InterfaceNumber int })

		pack *packet
	)

	if pack, err = a.sendCmd(req); err == nil {
		err = pack.Decode(resp)
	}
	return resp.InterfaceNumber, err
}
