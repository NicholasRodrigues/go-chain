package users

import (
	"fmt"
	"github.com/NicholasRodrigues/go-chain/pkg/crypto"
)

type User struct {
	Username   string
	Password   string
	PrivateKey *crypto.PrivateKey
	PublicKey  *crypto.PublicKey
	Balance    int
}

func NewUser(username, password string) (*User, error) {
	privKey, err := crypto.NewPrivateKey()
	if err != nil {
		return nil, err
	}
	pubKey := privKey.PublicKey()

	user := &User{
		Username:   username,
		Password:   password,
		PrivateKey: privKey,
		PublicKey:  pubKey,
		Balance:    0,
	}

	fmt.Printf("Debug: Created user %s with public key %x\n", username, pubKey.Bytes())
	return user, nil
}
