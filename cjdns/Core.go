package cjdns

// Core_exit tells cjdns to shutdown
func (a *Admin) Core_exit() (response map[string]interface{}, err error) {
	response, err = SendCmd(a, "Core_exit", nil)
	if err != nil {
		return
	}

	return
}
