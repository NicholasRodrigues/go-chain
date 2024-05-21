package crypto

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

const (
	privKeyLen = 64
	pubKeyLen  = 32
)

type PrivateKey struct {
	key ed25519.PrivateKey
}

func (p *PrivateKey) Bytes() []byte {
	return p.key[:]
}

func (p *PrivateKey) Sign(msg []byte) []byte {
	return ed25519.Sign(p.key, msg)
}

func (p *PrivateKey) PublicKey() *PublicKey {
	return &PublicKey{key: p.key.Public().(ed25519.PublicKey)}
}

func NewPrivateKey() (*PrivateKey, error) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &PrivateKey{key: priv}, nil
}

type PublicKey struct {
	key ed25519.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.key[:]
}

func (p *PublicKey) Verify(msg, sig []byte) bool {
	return ed25519.Verify(p.key, msg, sig)
}

func (p *PublicKey) String() string {
	return hex.EncodeToString(p.key)
}

func (p *PrivateKey) String() string {
	return hex.EncodeToString(p.key)
}

func PublicKeyFromString(keyStr string) (*PublicKey, error) {
	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != pubKeyLen {
		return nil, fmt.Errorf("invalid public key length")
	}
	return &PublicKey{key: ed25519.PublicKey(keyBytes)}, nil
}

func PrivateKeyFromString(keyStr string) (*PrivateKey, error) {
	keyBytes, err := hex.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	if len(keyBytes) != privKeyLen {
		return nil, fmt.Errorf("invalid private key length")
	}
	return &PrivateKey{key: ed25519.PrivateKey(keyBytes)}, nil
}
