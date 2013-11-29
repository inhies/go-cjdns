package key

import (
	"fmt"
)

var (
	ErrInvalidPubKey = fmt.Errorf("Invalid public key supplied")
)

// Returns the string representation of the public key ("<hex stuff>.k")
func makeString(k [32]byte) string {
	//func EncodePubKey(in []byte) (out []byte) {
	var wide, bits uint
	var i2b = []byte("0123456789bcdfghjklmnpqrstuvwxyz")
	in := k[:]
	out := []byte{}
	for len(in) > 0 {
		// Add the 8 bits of data from the next `in` byte above the existing bits
		wide, in, bits = wide|uint(in[0])<<bits, in[1:], bits+8
		for bits > 5 {
			// Remove the least significant 5 bits and add their character to out
			wide, out, bits = wide>>5, append(out, i2b[int(wide&0x1F)]), bits-5
		}
	}
	// If it wasn't a precise multiple of 40 bits, add some padding based on the remaining bits
	if bits > 0 {
		out = append(out, i2b[int(wide)])
	}
	return string(out)
}
