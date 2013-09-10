package cjdns_test

import (
	"fmt"
	"github.com/inhies/go-cjdns/cjdns"
)

func ExampleCheckPubkey() {
	rawKey, _ := cjdns.ParsePubKey("nss00kzfgckxzju9ynlk5qv9gdljh3wwhtc5dw6d7mk9m2fmgg40.k")
	valid := cjdns.CheckPubKey(rawKey)
	fmt.Println(valid)
	// Output:
	// true
}

func ExamplePadIPv6() {
	addr := cjdns.PadIPv6("fc2d:39c5:be:8dfa:db9d:bf12:7942:b806")
	fmt.Println(addr)
	// Output:
	// fc2d:39c5:00be:8dfa:db9d:bf12:7942:b806
}

func ExamplePubKeyToIP() {
	addr, _ := cjdns.PubKeyToIP("nss00kzfgckxzju9ynlk5qv9gdljh3wwhtc5dw6d7mk9m2fmgg40.k")
	fmt.Println(addr)
	// Output:
	// fc2d:39c5:00be:8dfa:db9d:bf12:7942:b806
}
