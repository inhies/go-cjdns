package admin

type SearchRunner struct{ client *Client }

func (s *SearchRunner) ShowActiveSearch(number int) error {
	_, err := s.client.sendCmd(&request{AQ: "SearchRunner",
		Args: &struct {
			Number int `bencode:"number"`
		}{number}})
	return err
}
