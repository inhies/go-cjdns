package key

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net"
	"testing"
)

var (
	pubkeyString = "r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k"
	pubkeyBytes  = [32]byte{
		0xd7, 0xc0, 0xdf, 0x45, 0x00, 0x1a, 0x5b, 0xe5,
		0xe8, 0x1c, 0x95, 0xe5, 0x19, 0xbe, 0x51, 0x99,
		0x05, 0x52, 0x37, 0xcb, 0x91, 0x16, 0x88, 0x2c,
		0xad, 0xce, 0xfe, 0x48, 0xab, 0x73, 0x51, 0x73,
	}
	pubkeyIPv6 = "fc68:cb2c:60db:cb96:19ac:34a8:fd34:03fc"
)

func Test_PubkeyFromString(t *testing.T) {
	Convey("Given a pubkey string", t, func() {
		key, err := NewPublicFromString(pubkeyString)
		Convey("It should convert to a Pubkey type", func() {
			So(err, ShouldBeNil)
		})

		Convey("The string representation should be \""+pubkeyString+"\"", func() {
			So(key.Encoded, ShouldEqual, pubkeyString)
		})

		Convey("The IPv6 address should be \""+pubkeyIPv6+"\"", func() {
			netIP := net.ParseIP(pubkeyIPv6)
			So(netIP.Equal(key.IPv6), ShouldBeTrue)
		})
	})
}

func Test_PubkeyFromBytes(t *testing.T) {
	Convey("Given a pubkey byte array", t, func() {
		key, err := NewPublicFromBytes(pubkeyBytes)
		Convey("It should convert to a Pubkey type", func() {
			So(err, ShouldBeNil)
		})

		Convey("The string representation should be \""+pubkeyString+"\"", func() {
			So(key.Encoded, ShouldEqual, pubkeyString)
		})

		Convey("The decoded value should be "+fmt.Sprintf("%x", pubkeyBytes), func() {
			So(key.Key, ShouldResemble, pubkeyBytes)
		})

		Convey("The IPv6 address should be \""+pubkeyIPv6+"\"", func() {
			netIP := net.ParseIP(pubkeyIPv6)
			So(netIP.Equal(key.IPv6), ShouldBeTrue)
		})
	})
}

/*
func ExampleCheckPubKey() {
	rawKey, _ := ParsePubKey("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	valid := CheckPubKey(rawKey)
	fmt.Println(valid)
	// Output:
	// true
}


func ExamplePadIPv6() {
	addr := PadIPv6("fc68:cb2c:60db:cb96:19ac:34a8:fd34:3fc")
	fmt.Println(addr)
	// Output:
	// fc68:cb2c:60db:cb96:19ac:34a8:fd34:03fc
}

func ExamplePubKeyToIP() {
	addr, _ := PubKeyToIP("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	fmt.Println(addr)
	// Output:
	// fc68:cb2c:60db:cb96:19ac:34a8:fd34:03fc
}
func ExampleParsePubKey() {
	key, err := Decode("r6jzx210usqbgnm3pdtm1z6btd14pvdtkn5j8qnpgqzknpggkuw0.k")
	if err != nil {
		fmt.Println(err)
		return
	}
	valid := bytes.Equal(key, []byte{
		0xd7, 0xc0, 0xdf, 0x45, 0x00, 0x1a, 0x5b, 0xe5,
		0xe8, 0x1c, 0x95, 0xe5, 0x19, 0xbe, 0x51, 0x99,
		0x05, 0x52, 0x37, 0xcb, 0x91, 0x16, 0x88, 0x2c,
		0xad, 0xce, 0xfe, 0x48, 0xab, 0x73, 0x51, 0x73,
	})
	fmt.Println(valid)
	// Output:
	// true
}
*/
