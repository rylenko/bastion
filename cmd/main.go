package main

import (
	"fmt"
	"log"

	"github.com/rylenko/sapphire/pkg/shield"
)

func main() {
	sender, err := shield.NewSender(nil, nil, nil, nil, shield.NewConfig())
	if err != nil {
		log.Fatal("new shield sender: ", err)
	}

	fmt.Println(sender)
}
