package admin

import (
	"encoding/binary"
	"encoding/hex"
	"net"
	"sort"
	"strconv"
	"strings"
)

const magicalLinkConstant = 5366870 //Determined by cjd way back in the dark ages.

type Route struct {
	IP      *net.IP
	Link    Link
	Path    *Path
	Version int
}

type (
	Link uint32
	Path uint64
)

func (l Link) String() string {
	return strconv.FormatUint(uint64(l)/magicalLinkConstant, 10)
}

func (p Path) String() string {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(p))
	text := make([]byte, 19)
	hex.Encode(text, b)
	copy(text[15:19], text[12:16])
	text[14] = '.'
	copy(text[10:14], text[8:12])
	text[9] = '.'
	copy(text[5:9], text[4:8])
	text[4] = '.'
	return string(text)
}

func (p Path) MarshalText() (text []byte, err error) {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(p))
	text = make([]byte, 19)
	hex.Encode(text, b)
	copy(text[15:19], text[12:16])
	text[14] = '.'
	copy(text[10:14], text[8:12])
	text[9] = '.'
	copy(text[5:9], text[4:8])
	text[4] = '.'
	return
}

func ParsePath(path string) Path {
	sPath := strings.Replace(path, ".", "", -1)
	bPath, err := hex.DecodeString(sPath)
	if err != nil || len(bPath) != 8 {
		//If we get an error, or the
		//path is not 64 bits, discard.
		//This should also prevent
		//runtime errors.
		return 0
	}
	return Path(binary.BigEndian.Uint64(bPath))
}

func (p *Path) UnmarshalText(text []byte) error {
	copy(text[4:8], text[5:9])
	copy(text[8:12], text[10:14])
	copy(text[12:16], text[15:19])
	text = text[:16]

	b := make([]byte, 16)

	_, err := hex.Decode(b, text)
	if err != nil {
		return err
	}
	*p = Path(binary.BigEndian.Uint64(b))
	return nil
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

func (s byPath) Less(i, j int) bool { return *s.Routes[i].Path < *s.Routes[j].Path }

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
	return isBehind(*r.Path, *node.Path)
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
func (rs Routes) Hops(path Path) (hops Routes) {
	for _, r := range rs {
		if isBehind(path, *r.Path) {
			hops = append(hops, r)
		}
	}
	return
}

// NodeStore_dumpTable will return cjdns's routing table.
func (c *Conn) NodeStore_dumpTable() (routingTable Routes, err error) {
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

	return resp.RoutingTable, err
}

// Peers returns a Routes object representing routes directly connected to a given IP.
func (rs Routes) Peers(ip net.IP) (peerRoutes Routes) {
	pm := make(map[string]*Route)

	for _, nodeA := range rs {
		if !nodeA.IP.Equal(ip) {
			continue
		}

		for _, nodeB := range rs {
			if isOneHop(*nodeA.Path, *nodeB.Path) || isOneHop(*nodeB.Path, *nodeA.Path) {
				if prev, ok := pm[nodeB.IP.String()]; !ok || *nodeB.Path < *prev.Path {
					// route has not be stored or it is shorter than the previous
					pm[nodeB.IP.String()] = nodeB
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
