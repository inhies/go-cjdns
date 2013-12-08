package admin

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

const magicalLinkConstant = 5366870 // Determined by cjd way back in the dark ages.

type Route struct {
	SwitchLabel string `bencode:"path"`
	IP          string
	Link        Link
	Path        uint64
	Version     int
}

func (r *Route) String() string {
	return r.SwitchLabel
}

type (
	Link uint32
)

func (l Link) String() string {
	return strconv.FormatUint(uint64(l)/magicalLinkConstant, 10)
}

func ParsePath(path string) uint64 {
	sPath := strings.Replace(path, ".", "", -1)
	bPath, err := hex.DecodeString(sPath)
	if err != nil || len(bPath) != 8 {
		//If we get an error, or the
		//path is not 64 bits, discard.
		//This should also prevent
		//runtime errors.
		return 0
	}
	return binary.BigEndian.Uint64(bPath)

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
func isBehind(destination, midPath uint64) bool {
	if midPath > destination {
		return false
	}
	mask := ^uint64(0) >> (64 - log2x64(midPath))
	return (destination & mask) == (midPath & mask)
}

// IsBehind returns true if packets destined for Route go through the specified node.
func (r *Route) IsBehind(node *Route) bool {
	return isBehind(r.Path, node.Path)
}

// Return true if destination is 1 hop away from midPath
// WARNING: this depends on implementation quirks of the router and will be broken in the future.
// NOTE: This may have false positives which isBehind() will remove.
func isOneHop(destination, midPath uint64) bool {
	if !isBehind(destination, midPath) {
		return false
	}

	// The "why" is here:
	// http://gitboria.com/cjd/cjdns/tree/master/switch/NumberCompress.h#L143
	c := destination >> log2x64(midPath)
	if c&1 != 0 {
		return log2x64(c) == 4
	}
	if c&3 != 0 {
		return log2x64(c) == 7
	}
	return log2x64(c) == 10
}

// Hops returns a Routes object representing a set of hops to a path
func (rs Routes) Hops(dest *Route) (hops Routes) {
	for _, node := range rs {
		if isBehind(dest.Path, node.Path) {
			hops = append(hops, node)
		}
	}
	return
}

// NodeStore_dumpTable will return cjdns's routing table.
func (c *Conn) NodeStore_dumpTable() (Routes, error) {
	var (
		args = new(struct {
			Page int `bencode:"page"`
		})
		req = &request{Q: "NodeStore_dumpTable", Args: args}

		resp = new(struct {
			More bool
			// skip this for now just to get the function to work
			RoutingTable Routes
		})

		pack *packet
		err  error
	)

	resp.More = true
	for resp.More {
		resp.More = false
		if pack, err = c.sendCmd(req); err == nil {
			err = pack.Decode(resp)
		}
		if err != nil {
			break
		}
		args.Page++
	}

	// Now sort through the table
	for _, r := range resp.RoutingTable {
		r.Path = ParsePath(r.SwitchLabel)
	}
	return resp.RoutingTable, err
}

// Peers returns a Routes object representing routes directly connected to a given IP.
func (rs Routes) Peers(ip string) (peerRoutes Routes) {
	if len(ip) != 39 {
		full := ip[:4]
		for _, couplet := range strings.SplitN(ip[5:], ":", 7) {
			if len(couplet) == 4 {
				full = full + ":" + couplet
			} else {
				full = full + fmt.Sprintf(":%04s", couplet)
			}
		}
		ip = full
	}

	pm := make(map[string]*Route)

	for _, nodeA := range rs {
		if nodeA.IP != ip {
			continue
		}

		for _, nodeB := range rs {
			if isOneHop(nodeA.Path, nodeB.Path) || isOneHop(nodeB.Path, nodeA.Path) {
				if prev, ok := pm[nodeB.IP]; !ok || nodeB.Path < prev.Path {
					// route has not be stored or it is shorter than the previous
					pm[nodeB.IP] = nodeB
				}
			}
		}
	}

	peerRoutes = make(Routes, len(pm))
	var i int
	for _, route := range pm {
		peerRoutes[i] = route
		i++
	}
	return
}
