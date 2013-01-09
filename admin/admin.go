// Package admin provides easy methods to the cjdns admin interface
package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/inhies/bencode"
	"io"
	"math/rand"
	"net"
	"strconv"
	"sync"
	"time"
)

// Contains the admin info for connecting to cjdns
type Admin struct {
	Address  string
	Password string
	Conn     net.Conn
	Mu       sync.Mutex
	Channels map[string]chan map[string]interface{}
}

type PingResponse struct {
	Time    int64
	Result  string
	Version string
	Error   string
}

const (
	readerChanSize       = 10
	socketReaderChanSize = 100
)

func SendCmd(user *Admin, cmd string, args map[string]interface{}) (response map[string]interface{}, err error) {
	query := make(map[string]interface{})
	enc := bencode.NewEncoder()

	// Generate a unique transaction ID for this request.
	// This tells Reader which channel to send the data back to us on.
	txid := strconv.FormatInt(rand.Int63(), 10) //randInt(10, 15)
	query["txid"] = txid

	// Check if we need to use authentication for this command.
	if cmd == "ping" || cmd == "cookie" {
		query["q"] = cmd
	} else {

		//Request a cookie
		cookie, err := ReqCookie(user)
		if err != nil {
			return nil, err
		}

		query["q"] = "auth"
		query["aq"] = cmd
		query["cookie"] = cookie
		if args != nil {
			query["args"] = args
		}
		//Generate the first hash we need
		query["hash"] = sha256hash(user.Password + cookie)

		//Encode the query and get the final hash
		enc.Encode(query)
		data := string(enc.Bytes)
		query["hash"] = sha256hash(data)
	}

	//Re-encode the query with the new hash
	enc = bencode.NewEncoder()
	enc.Encode(query)

	// If we are telling cjdns to exit, then we will not get a response, so there is no need to wait
	if cmd != "Core_exit" {
		//create the channel to receive data back on
		rChan := make(chan map[string]interface{}, 1) // use buffered channel to avoid blocking reader.

		user.Mu.Lock()
		user.Channels[txid] = rChan
		user.Mu.Unlock()

		// remove channel from map no matter how this function exits.
		defer func() {
			user.Mu.Lock()
			delete(user.Channels, txid)
			user.Mu.Unlock()
		}()

		//Send the query
		_, err = io.WriteString(user.Conn, string(enc.Bytes))
		if err != nil {
			return nil, err
		}

		output, ok := <-rChan
		if !ok {
			return nil, fmt.Errorf("Socket closed")
		}
		return output, nil
	}
	//Send the query
	_, err = io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return nil, err
	}
	return make(map[string]interface{}), nil
}

// Collects data from sockReader and sends it out on the correct channel as designated by 
// the txid or streamId fields.
func Reader(user *Admin) {
	//Create a channel and launch the go routine that actually reads from the socket
	inChan := make(chan map[string]interface{}, socketReaderChanSize)
	go sockReader(user, inChan)
	for input := range inChan {
		//Check for a txid and a streamId
		//both of which can appear
		var txid, sID string
		if input["txid"] != nil {
			txid = input["txid"].(string)
		}
		if input["streamId"] != nil {
			sID = input["streamId"].(string)
		}

		user.Mu.Lock()
		if _, ok := user.Channels[txid]; !ok {
			if _, ok := user.Channels[sID]; !ok {
				//We have no valid key!
				//panic("CHANNEL MISSING")
				continue

			} else {
				c := user.Channels[sID]
				user.Mu.Unlock()
				delete(input, "txid")
				c <- input
			}
		} else {
			c := user.Channels[txid]
			user.Mu.Unlock()
			delete(input, "txid")
			c <- input
		}
	}
}

// sockReader continually reads from the socket and sends the data out
func sockReader(user *Admin, out chan<- map[string]interface{}) {
	buf := make([]byte, 69632)
	remains := ""
	errCount := 0
	for {
		n, err := user.Conn.Read(buf[:])
		if err != nil {
			close(out)
			return
		}
		length := len(remains) + n
		newest := string(buf[0:n])

		// BUG(inhies): Switch from using strings as a workaround to null bytes
		// http://golang.org/pkg/bytes/#IndexByte
		combined := []byte(remains + newest)
		d := bencode.NewDecoder(combined[0:length])
		for !d.Consumed {
			o, err := d.Decode()
			if err != nil {
				errCount++
				// BUG(inhies): need to add an error recovery function where we will 
				// increment the start point of the read by 1 until we get a valid response,
				// then discard all previous data
				if errCount >= 10 {
					remains = ""
					break
				}
				remains = string(buf[d.Last:])
				break
			}
			remains = remains[:0]
			out <- o.(map[string]interface{})
		}
		errCount = 0
	}
}

// Writes data to the socket of the specified connection
func sendOut(user *Admin, query map[string]interface{}) error {
	enc := bencode.NewEncoder()
	enc.Encode(query)
	_, err := io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return err
	}
	return nil
}

// Connects to a running cjdns instance
func Connect(bind, pass string) (admin *Admin, err error) {
	conn, err := net.DialTimeout("tcp", bind, 2e9) // BUG(inhies): default timeout is 2 seconds. Add an option to make it user configurable
	if err != nil {
		return
	}
	var l sync.Mutex
	admin = &Admin{bind, pass, conn, l, make(map[string]chan map[string]interface{})}
	go Reader(admin)
	_, err = SendPing(admin, 1000)

	if err != nil {
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())

	return
}

// Hashes a string and returns it
func sha256hash(input string) string {
	h := sha256.New()
	io.WriteString(h, input)
	hex := hex.EncodeToString([]byte(h.Sum(nil)))
	return hex
}

// Returns a random alphanumeric string where length is <= max >= min
func randString(min, max int) string {
	r := myRand(min, max, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")
	return r
}

// Returns a random character from the specified string where length is <= max >= min
func myRand(min, max int, char string) string {

	var length int

	if min < max {
		length = min + rand.Intn(max-min)
	} else {
		length = min
	}

	buf := make([]byte, length)
	for i := 0; i < length; i++ {
		buf[i] = char[rand.Intn(len(char)-1)]
	}
	return string(buf)
}
