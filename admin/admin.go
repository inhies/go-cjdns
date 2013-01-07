// Package admin provides easy methods to the cjdns admin interface
package admin

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/inhies/bencode"
	"io"
	"net"
)

// Contains the admin info for connecting to cjdns
type Admin struct {
	Address  string
	Password string
	Port     string
	Cookie   string
	Conn     net.Conn
}

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
	println("Closing connection")
	/*channel := make(chan map[string]interface{})
	go getResponse(user, channel)
	response := <-channel
	fmt.Printf("END: %#v\n", response)*/
	user.Conn.Close()
}

func sendOut(user *Admin, query map[string]interface{}) bool {
	enc := bencode.NewEncoder()
	enc.Encode(query)
	println(string(enc.Bytes))
	_, err := io.WriteString(user.Conn, string(enc.Bytes))
	if err != nil {
		return false
	}

	return true
}
func finaliseQuery(user *Admin, query map[string]interface{}) map[string]interface{} {
	cookie, _ := ReqCookie(user)
	query["cookie"] = cookie
	query["hash"] = sha256hash(user.Password + cookie)

	enc := bencode.NewEncoder()
	enc.Encode(query)
	data := string(enc.Bytes)

	query["hash"] = sha256hash(data)
	return query
}
func Connect(bind, pass string) (*Admin, bool) {
	c, err := net.Dial("tcp", bind)
	if err != nil {
		panic(err)
		return &Admin{}, false
	}
	admin := &Admin{bind, pass, "", "", c}
	defer c.Close()

	ping, err := Ping(admin)
	if !ping {
		return &Admin{}, false
	}

	return admin, true
}
func getResponse(user *Admin, out chan map[string]interface{}) error {
	//println("get reponse called")
	buf := make([]byte, 69632)
	remains := ""
	//println("reading")
	n, err := user.Conn.Read(buf[:])
	//println("read")
	if err != nil {
		fmt.Printf("Error reading from socket: %+v\n\n", err)
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
			//	println("err")
			//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
			remains = string(buf[d.Last:]) //PROBABLY NEED TO FIX THIS
			//		println("***********************************************************************************************")
			fmt.Printf("ERROR! Read %v bytes, err: %+v\n", d.Last, err)
			//		println("***********************************************************************************************")

			//	fmt.Printf("RECEIVED: %+v\n\n", combined)
			//fmt.Printf("REMAINS: %+v\n\n", remains)
			break
		}
		remains = remains[:0]
		out <- o.(map[string]interface{})

	}
	return nil
	/*
		n, err := user.Conn.Read(buf[:])
		//fmt.Printf("RECEIVED: %+v\n\n", string(buf))
		if err != nil {
			fmt.Printf("Error reading from socket: %+v\n\n", err)
			return err
		}

		d := bencode.NewDecoder(buf[0:n])
		var result map[string]interface{}
		for !d.Consumed {
			o, err := d.Decode()
			//t := reflect.TypeOf(o)
			if err != nil {
				//fmt.Printf("\n\nERROR! Read %v bytes, o: %s, err: %+v\n\n", d.Last, t, err)
				buf = buf[d.Last:] //PROBABLY NEED TO FIX THIS
				//fmt.Printf("REMAINS: %+v\n\n", string(buf))
				break
			}
			result = o.(map[string]interface{})
			//	fmt.Printf("Read %d bytes of obj(%s): %#v\n\n", d.Last, t, o)
		}
		//fmt.Printf("MINE: %+v\n\n", result)
		out <- result

		return nil
	*/

}

func sha256hash(input string) string {
	h := sha256.New()
	io.WriteString(h, input)
	hex := hex.EncodeToString([]byte(h.Sum(nil)))
	return hex
}

/*
// Logger code with a bunch more junk in it
func Logger(user *Admin, info chan map[string]interface{}) {
	buf := make([]byte, 69632)
	remains := ""
	count, errcount := 0, 0
	for {

		n, err := user.Conn.Read(buf[:])
		//	fmt.Printf("RECEIVED: %+v\n\n", string(buf))
		if err != nil {
			fmt.Printf("Error reading from socket: %+v\n\n", err)
			return
		}
		length := len(remains) + n
		newest := string(buf[0:n])
		//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
		combined := []byte(remains + newest)
		//combined := append([]byte(remains), buf...)
		//	fmt.Printf("REMAINS %d: %+v\n\n", len(remains), string(remains))
		//	fmt.Printf("BUF %d: %+v\n\n", len(buf), string(buf))
		//	fmt.Printf("COMBINED %d: %+v\n\n", length, string(combined))
		d := bencode.NewDecoder(combined[0:length])
		//d := bencode.NewDecoder(buf[0:n])

		for !d.Consumed {
			o, err := d.Decode()
			//t := reflect.TypeOf(o)
			if err != nil {
				println("***********************************************************************************************")
				fmt.Printf("ERROR! Read %v bytes, err: %+v\n", d.Last, err)
				println("***********************************************************************************************")

				//BUG(inhies) use http://golang.org/pkg/bytes/#IndexByte instead of this hack
				remains = string(buf[d.Last:]) //PROBABLY NEED TO FIX THIS
				//		fmt.Printf("RECEIVED: %+v\n\n", combined)
				//		fmt.Printf("REMAINS: %+v\n\n", remains)
				//return
				errcount++
				break
			}
			remains = remains[:0]
			count++
			info <- o.(map[string]interface{})
			fmt.Printf("%d -- %d \n\n", count, errcount)
			//fmt.Printf("%d -- %d -- Read %d bytes: %#v\n\n", count, errcount, d.Last, o)
		}

	}
}

*/
