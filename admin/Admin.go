package admin

type Admin struct{ client *Client }

type AdminFunc struct {
	Type     string
	Required bool
}

// GetFunctions returns all available functions that cjdns supports
func (a *Admin) AvailableFunctions() (funcs map[string]map[string]AdminFunc, err error) {
	var (
		args = new(struct {
			Page int `bencode:"page"`
		})
		req = &request{Q: "Admin_availableFunctions", Args: args}

		resp = &struct {
			AvailableFunctions map[string]map[string]AdminFunc
			More               bool
		}{funcs, true}

		pack *packet
	)

	for resp.More {
		resp.More = false
		if pack, err = a.client.sendCmd(req); err == nil {
			err = pack.Decode(resp)
		}
		if err != nil {
			break
		}
		if len(resp.AvailableFunctions) == 0 {
			panic("empty response")
		}
		args.Page++
	}
	return resp.AvailableFunctions, err
}

// Checks with cjdns to see if asynchronous communication is allowed
func (a *Admin) AsyncEnabled() (bool, error) {
	res := new(struct{ AsyncEnabled bool })

	pack, err := a.client.sendCmd(&request{Q: "Admin_asyncEnabled"})
	if err == nil {
		err = pack.Decode(res)
	}
	return res.AsyncEnabled, err
}
