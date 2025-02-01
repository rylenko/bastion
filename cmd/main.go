package main

import (
	"fmt"
	"log"

	"github.com/rylenko/bastion/pkg/ratchet"
)

func main() {
	sender, err := ratchet.NewSender(
		nil,
		nil,
		nil,
		nil,
		ratchet.WithMessageKeysSkipLimit(0),
		ratchet.WithSkippedMessageKeysStorage(struct{}{}),
	)
	if err != nil {
		log.Fatal("new ratchet sender: ", err)
	}

	fmt.Println(sender)
}
