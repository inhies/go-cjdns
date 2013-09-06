package cjdns


// GetFunctions returns all available functions that cjdns supports
func Admin_availableFunctions(user *Conn, page int) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})

	args["page"] = page

	response, err = SendCmd(user, "Admin_availableFunctions", args)
	if err != nil {
		return
	}
	return
}


// Checks with cjdns to see if asynchronous communication is allowed
func Admin_asyncEnabled(user *Conn) (enabled bool, err error) {
	response, err := SendCmd(user, "Admin_asyncEnabled", nil)
	if err != nil {
		return
	}

	if response["asyncEnabled"].(int64) == 1 {
		enabled = true
	}

	return
}

