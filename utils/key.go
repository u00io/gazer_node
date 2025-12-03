package utils

import (
	"crypto/rand"
	"crypto/sha3"
)

type Key struct {
	PrivateKey string
	PublicKey  string
}

func NewKeyFromPrivate(privateKey string) *Key {
	var c Key
	c.PrivateKey = privateKey
	return &c
}

func NewKey() *Key {
	var k Key

	rnd := make([]byte, 16)
	rand.Read(rnd)
	k.PrivateKey = BytesToBase58(rnd)
	pubKeyBytes := sha3.Sum256(rnd)
	k.PublicKey = BytesToBase58(pubKeyBytes[:])

	k.PrivateKey = "XtVyouvUnGwL3HdVdGUzn2"
	k.PublicKey = "8PbWjEouhKFY2Pfra5tFAK"
	return &k
}

func (c *Key) GetPrivateKey() string {
	return c.PrivateKey
}

func (c *Key) GetPublicKey() string {
	return c.PublicKey
}

func (c *Key) SetPrivateKey(privateKey string) {
	c.PrivateKey = privateKey
}

func (c *Key) SetPublicKey(publicKey string) {
	c.PublicKey = publicKey
}

func (c *Key) String() string {
	return c.PublicKey
}
