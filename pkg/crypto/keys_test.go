package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrivateKey(t *testing.T) {
	privKey, err := NewPrivateKey()
	assert.NoError(t, err)
	assert.NotNil(t, privKey)
	assert.Equal(t, privKeyLen, len(privKey.Bytes()))
}

func TestPublicKey(t *testing.T) {
	privKey, err := NewPrivateKey()
	assert.NoError(t, err)

	pubKey := privKey.PublicKey()
	assert.NotNil(t, pubKey)
	assert.Equal(t, pubKeyLen, len(pubKey.Bytes()))
}

func TestSignAndVerify(t *testing.T) {
	privKey, err := NewPrivateKey()
	assert.NoError(t, err)

	pubKey := privKey.PublicKey()

	msg := []byte("hello, blockchain")
	sig := privKey.Sign(msg)

	assert.True(t, pubKey.Verify(msg, sig))
}

func TestPublicKeyFromString(t *testing.T) {
	privKey, err := NewPrivateKey()
	assert.NoError(t, err)

	pubKey := privKey.PublicKey()
	pubKeyStr := pubKey.String()

	deserializedPubKey, err := PublicKeyFromString(pubKeyStr)
	assert.NoError(t, err)
	assert.Equal(t, pubKey.Bytes(), deserializedPubKey.Bytes())
}

func TestPrivateKeyFromString(t *testing.T) {
	privKey, err := NewPrivateKey()
	assert.NoError(t, err)

	privKeyStr := privKey.String()

	deserializedPrivKey, err := PrivateKeyFromString(privKeyStr)
	assert.NoError(t, err)
	assert.Equal(t, privKey.Bytes(), deserializedPrivKey.Bytes())
}
