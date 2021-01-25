package main

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/gbaranski/lightmq/pkg/broker"
)

func main() {
	pkey, skey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	b, err := broker.New(broker.Config{
		Hostname:   "0.0.0.0",
		Port:       997,
		PrivateKey: skey,
		PublicKey:  pkey,
	})
	if err != nil {
		panic(err)
	}
	panic(b.ListenTCP())
}
