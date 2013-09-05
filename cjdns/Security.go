package cjdns



// Security_noFiles removes the ability to create new files from cjdns.
// It is recommended to always set this.
func (a *Admin) Security_noFiles(user *Admin) (response map[string]interface{}, err error) {
	response, err = SendCmd(a, "Security_noFiles", nil)
	if err != nil {
		return
	}
	return
}


//Security_setUser(user)
