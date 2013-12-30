package admin

import "errors"

// Memory returns the number of bytes allocated by all memory allocators
// in the router.
func (a *Conn) Memory() (int, error) {
	r := new(struct{ Bytes int })
	pack, err := a.sendCmd(&request{Q: "memory"})
	if err == nil {
		err = pack.Decode(r)
	}

	return r.Bytes, err
}

// Ping sends a ping to cjdns and returns true if a pong was received.
func (a *Conn) Ping() error {
	pack, err := a.sendCmd(&request{Q: "ping"})
	if err == nil {
		r := new(struct {
			Q     string
			Error string
		})
		err = pack.Decode(r)
		if r.Q != "pong" {
			err = errors.New("did not receive pong.")
		}
	}
	return nil
}

// authedPing sends an "authorized" ping to cjdns and returns an error if a
// pong is not recieved
func (a *Conn) authedPing() error {
	pack, err := a.sendCmd(&request{AQ: "ping"})
	if err == nil {
		r := new(struct {
			Q     string
			Error string
		})
		err = pack.Decode(r)
		if r.Q != "pong" {
			err = errors.New("did not receive pong.")
		}
	}
	return nil
}
