package main

import (
	"fmt"
	"github.com/inhies/go-cjdns/admin"
	"github.com/kylelemons/godebug/pretty"
)

func main() {
	fmt.Println("Hello World!")
	conn, err := admin.Connect(nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	results, err := conn.InterfaceController.PeerStats()
	if err != nil {
		fmt.Println(err)
		return
	}
	pretty.Print(results)

}
