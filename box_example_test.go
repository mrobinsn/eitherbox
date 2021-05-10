package eitherbox

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/nacl/box"
)

func ExampleEncrypt() {
	// Create keys for Alice
	alicePublic, alicePrivate, _ := box.GenerateKey(rand.Reader)

	// Create keys for Bob
	bobPublic, bobPrivate, _ := box.GenerateKey(rand.Reader)

	// Create keys for Eve
	evePublic, evePrivate, _ := box.GenerateKey(rand.Reader)

	secret := []byte("hello world")

	twokeyBox := Encrypt(secret, alicePublic, bobPublic)

	// Alice can decrypt
	aliceMsg, _ := twokeyBox.Decrypt(alicePublic, alicePrivate)

	// Bob can decrypt
	bobMsg, _ := twokeyBox.Decrypt(bobPublic, bobPrivate)

	// Eve can't decrypt
	eveMsg, _ := twokeyBox.Decrypt(evePublic, evePrivate)

	fmt.Println("Alice got:", string(aliceMsg))
	fmt.Println("Bob got:", string(bobMsg))
	fmt.Println("Eve got:", string(eveMsg))
	// Output: Alice got: hello world
	// Bob got: hello world
	// Eve got:
}
