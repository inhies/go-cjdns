package cjdns

// Subscribes you to receive logging updates based on the parameters that are set.
// Returns a map[string]interface channel where all logging data will be sent
// and the stream ID cjdns uses to identify the subscription.
func (c *Conn) AdminLog_subscribe(file string, level string, line int) (channel chan map[string]interface{}, streamId string, err error) {

	args := make(map[string]interface{})

	args["file"] = file
	args["level"] = level
	args["line"] = line

	response, err := SendCmd(c, "AdminLog_subscribe", args)

	if err != nil {
		return
	}

	streamId = response["streamId"].(string)
	channel = make(chan map[string]interface{}, 100) // use buffered channel to avoid blocking reader.
	c.Mu.Lock()
	c.Channels[streamId] = channel
	c.Mu.Unlock()
	return

}

// Removes the logging subscription so that you no longer recieve log info.
func (c *Conn) AdminLog_unsubscribe(streamid string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["streamId"] = streamid
	response, err = SendCmd(c, "AdminLog_unsubscribe", args)
	if err != nil {
		return
	}
	return
}
