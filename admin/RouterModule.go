package admin

import "errors"

type RouterModule struct{ client *Client }

//RouterModule_lookup returns a single path for an address. Not sure what this is used for
func (r *RouterModule) Lookup(address string) (response map[string]interface{}, err error) {
	var (
		args = &struct {
			Address string `bencode:"address"`
		}{address}

		pack *packet
	)

	pack, err = r.client.sendCmd(&request{AQ: "RouterModule_lookup", Args: args})
	if err == nil {
		err = pack.Decode(response)
	}
	return
}

// Pings the specified IPv6 address or switch label and will timeout if it takes longer than the specified timeout period.
// CJDNS will fallback to its own timeout if the a zero timeout is given.
func (r *RouterModule) PingNode(addr string, timeout int) (ms int, version string, err error) {
	args := &struct {
		Path    string `bencode:"path"`
		Timeout int    `bencode:"timeout,omitempty"`
	}{addr, timeout}

	resp := new(struct {
		Ms      int    // number of milliseconds since the original ping
		Result  string // set when ping times out
		Version string // git hash of the source code which the node was built on
	})

	var pack *packet
	pack, err = r.client.sendCmd(&request{AQ: "RouterModule_pingNode", Args: args})
	if err == nil {
		err = pack.Decode(resp)
	}
	if err == nil && resp.Ms == 0 {
		err = errors.New(resp.Result)
	}
	return resp.Ms, resp.Version, err
}
