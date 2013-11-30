package key

import (
	"code.google.com/p/go.crypto/curve25519"
	"crypto/rand"
	"encoding/hex"
)

type (
	// Represents a cjdns private key.
	Private [32]byte
)

// Deocdes the hex representation of a private key.
func DecodePrivate(key string) (*Private, error) {
	keyBytes, err := hex.DecodeString(key)
	if err != nil {
		return nil, err
	}

	var keyOut Private
	copy(keyOut[:], keyBytes[:32])
	return &keyOut, nil
}

// Returns a new randomly generated private key.
func Generate() *Private {
	var priv Private
	privkey := Private(priv)

Start:
	rand.Read(privkey[:])

	/*
		    // not sure if needed
			key[0] &= 248
			key[31] &= 127
			key[31] |= 64
	*/
	pubkey := privkey.makePub()

	// Loop until we get a private key that will create a valid IPv6 address.
	if !pubkey.Valid() {
		goto Start
	}

	return &privkey
}

// Performs ScalarBaseMult on the supplied private key, returning the public key
func (privkey *Private) makePub() *Public {
	var pub [32]byte
	priv := [32]byte(*privkey)

	curve25519.ScalarBaseMult(&pub, &priv)

	pubkey := Public(pub)
	return &pubkey
}

// Returns true if the private key is valid.
func (k *Private) Valid() bool {
	pubkey := k.makePub()
	if !pubkey.Valid() {
		return false
	}
	return true
}

// Returns the public key in base32 format.
func (k *Private) String() string {
	return hex.EncodeToString(k[:])
}

// Returns the associated public key for the supplied private key.
func (k *Private) Pubkey() *Public {
	return k.makePub()
}
