package admin

// Security_noFiles removes the ability to create new files from cjdns.
// It is recommended to always set this.
func (c *Conn) Security_noFiles() error {
	_, err := c.sendCmd(&request{AQ: "Security_noFiles"})
	return err
}

func (c *Conn) Security_setUser(user string) error {
	_, err := c.sendCmd(&request{
		AQ: "Security_setUser",
		Args: &struct {
			User string `bencode:"user"`
		}{user}})
	return err
}
