package main

import (
	"fmt"
	"log"

	"github.com/rylenko/sapphire/pkg/shield"
)

func main() {
	sender, err := shield.NewSender(
		nil,
		nil,
		nil,
		nil,
		shield.WithMessageKeysSkipLimit(0),
		shield.WithSkippedMessageKeys(struct{}{}),
	)
	if err != nil {
		log.Fatal("new shield sender: ", err)
	}

	fmt.Println(sender)
}
