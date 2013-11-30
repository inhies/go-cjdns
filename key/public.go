package key

import (
	"net"
)

// Represents a cjdns public key.
type Public [32]byte

// Takes the string representation of a public key and returns a new Public
func DecodePublic(key string) (*Public, error) {
	raw := []byte{}
	out := [32]byte{}

	// Check for the trailing .k
	if key[len(key)-2:] == ".k" {
		key = key[0 : len(key)-2]
	}

	// Decode the key
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

	// Convert the slice to an array
	copy(out[:], raw[:])
	keyOut := Public(out)

	// Check the key for validitiy
	if !keyOut.Valid() {
		return nil, ErrInvalidPubKey
	}

	return &keyOut, nil
}

// Returns true if k is a valid public key.
func (k *Public) Valid() bool {
	// It's a valid key if the IP address begins with FC
	v := hashTwice(*k)
	return v[0] == 0xFC
}

// Returns the public key in base32 format ending with .k
func (k *Public) String() string {
	return makeString(*k) + ".k"
}

// Retusn the cjdns IPv6 address of the key.
func (k *Public) IP() net.IP {
	return k.makeIPv6()
}

// Returns a string containing the IPv6 address for the public key
func (k *Public) makeIPv6() net.IP {
	out := hashTwice(*k)
	return net.IP(out)
}
