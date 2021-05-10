// Package eitherbox provides a cryptographically secure composition of NaCL secretbox that can be opened by either of two key holders.
package eitherbox

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/box"
	"golang.org/x/crypto/nacl/secretbox"
)

// Box is a cryptographically secure box that can be opened with one of two keys.
type Box []byte

// Key represents a NaCL compatible key
type Key *[32]byte

// Encrypt encrypts `b` into a package that can be decrypted by either public key.
// Public/private keypairs should be generated with box.GenerateKey(..).
func Encrypt(b []byte, k1, k2 Key) Box {
	// encrypt `b` with a new random key
	k3 := randomKey()
	nonce := randomNonce()
	ct := secretbox.Seal(nonce[:], b, nonce, k3)

	// encrypt k3 for {k1,k2}
	k3k1, err := box.SealAnonymous(nil, k3[:], k1, rand.Reader)
	if err != nil {
		panic(err)
	}
	k3k2, err := box.SealAnonymous(nil, k3[:], k2, rand.Reader)
	if err != nil {
		panic(err)
	}

	// format is [k3k1,k3k2,ct]
	return append(append(k3k1, k3k2...), ct...)
}

// Decrypt returns the original plaintext, given that this key is one of the original two keys used to create the box.
func (b Box) Decrypt(pub, prv Key) ([]byte, error) {
	if len(b) <= keyBoxSize*2 {
		return nil, fmt.Errorf("double keyed box is invalid")
	}

	// format is [k3k1,k3k2,ct]
	k3k1b := b[0:keyBoxSize]
	k3k2b := b[keyBoxSize : keyBoxSize*2]
	ct := b[keyBoxSize*2:]

	// try both keys
	k3, err := tryDecryptAnonymous(pub, prv, k3k1b, k3k2b)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt double keyed box {k1,k2}: %w", err)
	}

	if len(ct) < 24 {
		return nil, fmt.Errorf("ciphertext is too small") // nonce couldn't even fit
	}

	nonce := new([24]byte)
	copy(nonce[:], ct[:24]) // first 24 bytes of ciphertext is the nonce
	d, ok := secretbox.Open(nil, ct[24:], nonce, k3)
	if !ok {
		return nil, fmt.Errorf("failed to decrypt double keyed box")
	}
	return d, nil
}

const keyBoxSize = 32 + box.AnonymousOverhead

func tryDecryptAnonymous(pub, prv Key, candidates ...[]byte) (Key, error) {
	for _, c := range candidates {
		b, ok := box.OpenAnonymous(nil, c, pub, prv)
		if ok {
			k := new([32]byte)
			copy(k[:], b)
			return k, nil
		}
	}
	return nil, fmt.Errorf("failed to decrypt any candidates")
}

func randomNonce() *[24]byte {
	no := new([24]byte)
	n, err := rand.Read(no[:])
	if err != nil {
		panic(err)
	}
	if n != 24 {
		panic("couldn't fill random bytes")
	}
	return no
}

func randomKey() Key {
	k := new([32]byte)
	n, err := rand.Read(k[:])
	if err != nil {
		panic(err)
	}
	if n != 32 {
		panic("couldn't fill random bytes")
	}
	return k
}
