// Package cjdns provides easy methods to the cjdns admin interface
package cjdns

import (
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/inhies/bencode"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os/user"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type CjdnsAdminConfig struct {
	Addr     string `json:"addr"`
	Port     int    `json:"port"`
	Password string `json:"password"`
}

// Contains the admin info for connecting to cjdns
type Conn struct {
	Conn     net.Conn
	Mu       sync.Mutex
	Channels map[string]chan map[string]interface{}
	password string
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

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}
func stripComments(b []byte) ([]byte, error) {
	regComment, err := regexp.Compile("(?s)//.*?\n|/\\*.*?\\*/")
	if err != nil {
		return nil, err
	}
	out := regComment.ReplaceAllLiteral(b, nil)
	return out, nil
}

func Connect(config *CjdnsAdminConfig) (admin *Conn, err error) {
	if config == nil {
		config = new(CjdnsAdminConfig)
		u, err := user.Current()
		if err != nil {
			return nil, err
		}

		rawFile, err := ioutil.ReadFile(u.HomeDir + "/.cjdnsadmin")
		if err != nil {
			return nil, err
		}

		raw, err := stripComments(rawFile)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(raw, &config)
		if err != nil {
			return nil, err
		}
	}

	addr := &net.UDPAddr{
		IP:   net.ParseIP(config.Addr),
		Port: config.Port,
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}

	var l sync.Mutex

	admin = &Conn{
		password: config.Password,
		Conn:     conn,
		Mu:       l,
		Channels: make(map[string]chan map[string]interface{}),
	}

	go Reader(admin)
	_, err = admin.SendPing(1000)

	if err != nil {
		return
	}
	return
}

func SendCmd(user *Conn, cmd string, args map[string]interface{}) (response map[string]interface{}, err error) {
	query := make(map[string]interface{})
	enc := bencode.NewEncoder()

	// Generate a unique transaction ID for this request.
	// This tells Reader which channel to send the data back to us on.
	txid := strconv.FormatInt(rand.Int63(), 10) //randInt(10, 15)
	query["txid"] = txid

	// Check if we need to use authentication for this command.
	if cmd == "ping" || cmd == "cookie" || cmd == "Admin_asyncEnabled" {
		query["q"] = cmd
	} else {

		//Request a cookie
		cookie, err := user.ReqCookie()
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
		query["hash"] = sha256hash(user.password + cookie)

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

		// If an error field exists, and we have an error, return it
		if _, ok := response["error"]; ok {
			if response["error"] != "none" {
				err = fmt.Errorf(response["error"].(string))
				return
			}
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
func Reader(user *Conn) {
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
func sockReader(user *Conn, out chan<- map[string]interface{}) {
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
func sendOut(user *Conn, query map[string]interface{}) error {
	enc := bencode.NewEncoder()
	enc.Encode(query)
	_, err := io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return err
	}
	return nil
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

// base32 encodes the public key for use in the decoder
func EncodePubKey(in []byte) (out []byte) {
	var wide, bits uint
	var i2b = []byte("0123456789bcdfghjklmnpqrstuvwxyz")

	for len(in) > 0 {
		// Add the 8 bits of data from the next `in` byte above the existing bits
		wide, in, bits = wide|uint(in[0])<<bits, in[1:], bits+8
		for bits > 5 {
			// Remove the least significant 5 bits and add their character to out
			wide, out, bits = wide>>5, append(out, i2b[int(wide&0x1F)]), bits-5
		}
	}
	// If it wasn't a precise multiple of 40 bits, add some padding based on the remaining bits
	if bits > 0 {
		out = append(out, i2b[int(wide)])
	}
	return out
}

//// Converts a cjdns public key to an ipv6 address
//func PubKeyToIP(in []byte) (outString string, err error) {
//	// Check for the trailing .k
//	if in[len(in)-2] == '.' && in[len(in)-1] == 'k' {
//		in = in[0 : len(in)-2]
//	}

//	var wide, bits uint
//	var out []byte
//	var i2b = []byte("0123456789bcdfghjklmnpqrstuvwxyz")
//	var b2i = func() []byte {
//		var ascii [256]byte
//		for i := range ascii {
//			ascii[i] = 255
//		}
//		for i, b := range i2b {
//			ascii[b] = byte(i)
//		}
//		return ascii[:]
//	}()

//	for len(in) > 0 && in[0] != '=' {
//		// Add the 5 bits of data corresponding to the next `in` character above existing bits
//		wide, in, bits = wide|uint(b2i[int(in[0])])<<bits, in[1:], bits+5
//		if bits >= 8 {
//			// Remove the least significant 8 bits of data and add it to out
//			wide, out, bits = wide>>8, append(out, byte(wide)), bits-8
//		}
//	}

//	// If there was padding, there will be bits left, but they should be zero
//	if wide != 0 {
//		err = fmt.Errorf("extra data at end of decode")
//		return
//	}

//	// Do the hashing that generates the IP
//	out = sha512hash(sha512hash(out))
//	if out[0] != 0xfc {
//		err = fmt.Errorf("invalid")
//		return
//	}

//	out = out[0:16]

//	// Assemble the IP
//	for i := 0; i < 16; i++ {
//		if i > 0 && i < 16 && i%2 == 0 {
//			outString += ":"
//		}
//		outString += fmt.Sprintf("%02x", out[i])
//	}
//	return
//}
func sha512hash(input []byte) []byte {
	h := sha512.New()
	h.Write(input)
	return []byte(h.Sum(nil))
}
