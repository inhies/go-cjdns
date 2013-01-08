// Package admin provides easy methods to the cjdns admin interface
package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/inhies/bencode"
	"github.com/kylelemons/godebug/pretty"
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

const (
	readerChanSize       = 10
	socketReaderChanSize = 100
	defaultPingTimeout   = 10000 //10 seconds
)

// Logger loops indefinitely and pushes log lines out on the channel
func Logger(user *Admin, info chan map[string]interface{}) {
	buf := make([]byte, 69632)
	remains := ""
	for {
		n, err := user.Conn.Read(buf[:])
		if err != nil {
			fmt.Printf("Error reading from socket: %+v\n\n", err)
			return
		}
		length := len(remains) + n
		newest := string(buf[0:n])
		//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
		combined := []byte(remains + newest)
		d := bencode.NewDecoder(combined[0:length])
		for !d.Consumed {
			o, err := d.Decode()
			if err != nil {
				//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
				remains = string(buf[d.Last:]) //PROBABLY NEED TO FIX THIS
				break
			}
			remains = remains[:0]
			info <- o.(map[string]interface{})
		}

	}
}

// Disconnect
func Disconnect(user *Admin) {
	//unsubscribe from logging and do other housekeeping
	user.Conn.Close()

}

func SendCmd(user *Admin, cmd string, args map[string]interface{}) (response map[string]interface{}, err error) {

	query := make(map[string]interface{})
	enc := bencode.NewEncoder()

	txid := strconv.FormatInt(rand.Int63(), 10) //randInt(10, 15)

	//Generate a unique transaction ID for this request
	query["txid"] = txid

	//Check if we need to use authentication for this command
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
		query["args"] = args

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

	//create the channel to receive data back on
	rChan := make(chan map[string]interface{}, 1) // use buffered channel to avoid blocking reader.

	user.Mu.Lock()
	//	println("NEW CHANNEL " + txid)
	user.Channels[txid] = rChan
	user.Mu.Unlock()

	// remove channel from map no matter how this function exits.
	defer func() {
		//println("DELETEING CHANNEL " + txid)
		user.Mu.Lock()
		delete(user.Channels, txid)
		user.Mu.Unlock()
	}()

	//Send the query

	//fmt.Printf("SendCmd SENDING: %v\n", string(enc.Bytes))

	_, err = io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return nil, err
	}

	output := <-rChan
	//	fmt.Printf("SendCmd RECEIVED: %v\n", output)
	delete(output, "txid")
	return output, nil

}

func Reader(user *Admin) {
	inChan := make(chan map[string]interface{}, socketReaderChanSize) //Channel we will receive data from the socket on
	//outChans := make(map[string]Response)                             //Map of channels we will send data out on where [] is unique ID

	go sockReader(user, inChan)
	for {
		//println("STATUS:")
		//pretty.Print(user.Channels)
		select {

		//outChan is a map of channels we will use to send data back on

		/* case newChan, ok := <-user.ReaderChan:
		if !ok {
			println("input channel closed, exiting...")
			return
		} else {
			fmt.Printf("NEWCHAN: %v\n", newChan)
			//get the 
			outChans[newChan.Key] = newChan
			newChan.Channel <- make(map[string]interface{}) //"OKLOL"
		}
		*/

		//inChan is decoded data from the socket
		case input, ok := <-inChan:
			if !ok {
				println("sockReader closed, re-starting...")
				go sockReader(user, inChan)
			} else {
				//fmt.Printf("RECEIVED: %v\n", input)
				var txid, sID string
				if input["txid"] != nil {
					txid = input["txid"].(string)
					//println("TXID " + txid)

				}
				if input["streamId"] != nil {
					sID = input["streamId"].(string)
					//println("STREAMID " + sID)
				}

				user.Mu.Lock()
				if _, ok := user.Channels[txid]; !ok {
					if _, ok := user.Channels[sID]; !ok {
						//We have no valid key!
						println("CHANNEL MISSING")
						pretty.Print(user.Channels)
						pretty.Print(input)
					} else {
						c := user.Channels[sID]
						user.Mu.Unlock()
						c <- input
					}
				} else {
					c := user.Channels[txid]
					user.Mu.Unlock()
					c <- input
				}

				//if c != nil {
				//	c <- input
				//} else {
				//	//pretty.Print(key)

				//	// worker exited. you might not care
				//}
			}

		}
	}

}

// sockReader continually reads from the socket and sends the data out
func sockReader(user *Admin, out chan<- map[string]interface{}) {
	buf := make([]byte, 69632)
	remains := ""
	for {

		n, err := user.Conn.Read(buf[:])

		if err != nil {
			return
		}
		length := len(remains) + n
		newest := string(buf[0:n])

		//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
		combined := []byte(remains + newest)
		d := bencode.NewDecoder(combined[0:length])
		for !d.Consumed {

			o, err := d.Decode()
			if err != nil {
				//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
				remains = string(buf[d.Last:]) //PROBABLY NEED TO FIX THIS
				fmt.Printf("ERROR! Read %v bytes, err: %+v\n", d.Last, err)
				break
			}
			remains = remains[:0]
			out <- o.(map[string]interface{})
			//t.Printf("RECEIVED: %v\n", o)

		}
	}

}

func sendOut(user *Admin, query map[string]interface{}) error {
	enc := bencode.NewEncoder()
	enc.Encode(query)
	//	println("MINE: " + string(enc.Bytes))

	/*data, errz := zeebo.EncodeString(query)
	if errz != nil {
		panic(errz)
	}
	println("ZEEB: " + data) */

	_, err := io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return err
	}

	return nil
}
func finaliseQuery(user *Admin, query map[string]interface{}) map[string]interface{} {
	cookie, _ := ReqCookie(user)
	query["q"] = "auth"
	query["cookie"] = cookie
	query["hash"] = sha256hash(user.Password + cookie)

	enc := bencode.NewEncoder()
	enc.Encode(query)
	data := string(enc.Bytes)

	query["hash"] = sha256hash(data)
	return query
}
func Connect(bind, pass string) (admin *Admin, success bool) {

	conn, err := net.Dial("tcp", bind)
	if err != nil {
		panic(err)
		return
	}
	//readerChannel := make(chan Response, readerChanSize)
	var l sync.Mutex
	admin = &Admin{bind, pass, conn, l, make(map[string]chan map[string]interface{})}
	go Reader(admin)
	_, err = sendPing(admin)

	if err != nil {
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())
	success = true

	return
}

func getResponse(user *Admin, out chan map[string]interface{}) error {

	buf := make([]byte, 69632)
	remains := ""

	n, err := user.Conn.Read(buf[:])

	if err != nil {
		return err
	}
	length := len(remains) + n
	newest := string(buf[0:n])

	//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
	combined := []byte(remains + newest)
	d := bencode.NewDecoder(combined[0:length])
	for !d.Consumed {
		o, err := d.Decode()
		if err != nil {
			//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
			remains = string(buf[d.Last:]) //PROBABLY NEED TO FIX THIS
			fmt.Printf("ERROR! Read %v bytes, err: %+v\n", d.Last, err)
			break
		}
		remains = remains[:0]
		out <- o.(map[string]interface{})

	}
	return nil
}

func sha256hash(input string) string {
	h := sha256.New()
	io.WriteString(h, input)
	hex := hex.EncodeToString([]byte(h.Sum(nil)))
	return hex
}

func randString(min, max int) string {
	r := myRand(min, max, "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

	return r
}

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
