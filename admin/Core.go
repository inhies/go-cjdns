package admin

import "errors"

type Core struct{ client *Client }

// Core_exit tells cjdns to shutdown
func (c *Core) Exit() error {
	resp := new(struct{ Error string })

	pack, err := c.client.sendCmd(&request{AQ: "Core_exit"})
	if err == nil {
		err = pack.Decode(resp)
		if err == nil && resp.Error != "none" {
			err = errors.New(resp.Error)
		}
	}
	return err
}
