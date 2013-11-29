package key

import (
	"crypto/sha512"
	"net"
)

type (
	// Represents a cjdns public key.
	Public struct {
		Key     [32]byte // Raw public key, used for doing crypto
		Encoded string   // String representation ending in .k
		IPv6    net.IP   // IPv6 address
	}
)

// Takes a raw public key and returns a new Public
func NewPublicFromBytes(key [32]byte) (*Public, error) {
	var err error
	keyOut := &Public{Key: key}
	keyOut.Encoded = makeString(key) + ".k"
	keyOut.IPv6, err = keyOut.makeIPv6()
	if err != nil {
		return nil, err
	}
	return keyOut, nil
}

// Takes the string representation of a public key and returns a new Public
func NewPublicFromString(key string) (*Public, error) {
	keyIn := key
	raw := []byte{}
	out := [32]byte{}
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
		return nil, ErrInvalidPubKey
	}
	copy(out[:], raw[:])

	var err error
	keyOut := &Public{Key: out}
	keyOut.Encoded = keyIn
	keyOut.IPv6, err = keyOut.makeIPv6()
	if err != nil {
		return nil, err
	}
	return keyOut, nil
}

// Returns a string containing the IPv6 address for the public key
func (k Public) makeIPv6() (net.IP, error) {
	// Do the hashing that generates the IP
	var out []byte
	h := sha512.New()
	h.Write(k.Key[:])
	out = h.Sum(out[:0])
	h.Reset()
	h.Write(out)
	out = h.Sum(out[:0])[0:16]
	if out[0] != 0xFC {
		return nil, ErrInvalidPubKey
	}
	ip := net.IP(out)
	return ip, nil
}
