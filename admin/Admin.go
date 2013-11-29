package admin

// GetFunctions returns all available functions that cjdns supports
func (c *Conn) Admin_availableFunctions() (functions map[string]string, err error) {
	var page int
	more := true
	args := make(map[string]interface{})
	var response map[string]interface{}

	for more {
		args["page"] = page
		response, err = SendCmd(c, "Admin_availableFunctions", args)
		if err != nil {
			return
		}
		more = (response["more"].(int64) == 1)
		page++
	}

	functions = make(map[string]string)
	for k, v := range response["availableFunctions"].(map[string]string) {
		functions[k] = v
	}
	return
}

// Checks with cjdns to see if asynchronous communication is allowed
func (c *Conn) Admin_asyncEnabled() (enabled bool, err error) {
	response, err := SendCmd(c, "Admin_asyncEnabled", nil)
	if err != nil {
		return
	}

	if response["asyncEnabled"].(int64) == 1 {
		enabled = true
	}
	return
}
