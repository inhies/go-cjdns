package cjdns

import (
	"crypto/sha512"
	"fmt"
	"strings"
)

// ParsePubKey decodes a public key string to a byte slice
func ParsePubKey(key string) (raw []byte, err error) {
	// Check for the trailing .k
	if key[len(key)-2] == '.' && key[len(key)-1] == 'k' {
		key = key[0 : len(key)-2]
	}

	var wide, bits uint
	var i2b = []byte("0123456789bcdfghjklmnpqrstuvwxyz")
	var b2i = func() []byte {
		var ascii [256]byte
		for i := range ascii {
			ascii[i] = 255
		}
		for i, b := range i2b {
			ascii[b] = byte(i)
		}
		return ascii[:]
	}()

	for len(key) > 0 && key[0] != '=' {
		// Add the 5 bits of data corresponding to the next `in` character above existing bits
		wide, key, bits = wide|uint(b2i[int(key[0])])<<bits, key[1:], bits+5
		if bits >= 8 {
			// Remove the least significant 8 bits of data and add it to out
			wide, raw, bits = wide>>8, append(raw, byte(wide)), bits-8
		}
	}

	// If there was padding, there will be bits left, but they should be zero
	if wide != 0 {
		err = fmt.Errorf("extra data at end of decode")
	}
	return
}

// CheckPubKey returns true if a public key is valid.
func CheckPubKey(in []byte) bool {
	// Do the hashing that generates the IP
	var out []byte
	h := sha512.New()
	h.Write(in)
	out = h.Sum(out[:0])
	h.Reset()
	h.Write(out)
	out = h.Sum(out[:0])
	return (out[0] == 0xfc)
}

// Fills out an IPv6 address to the full 32 bytes
// Usefull for string comparisons
func PadIPv6(truncated string) (full string) {
	full = truncated[:4]
	for _, couplet := range strings.SplitN(truncated[5:], ":", 7) {
		if len(couplet) == 4 {
			full = full + ":" + couplet
		} else {
			full = full + fmt.Sprintf(":%04s", couplet)
		}
	}
	return
}

// PubKeyToIP and returns the the IPv6 address for a public key
func PubKeyToIP(key string) (ip string, err error) {
	rawKey, err := ParsePubKey(key)
	if err != nil {
		return "", err
	}

	// Do the hashing that generates the IP
	var out []byte
	h := sha512.New()
	h.Write(rawKey)
	out = h.Sum(out[:0])
	h.Reset()
	h.Write(out)
	out = h.Sum(out[:0])[0:16]

	// Assemble the IP
	for i := 0; i < 16; i++ {
		if i > 0 && i < 16 && i%2 == 0 {
			ip += ":"
		}
		ip += fmt.Sprintf("%02x", out[i])
	}
	return
}
