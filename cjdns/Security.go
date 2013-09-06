package cjdns



// Security_noFiles removes the ability to create new files from cjdns.
// It is recommended to always set this.
func Security_noFiles(user *Conn) (response map[string]interface{}, err error) {
	response, err = SendCmd(user, "Security_noFiles", nil)
	if err != nil {
		return
	}
	return
}


//Security_setUser(user)
