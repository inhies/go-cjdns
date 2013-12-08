package admin

type AdminFunc struct {
	Type     string
	Required bool
}

// GetFunctions returns all available functions that cjdns supports
func (a *Conn) Admin_availableFunctions() (functions map[string]map[string]AdminFunc, err error) {
	var (
		args = new(struct {
			Page int `bencode:"page"`
		})

		req = &request{Q: "Admin_availableFunctions", Args: args}

		res = new(struct {
			AvailableFunctions map[string]map[string]AdminFunc
			More               bool
		})

		pack *packet
	)

	functions = res.AvailableFunctions

	res.More = true
	for res.More {
		res.More = false
		if pack, err = a.sendCmd(req); err == nil {
			err = pack.Decode(res)
		}
		if err != nil {
			break
		}
		args.Page++
	}
	return
}

// Checks with cjdns to see if asynchronous communication is allowed
func (c *Conn) Admin_asyncEnabled() (bool, error) {
	res := new(struct{ AsyncEnabled bool })

	pack, err := c.sendCmd(&request{Q: "Admin_asyncEnabled"})
	if err == nil {
		err = pack.Decode(res)
	}
	return res.AsyncEnabled, err
}
