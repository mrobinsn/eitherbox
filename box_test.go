package eitherbox

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/nacl/box"
)

func TestBox(t *testing.T) {
	k1Pub, k1Prv, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	k2Pub, k2Prv, err := box.GenerateKey(rand.Reader)
	require.NoError(t, err)
	b := []byte("secret")

	dkb := Encrypt(b, k1Pub, k2Pub)
	require.NotEmpty(t, dkb)

	t.Run("k1 can decrypt", func(t *testing.T) {
		result, err := dkb.Decrypt(k1Pub, k1Prv)
		require.NoError(t, err)
		require.Equal(t, b, result)
	})

	t.Run("k2 can decrypt", func(t *testing.T) {
		result, err := dkb.Decrypt(k2Pub, k2Prv)
		require.NoError(t, err)
		require.Equal(t, b, result)
	})

	t.Run("k4 can't decrypt", func(t *testing.T) {
		k4Pub, k4Prv, err := box.GenerateKey(rand.Reader)
		require.NoError(t, err)
		result, err := dkb.Decrypt(k4Pub, k4Prv)
		require.Error(t, err)
		require.Empty(t, result)
	})

	t.Run("tampering with ct", func(t *testing.T) {
		dkb[keyBoxSize*2+1] = byte(0xff)
		_, err := dkb.Decrypt(k1Pub, k1Prv)
		require.Error(t, err)
		_, err = dkb.Decrypt(k2Pub, k2Prv)
		require.Error(t, err)
	})
}
