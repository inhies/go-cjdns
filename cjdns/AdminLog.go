package cjdns



// Subscribes you to receive logging updates based on the parameters that are set.
// Returns a map[string]interface channel where all logging data will be sent
// and the stream ID cjdns uses to identify the subscription.
func AdminLog_subscribe(user *Conn, file string, level string, line int) (channel chan map[string]interface{}, streamId string, err error) {

	args := make(map[string]interface{})

	args["file"] = file
	args["level"] = level
	args["line"] = line

	response, err := SendCmd(user, "AdminLog_subscribe", args)

	if err != nil {
		return
	}

	streamId = response["streamId"].(string)
	channel = make(chan map[string]interface{}, 100) // use buffered channel to avoid blocking reader.
	user.Mu.Lock()
	user.Channels[streamId] = channel
	user.Mu.Unlock()
	return

}

// Removes the logging subscription so that you no longer recieve log info.
func AdminLog_unsubscribe(user *Conn, streamid string) (response map[string]interface{}, err error) {
	args := make(map[string]interface{})
	args["streamId"] = streamid
	response, err = SendCmd(user, "AdminLog_unsubscribe", args)
	if err != nil {
		return
	}
	return
}
