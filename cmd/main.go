package main

import (
	"fmt"
	"log"

	"github.com/rylenko/bastion/pkg/ratchet"
	"github.com/rylenko/bastion/pkg/ratchet/keys"
)

func main() {
	sender, err := ratchet.NewSender(
		keys.Public{},
		keys.Root{},
		keys.Header{},
		keys.Header{},
		ratchet.WithMessageKeysSkipLimit(0),
	)
	if err != nil {
		log.Fatal("new ratchet sender: ", err)
	}

	fmt.Println(sender)
}
