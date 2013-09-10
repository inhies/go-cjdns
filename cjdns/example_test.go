package cjdns

import (
	"fmt"
)

func ExampleCheckPubkey() {
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
