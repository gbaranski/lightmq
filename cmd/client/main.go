package main

import (
	"crypto/ed25519"
	"crypto/rand"

	"github.com/gbaranski/lightmq/pkg/client"
	log "github.com/sirupsen/logrus"
)

func main() {
	pkey, skey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	c := client.New(client.Config{
		ClientID:   "someClientID",
		PublicKey:  pkey,
		PrivateKey: skey,
		Hostname:   "localhost",
		Port:       997,
	})
	err = c.Connect()
	if err != nil {
		panic(err)
	}
	log.Info("Successfully connected!")
	select {}

}
