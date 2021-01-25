package main

import "github.com/gbaranski/lightmq/pkg/broker"

func main() {
	b, err := broker.New(broker.Config{
		Hostname: "0.0.0.0",
		Port:     1883,
	})
	if err != nil {
		panic(err)
	}
	panic(b.Listen())
}
