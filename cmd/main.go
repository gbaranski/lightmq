package main

import "github.com/gbaranski/lightmq"

func main() {
	b, err := lightmq.NewBroker(lightmq.Config{
		Hostname: "0.0.0.0",
		Port:     1883,
	})
	if err != nil {
		panic(err)
	}
	panic(b.Listen())
}
