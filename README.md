# eitherbox

[![Go Reference](https://pkg.go.dev/badge/github.com/mrobinsn/eitherbox.svg)](https://pkg.go.dev/github.com/mrobinsn/eitherbox)
[![Go Report Card](https://goreportcard.com/badge/github.com/mrobinsn/eitherbox)](https://goreportcard.com/report/github.com/mrobinsn/eitherbox)
[![CircleCI](https://circleci.com/gh/mrobinsn/eitherbox.svg?style=svg&circle-token=80fe66828c2588e5cd0d4c80441c1b441d0ba248)](https://circleci.com/gh/mrobinsn/eitherbox)
[![codecov](https://codecov.io/gh/mrobinsn/eitherbox/branch/main/graph/badge.svg?token=5U8FBOQQYC)](https://codecov.io/gh/mrobinsn/eitherbox)

A NaCL based secret box that can be opened with either of two keys.

## Example

```golang
package main

import (
	"crypto/rand"
	"fmt"

	"github.com/mrobinsn/eitherbox"
	"golang.org/x/crypto/nacl/box"
)

func main() {
	// Create keys for Alice
	alicePublic, alicePrivate, _ := box.GenerateKey(rand.Reader)

	// Create keys for Bob
	bobPublic, bobPrivate, _ := box.GenerateKey(rand.Reader)

	// Create keys for Eve
	evePublic, evePrivate, _ := box.GenerateKey(rand.Reader)

	secret := []byte("hello world")

	ebox := eitherbox.Encrypt(secret, alicePublic, bobPublic)

	// Alice can decrypt
	aliceMsg, _ := ebox.Decrypt(alicePublic, alicePrivate)

	// Bob can decrypt
	bobMsg, _ := ebox.Decrypt(bobPublic, bobPrivate)

	// Eve can't decrypt
	eveMsg, _ := ebox.Decrypt(evePublic, evePrivate)

	fmt.Println("Alice got:", string(aliceMsg))
	fmt.Println("Bob got:", string(bobMsg))
	fmt.Println("Eve got:", string(eveMsg))
	// Output: Alice got: hello world
	// Bob got: hello world
	// Eve got:
}
```
