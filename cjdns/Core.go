package cjdns


// Core_exit tells cjdns to shutdown
func Core_exit(user *Admin) (response map[string]interface{}, err error) {
	response, err = SendCmd(user, "Core_exit", nil)
	if err != nil {
		return
	}

	return
}

