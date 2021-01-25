package main

import (
	"bufio"
	"crypto/ed25519"
	"crypto/rand"
	"os"

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
	r := bufio.NewReader(os.Stdin)
	for {
		text, _ := r.ReadString('\n')
		c.Send([]byte(text))
	}
}
