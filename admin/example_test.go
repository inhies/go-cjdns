package admin

import (
	"fmt"
	"sort"
)

func ExampleAvailableFunctions() {
	funcs, err := c.Admin_availableFunctions()
	if err != nil {
		fmt.Println(err)
	}
	if len(funcs) == 0 {
		fmt.Println("funcs was zero-length")
	}

	cmdList := make([]string, len(funcs))

	var i int
	for cmd, _ := range funcs {
		cmdList[i] = cmd
		i++
	}

	sort.Strings(cmdList)

	for _, cmd := range cmdList {
		fmt.Printf("%s(", cmd)
		args := funcs[cmd]
		if len(args) > 0 {
			for k, arg := range args {
				fmt.Printf(" %s %s, ", k, arg.Type)
				if arg.Required {
					fmt.Print("(required)")
				}
			}
		}
		fmt.Println(")")
	}
	// Output:
	// AdminLog_subscribe( file String, (required) level String, (required) line Int, (required))
	// AdminLog_unsubscribe( streamId String, (required))
	// Admin_asyncEnabled()
	// Admin_availableFunctions( page Int, (required))
	// AuthorizedPasswords_add( authType Int, (required) password String, (required) user String, (required))
	// AuthorizedPasswords_list()
	// AuthorizedPasswords_remove( user String, (required))
	// Core_exit()
	// Core_initTunnel( desiredTunName String, (required))
	// ETHInterface_beacon( interfaceNumber Int, (required) state Int, (required))
	// ETHInterface_beginConnection( interfaceNumber Int, (required) macAddress String, (required) password String, (required) publicKey String, (required))
	// ETHInterface_new( bindDevice String, (required))
	// InterfaceController_disconnectPeer( pubkey String, (required))
	// InterfaceController_peerStats( page Int, (required))
	// IpTunnel_allowConnection( ip4Address String, (required) ip6Address String, (required) publicKeyOfAuthorizedNode String, (required))
	// IpTunnel_connectTo( publicKeyOfNodeToConnectTo String, (required))
	// IpTunnel_listConnections()
	// IpTunnel_removeConnection( connection Int, (required))
	// IpTunnel_showConnection( connection Int, (required))
	// NodeStore_dumpTable( page Int, (required))
	// RouterModule_lookup( address String, (required))
	// RouterModule_pingNode( path String, (required) timeout Int, (required))
	// SearchRunner_showActiveSearch( number Int, (required))
	// Security_noFiles()
	// Security_setUser( user String, (required))
	// SessionManager_getHandles( page Int, (required))
	// SessionManager_sessionStats( handle Int, (required))
	// SwitchPinger_ping( data String, (required) path String, (required) timeout Int, (required))
	// UDPInterface_beginConnection( address String, (required) interfaceNumber Int, (required) password String, (required) publicKey String, (required))
	// UDPInterface_new( bindAddress String, (required))
	// memory()
	// ping()
}
