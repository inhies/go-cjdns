package cjdns

import (
	"fmt"
	"encoding/binary"
	"encoding/hex"
	"sort"
	"strings"
)

var (
	magicalLinkConstant = 5366870.0 //Determined by cjd way back in the dark ages.
)

type Route struct {
	IP      string
	Link    float64
	Path    string
	rawLink int64
	rawPath uint64
}

type Routes []*Route

func (rs Routes) Len() int      { return len(rs) }
func (rs Routes) Swap(i, j int) { rs[i], rs[j] = rs[j], rs[i] }

/*
func (rs Routes) parsePaths() {
	if rs[0].rawPath != 0 {
		return
	}
	for _, r := range rs {
		h, _ := hex.DecodeString(strings.Replace(r.Path, ".", "", -1))
		r.rawPath = binary.BigEndian.Uint64(h)
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

func (s byPath) Less(i, j int) bool { return s.Routes[i].rawPath < s.Routes[j].rawPath }

// SortByQuality sorts Routes by link quality.
func (r Routes) SortByQuality() {
	if len(r) < 2 {
		return
	}
	sort.Sort(byQuality{r})
}

type byQuality struct{ Routes }

func (s byQuality) Less(i, j int) bool { return s.Routes[i].Link > s.Routes[j].Link }

func log2x64(number uint64) uint {
	var out uint
	for number != 0 {
		number = number >> 1
		out++
	}
	return out
}

func isBehind(destination uint64, midPath uint64) bool {
	if midPath > destination {
		return false
	}
	mask := ^uint64(0) >> (64 - log2x64(midPath))
	return (destination & mask) == (midPath & mask)
}

// Hops returns a Routes object representing a set of hops to a path
func (rs Routes) Hops(path string) (hops Routes) {
	h, _ := hex.DecodeString(strings.Replace(path, ".", "", -1))
	p := binary.BigEndian.Uint64(h)
	//rs.parsePaths()
	for _, r := range rs {
		if isBehind(p, r.rawPath) {
			hops = append(hops, r)
		}
	}
	hops.SortByPath()
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
				err = fmt.Errorf("NodeStore_dumpTable:", e.(string))
				return
			}
		}

		rawTable := response["routingTable"].([]interface{})
		for i := range rawTable {
			r := rawTable[i].(map[string]interface{})
			rPath := r["path"].(string)
			sPath := strings.Replace(rPath, ".", "", -1)
			bPath, err := hex.DecodeString(sPath)
			if err != nil || len(bPath) != 8 {
				//If we get an error, or the
				//path is not 64 bits, discard.
				//This should also prevent
				//runtime errors.
				continue
			}
			path := binary.BigEndian.Uint64(bPath)
			routingTable = append(routingTable, &Route{
				IP:      r["ip"].(string),
				Link:    float64(r["link"].(int64)) / magicalLinkConstant,
				Path:    rPath,
				rawPath: path,
				rawLink: r["link"].(int64),
			})

		}

		if more = (response["more"] != nil); more {
			page++
		}
	}
	return
}
