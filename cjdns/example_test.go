package cjdns

import (
	"fmt"
)

func ExampleCheckPubkey() {
	rawKey, _ := ParsePubKey("nss00kzfgckxzju9ynlk5qv9gdljh3wwhtc5dw6d7mk9m2fmgg40.k")
	valid := CheckPubKey(rawKey)
	fmt.Println(valid)
	// Output:
	// true
}

func ExamplePadIPv6() {
	addr := PadIPv6("fc2d:39c5:be:8dfa:db9d:bf12:7942:b806")
	fmt.Println(addr)
	// Output:
	// fc2d:39c5:00be:8dfa:db9d:bf12:7942:b806
}

func ExamplePubKeyToIP() {
	addr, _ := PubKeyToIP("nss00kzfgckxzju9ynlk5qv9gdljh3wwhtc5dw6d7mk9m2fmgg40.k")
	fmt.Println(addr)
	// Output:
	// fc2d:39c5:00be:8dfa:db9d:bf12:7942:b806
}
