package admin

// Security_noFiles removes the ability to create new files from cjdns.
// It is recommended to always set this.
func (c *Conn) Security_noFiles() (response map[string]interface{}, err error) {
	var pack *packet
	pack, err = c.sendCmd(&request{AQ: "Security_noFiles"})
	if err == nil {
		err = pack.Decode(&response)
	}
	return
}

//Security_setUser(user)
