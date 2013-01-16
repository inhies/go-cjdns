package config

// Thanks to SashaCrofter for the original layout of these structures
type Config struct {
	CorePath                    string         `json:"corePath"`                    //the path to the cjdns executable
	PrivateKey                  string         `json:"privateKey"`                  //the private key for this node (keep it safe)
	PublicKey                   string         `json:"publicKey"`                   //the public key for this node
	IPv6                        string         `json:"ipv6"`                        //this node's IPv6 address as (derived from publicKey)
	AuthorizedPasswords         []AuthPass     `json:"authorizedPasswords"`         //authorized passwords
	Admin                       AdminBlock     `json:"admin"`                       //information for RCP server
	Interfaces                  InterfacesList `json:"interfaces"`                  //interfaces for the switch core
	Router                      RouterBlock    `json:"router"`                      //configuration for the router
	ResetAfterInactivitySeconds int            `json:"resetAfterInactivitySeconds"` //remove cryptoauth sessions after this number of seconds
	PidFile                     string         `json:"pidFile,omitempty"`           //the file to write the PID to, if enabled (disabled by default)
	RawSecurity                 interface{}    `json:"security"`                    //contains raw security info from the JSON that is not nicely unmarshalled
     Security                    SecurityBlock  `json:"-"`                           //usable representation of the security info that can not be saved to JSON
	Version                     int            `json:"version"`                     //the internal config file version (mostly unused)
}

type AuthPass struct {
	Password string `json:"password"` //the password for incoming authorization
}

type AdminBlock struct {
	Bind     string `json:"bind"`     //the port to bind the RCP server to
	Password string `json:"password"` //the password for the RCP server
}

type InterfacesList struct {
	UDPInterface []UDPInterfaceBlock `json:"UDPInterface,omitempty"` //Network interface
	ETHInterface []EthInterfaceBlock `json:"ETHInterface,omitempty"` //Ethernet interface
}

type UDPInterfaceBlock struct {
	Bind      string                `json:"bind"`      //Address to bind to ("0.0.0.0:port")
	ConnectTo map[string]Connection `json:"connectTo"` //Maps connection information to peer details, where the Key is the peer's IPv4 address and port and the Connection contains all of the information about the peer, such as password and public key
}

type EthInterfaceBlock struct {
	Bind      string                `json:"bind"`      //Interface to bind to ("eth0")
	ConnectTo map[string]Connection `json:"connectTo"` //Maps connection information to peer details, where the Key is the peer's MAC address and the Connection contains all of the information about the peer, such as password and public key
     Beacon    int                   `json:"beacon"`    //Sets the beacon state for the ether interface. 0 = disabled, 1 = accept beacons, 2 = send and accept beacons.
}

type Connection struct {
	Password  string `json:"password"`  //the password to connect to the peer node
	PublicKey string `json:"publicKey"` //the peer node's public key
}

type RouterBlock struct {
	Interface RouterInterface `json:"interface"` //interface used for connecting to the cjdns network
	IPTunnel  TunnelInterface `json:"ipTunnel"`  //interface used for connecting to the cjdns network
}

type RouterInterface struct {
	Type      string `json:"type"`                //the type of interface
	TunDevice string `json:"tunDevice,omitempty"` //the persistent interface to use for cjdns (not usually used)

}
type TunnelInterface struct {
	AllowedConnections  []TunnelAllowed `json:"allowedConnections"`  //A list of details for users connecting to us to form an IP tunnel
	OutgoingConnections []string        `json:"outgoingConnections"` //A list of nodes we will connect to in order to form an IP tunnel
}
type TunnelAllowed struct {
	Publickey  string `json:"publicKey"`  //the peer node's public key
	IP4Address string `json:"ip4Address"` //the IPv4 address we will assign to the peer's tunnel (we only need to specify either the IPv4 or IPv6 addresses)
	IP6Address string `json:"ip6Address"` //the IPv6 address we will assign to the peer's tunnel (we only need to specify either the IPv4 or IPv6 addresses)
}

//We can not unmarshall the security section of the config directly to a useable structure, so we manually save and restore the values using the SecurityBlock
//We set them by parsing the RawSecurity interface{} 
//This allows us to easily edit these values in our program. Note that the RawSecurity interface{} is what actaully gets marshalled back in to JSON
//We must parse SecurityBlock and create the proper RawSecurity interface{} before marshalling
type SecurityBlock struct {
	NoFiles bool
	SetUser string
}
