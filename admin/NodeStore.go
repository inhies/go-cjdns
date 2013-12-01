package admin

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const magicalLinkConstant = 5366870 //Determined by cjd way back in the dark ages.

type Route struct {
	IP      string
	Link    Link
	Path    Path
	Version int
}

type (
	Link uint32
	Path int64
)

func (l Link) String() string {
	return strconv.FormatUint(uint64(l)/magicalLinkConstant, 10)
}

func (p Path) String() string {
	str := strconv.FormatInt(int64(p), 10)
	str = strings.Repeat("0", 16-len(str)) + str
	var out string
	for i, s := range str {
		if i > 0 && i%4 == 0 {
			out += "."
		}
		out += string(s)
	}
	return out
}

func ParsePath(path string) Path {
	sPath := strings.Replace(path, ".", "", -1)
	fmt.Println(sPath)
	bPath, err := hex.DecodeString(sPath)
	if err != nil || len(bPath) != 8 {
		fmt.Println("ERROR", err)
		//If we get an error, or the
		//path is not 64 bits, discard.
		//This should also prevent
		//runtime errors.
		return 0
	}
	return Path(binary.BigEndian.Uint64(bPath))

}

type Routes []*Route

func (rs Routes) Len() int      { return len(rs) }
func (rs Routes) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

/*
func (rs Routes) parsePaths() {
	if rs[0].Path != 0 {
		return
	}
	for _, r := range rs {
		h, _ := hex.DecodeString(strings.Replace(r.Path, ".", "", -1))
		r.Path = binary.BigEndian.Uint64(h)
	}
}
*/

// SortByPath sorts Routes by link quality.
func (r Routes) SortByPath() {
	if len(r) < 2 {
		return
	}
	//r.parsePaths()
	sort.Sort(byPath{r})
}

type byPath struct{ Routes }

func (s byPath) Less(i, j int) bool { return s.Routes[i].Path < s.Routes[j].Path }

// SortByQuality sorts Routes by link quality.
func (r Routes) SortByQuality() {
	if len(r) < 2 {
		return
	}
	sort.Sort(byQuality{r})
}

type byQuality struct{ Routes }

func (s byQuality) Less(i, j int) bool { return s.Routes[i].Link > s.Routes[j].Link }

// Log base 2 of a uint64
func log2x64(in uint64) (out uint) {
	for in != 0 {
		in = in >> 1
		out++
	}
	return
}

// return true if packets destined for destination go through midPath.
func isBehind(destination, midPath Path) bool {
	if midPath > destination {
		return false
	}
	mask := ^uint64(0) >> (64 - log2x64(uint64(midPath)))
	return (uint64(destination) & mask) == (uint64(midPath) & mask)
}

// IsBehind returns true if packets destined for Route go through the specified node.
func (r *Route) IsBehind(node *Route) bool {
	return isBehind(r.Path, node.Path)
}

// Return true if destination is 1 hop away from midPath
// WARNING: this depends on implementation quirks of the router and will be broken in the future.
// NOTE: This may have false positives which isBehind() will remove.
func isOneHop(destination, midPath Path) bool {
	if !isBehind(destination, midPath) {
		return false
	}

	// The "why" is here:
	// http://gitboria.com/cjd/cjdns/tree/master/switch/NumberCompress.h#L143
	c := uint64(destination) >> log2x64(uint64(midPath))
	if c&1 != 0 {
		return log2x64(c) == 4
	}
	if c&3 != 0 {
		return log2x64(c) == 7
	}
	return log2x64(c) == 10
}

// Hops returns a Routes object representing a set of hops to a path
func (rs Routes) Hops(path string) (hops Routes) {
	h, _ := hex.DecodeString(strings.Replace(path, ".", "", -1))
	p := binary.BigEndian.Uint64(h)
	//rs.parsePaths()
	for _, r := range rs {
		if isBehind(Path(p), r.Path) {
			hops = append(hops, r)
		}
	}
	return
}

// NodeStore_dumpTable will return cjdns's routing table.
func (c *Conn) NodeStore_dumpTable() (routingTable Routes, err error) {
	more := true
	var page int
	var response map[string]interface{}
	for more {
		args := make(map[string]interface{})
		args["page"] = page

		response, err = SendCmd(c, "NodeStore_dumpTable", args)
		if err != nil {
			return
		}
		if e, ok := response["error"]; ok {
			if e.(string) != "none" {
				err = errors.New("NodeStore_dumpTable: " + e.(string))
				return
			}
		}

		rawTable := response["routingTable"].([]interface{})
		for i := range rawTable {
			r := rawTable[i].(map[string]interface{})
			routingTable = append(routingTable, &Route{
				IP:      r["ip"].(string),
				Link:    Link(r["link"].(int64)),
				Path:    ParsePath(r["path"].(string)),
				Version: int(r["version"].(int64)),
			})
		}

		if more = (response["more"] != nil); more {
			page++
		}
	}
	return
}

/*
// NodePeers returns a Routes object representing the direct peers of target.
func (a *Admin) NodePeers(IP string) (directPeers Routes, err error) {
	if l := len(IP); l > 40 {
		err = errors.New(IP + " is not a valid address")
		return
	} else if l < 40 {
		IP = PadIPv6(IP)
	}

	var table Routes
	table, err = a.NodeStore_dumpTable()
	if err != nil {
		return
	}

	m := make(map[string]*Route)

	for _, nodeA := range table {
		if nodeA.IP != IP {
			continue
		}
		fmt.Println("found", nodeA.IP, "at", nodeA.Path, "in table")

		for _, nodeB := range table {
			if nodeB.IP == IP {
				continue
			}
			fmt.Println("looking at", nodeB.IP, nodeB.Path)
			if isOneHop(nodeA.Path, nodeB.Path) || isOneHop(nodeB.Path, nodeA.Path) {
				fmt.Println(nodeA.Path, "is next to", nodeB.Path)
				if previous, ok := m[nodeB.IP]; !ok || previous.Path > nodeB.Path {
					m[nodeB.IP] = nodeB
				}
			}
		}
	}
	directPeers = make(Routes, len(m))
	var i int
	for _, r := range m {
		directPeers[i] = r
		i++
	}
	return
}
*/
