package admin

type Security struct{ client *Client }

// Security_setUser sets the user ID which cjdns is running under to a different user.
// This function allows cjdns to shed privileges after starting up.
func (s *Security) SetUser(user string) error {
	_, err := s.client.sendCmd(&request{
		AQ: "Security_setUser",
		Args: &struct {
			User string `bencode:"user"`
		}{user}})
	return err
}
