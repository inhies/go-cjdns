package cjdns

// Subscribes you to receive logging updates based on the parameters that are set.
// Returns a map[string]interface channel where all logging data will be sent
// and the stream ID cjdns uses to identify the subscription.
func (a *Admin) AdminLog_subscribe(file string, level string, line int) (channel chan map[string]interface{}, streamId string, err error) {

	args := make(map[string]interface{})

	args["file"] = file
	args["level"] = level
	args["line"] = line

	response, err := SendCmd(a, "AdminLog_subscribe", args)

	if err != nil {
		return
	}

	streamId = response["streamId"].(string)
	channel = make(chan map[string]interface{}, 100) // use buffered channel to avoid blocking reader.
	a.Mu.Lock()
	a.Channels[streamId] = channel
	a.Mu.Unlock()
	return

}

// Removes the logging subscription so that you no longer recieve log info.
func (a *Admin) AdminLog_unsubscribe(streamid string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["streamId"] = streamid
	response, err = SendCmd(a, "AdminLog_unsubscribe", args)
	if err != nil {
		return
	}
	return
}
