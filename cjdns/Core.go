package cjdns

// Core_exit tells cjdns to shutdown
func (c *Conn) Core_exit() (response map[string]interface{}, err error) {
	response, err = SendCmd(c, "Core_exit", nil)
	if err != nil {
		return
	}

	return
}
