package key

import ()

var (
	privkeyBytes = [32]byte{
		0x75, 0x1d, 0x3d, 0xb8, 0x5b, 0x84, 0x8d, 0xea,
		0xf2, 0x21, 0xe0, 0xed, 0x2b, 0x6c, 0xc1, 0x7f,
		0x58, 0x7b, 0x29, 0x05, 0x7d, 0x74, 0xcd, 0xd4,
		0xdc, 0x0b, 0xd1, 0x8b, 0x71, 0x57, 0x28, 0x8e,
	}
)

type (
	// Represents a cjdns private key.
	Private struct {
		Key     [32]byte
		Encoded string
	}
)

func NewPrivateFromBytes(key [32]byte) (*Private, error) {
	var err error
	keyOut := Private{Key: key}
	keyOut.Encoded = makeString(key)
	if err != nil {
		return nil, err
	}
	return &keyOut, nil
}

func NewPrivateFromString(key string) (*Private, error) {
	return nil, nil
}
